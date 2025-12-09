package zero

import (
	"bytes"
	"compress/gzip"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Fx interface {
	AddFile(filePath string, prefix string) error
	AddPath(dir string) string
	// EmbedPath(efs embed.FS, root string) error
	Pathless() string
	Api() string
	Serve()
	Router() *mux.Router
	ToBytes(input string) []byte
}

type fx struct {
	router   *mux.Router
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
		router:   mux.NewRouter(),
		pathless: pathless,
		api:      apiUrl,
	}
	f.router.Use(f.cors())
	return f
}

func (f *fx) Router() *mux.Router {
	return f.router
}

func (f *fx) Pathless() string {
	return f.pathless
}

func (f *fx) Api() string {
	return f.api
}

func (f *fx) Serve() {
	go func() {
		http.ListenAndServe(":1001", f.Router())
	}()
}

func (f *fx) cors() mux.MiddlewareFunc {
	return handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "X-Frame"}),
		handlers.AllowedOrigins([]string{f.pathless}),
		handlers.AllowedMethods([]string{"GET", "OPTIONS"}),
		handlers.ExposedHeaders([]string{"X-Frame", "X-Frames"}),
	)
}

// Add a single file to the frame with a prefix path
func (f *fx) AddFile(filePath string, prefix string) error {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	base := filepath.Base(filePath)
	name := base[:len(base)-len(filepath.Ext(base))]
	contentType := f.getType(base, fileData)
	routePath := "/" + strings.Trim(prefix, "/") + "/" + name

	compress := !strings.HasPrefix(contentType, "video/")
	f.addRoute(routePath, fileData, contentType, compress)
	return nil
}

// Walk directory and load files into memory, determine Content-Type based on file extension, register routes as /<dirname>/<file without extension>
func (f *fx) AddPath(dir string) string {
	prefix := filepath.Base(dir)
	var orderData []byte
	var orderContentType string

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		base := filepath.Base(path)
		name := base[:len(base)-len(filepath.Ext(base))]
		contentType := f.getType(base, nil)

		if base == "sort.json" {
			orderData, _ = os.ReadFile(path)
			orderContentType = contentType
		} else {
			fileData, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			contentType = f.getType(base, fileData)
			routePath := "/" + prefix + "/" + name
			compress := !strings.HasPrefix(contentType, "video/")
			f.addRoute(routePath, fileData, contentType, compress)
		}
		return nil
	})

	if orderData != nil {
		routePath := "/" + prefix
		compress := !strings.HasPrefix(orderContentType, "video/")
		f.addRoute(routePath, orderData, orderContentType, compress)
	}
	return prefix
}

// New method for embedded filesystem
// func (f *fx) EmbedPath(efs embed.FS, root string) error {
// 	return fs.WalkDir(efs, root, func(path string, d fs.DirEntry, err error) error {
// 		if err != nil || d.IsDir() {
// 			return err
// 		}

// 		fileData, err := efs.ReadFile(path)
// 		if err != nil {
// 			return err
// 		}

// 		base := filepath.Base(path)
// 		name := base[:len(base)-len(filepath.Ext(base))]
// 		contentType := f.getType(base, fileData)
// 		routePath := "/" + filepath.Base(root) + "/" + name

// 		f.addRoute(routePath, fileData, contentType)
// 		return nil
// 	})
// }

func (f *fx) getType(filename string, data []byte) string {
	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}
	return contentType
}

func (f *fx) addRoute(path string, data []byte, contentType string, compress bool) {
	f.Router().HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		if compress {
			var buf bytes.Buffer
			gzipWriter := gzip.NewWriter(&buf)
			gzipWriter.Write(data)
			gzipWriter.Close()
			w.Header().Set("Content-Encoding", "gzip")
			w.Write(buf.Bytes())
		} else {
			http.ServeContent(w, r, path, time.Time{}, bytes.NewReader(data))
		}
	})
}

func (f *fx) ToBytes(input string) []byte {
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		resp, err := http.Get(input)
		if err != nil {
			return nil
		}
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil
		}
		return b
	}
	b, err := os.ReadFile(input)
	if err != nil {
		return nil
	}
	return b
}
