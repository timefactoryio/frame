package zero

import (
	"bytes"
	"compress/gzip"
	"net/http"
)

type Fx interface {
	AddFile(filePath string, prefix string) error
	AddPath(dir string) string
	AddRoute(path string, data []byte, contentType string)
	Pathless() string
	Api() string
	Router() *http.ServeMux
	ToBytes(input string) []byte
	ToString(input string) string
}

type fx struct {
	mux      *http.ServeMux
	pathless string
	api      string
}

func NewFx(pathless, apiUrl string) Fx {
	if pathless == "" {
		pathless = "http://localhost:1000"
	}
	if apiUrl == "" {
		apiUrl = "http://localhost:1001"
	}
	f := &fx{
		mux:      http.NewServeMux(),
		pathless: pathless,
		api:      apiUrl,
	}
	return f
}

func (f *fx) Router() *http.ServeMux {
	return f.mux
}

func (f *fx) Pathless() string {
	return f.pathless
}

func (f *fx) Api() string {
	return f.api
}

func (f *fx) AddRoute(path string, data []byte, contentType string) {
	compressed := f.Compress(data)
	f.Router().HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Encoding", "gzip")
		w.Write(compressed)
	})
}

func (f *fx) Compress(data []byte) []byte {
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	gzipWriter.Write(data)
	gzipWriter.Close()
	return buf.Bytes()
}

func (f *fx) Serve() {
	go func() {
		http.ListenAndServe(":1001", f.cors(f.Router()))
	}()
}

func (f *fx) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", f.pathless)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
