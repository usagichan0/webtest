package main

import (
	"archive/zip"
	"bytes"
	"encoding/hex"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type PageData struct {
	Title         string
	PostData      map[string][]string
	UploadResults []UploadResult
}

type UploadResult struct {
	Filename    string
	ContentType string
	Size        int64
	HexDump     string
}

// ヘルパー：レイアウトと各コンテンツテンプレートを結合してレンダリング
func renderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	tmpl, err := template.ParseFiles(
		filepath.Join("templates", "layout.html"),
		filepath.Join("templates", tmplName),
	)
	if err != nil {
		http.Error(w, "Template Parse Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "Template Execute Error: "+err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// Render.com は環境変数 PORT に動的にポート番号を割り当てます
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// indexページ
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		renderTemplate(w, "index.html", nil)
	})

	// オートフィルテスト用ページ（GET/POST両対応）
	http.HandleFunc("/autofill", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{}
		if r.Method == http.MethodPost {
			r.ParseForm()
			data.PostData = r.Form
		}
		renderTemplate(w, "autofill.html", data)
	})

	// Web会議機能テスト用ページ
	http.HandleFunc("/media", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "media.html", nil)
	})

	// ファイル操作テスト用ページ
	http.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "files.html", nil)
	})

	// ポップアップテスト用ページ
	http.HandleFunc("/popups", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "popups.html", nil)
	})

	// ポップアップ先ページ
	http.HandleFunc("/popups/target", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "popup_target.html", nil)
	})

	// ダミーZIPダウンロードエンドポイント
	http.HandleFunc("/download/zip", func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		zw := zip.NewWriter(buf)
		
		f, err := zw.Create("test_document.txt")
		if err == nil {
			f.Write([]byte("This is a dummy file included in the downloaded ZIP archive.\n"))
		}
		zw.Close()

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename=\"dummy_download.zip\"")
		w.Write(buf.Bytes())
	})

	// ダミーPDFダウンロードエンドポイント
	http.HandleFunc("/download/pdf", func(w http.ResponseWriter, r *http.Request) {
		// 最小構成のダミーPDFデータ
		pdfData := `%PDF-1.4
1 0 obj <</Type /Catalog /Pages 2 0 R>> endobj
2 0 obj <</Type /Pages /Kids [3 0 R] /Count 1>> endobj
3 0 obj <</Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Resources <<>> /Contents 4 0 R>> endobj
4 0 obj <</Length 51>> stream
BT /F1 24 Tf 100 700 Td (This is a dummy PDF file.) Tj ET
endstream endobj
xref
0 5
0000000000 65535 f 
0000000009 00000 n 
0000000056 00000 n 
0000000111 00000 n 
0000000212 00000 n 
trailer <</Size 5 /Root 1 0 R>>
startxref
314
%%EOF`
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=\"dummy_download.pdf\"")
		w.Write([]byte(pdfData))
	})

	// ファイルアップロードハンドラ
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/files", http.StatusSeeOther)
			return
		}

		err := r.ParseMultipartForm(50 << 20) // 最大50MBをメモリパース
		if err != nil {
			http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		data := PageData{UploadResults: []UploadResult{}}

		// アップロードされた全ファイルを処理 (単一、複数、DnDいずれも対象)
		for _, fileHeaders := range r.MultipartForm.File {
			for _, hdr := range fileHeaders {
				file, err := hdr.Open()
				if err != nil {
					continue
				}
				
				// 先頭32バイトを読み込む
				head := make([]byte, 32)
				n, _ := file.Read(head)
				head = head[:n]
				file.Close()

				contentType := hdr.Header.Get("Content-Type")
				if contentType == "" {
					contentType = "application/octet-stream"
				}

				res := UploadResult{
					Filename:    hdr.Filename,
					ContentType: contentType,
					Size:        hdr.Size,
					// encoding/hex の Dump関数でアスキー文字列も含めたHEXダンプを生成する
					HexDump:     hex.Dump(head),
				}
				data.UploadResults = append(data.UploadResults, res)
			}
		}

		renderTemplate(w, "files.html", data)
	})

	log.Printf("サーバーをポート %s で起動しています...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
