package middleware

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()

		writer := gzipWriter{
			ResponseWriter: w,
			Writer:         gz,
		}
		next.ServeHTTP(writer, r)
	})
}

func DecompressGZIP(next http.Handler) http.Handler {
	// приводим возвращаемую функцию к типу функций HandlerFunc
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(`Content-Encoding`) == `gzip` { //	если входящий пакет сжат GZIP
			gz, err := gzip.NewReader(r.Body) //	изготавливаем reader-декомпрессор GZIP
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				log.Println("Request Body decompression error: " + err.Error())
				return
			}
			r.Body = gz //	подменяем стандартный reader из Request на декомпрессор GZIP
			defer gz.Close()
		}
		next.ServeHTTP(w, r) // передаём управление следующему обработчику
	})
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}
