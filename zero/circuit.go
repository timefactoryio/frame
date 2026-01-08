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
)

type Circuit interface {
	Router() *http.ServeMux
	Handle(data []byte) http.HandlerFunc
	Reader(path string) string
	Pathless() string
	ToBytes(input string) []byte
	Compress(data []byte) []byte
}

type circuit struct {
	router   *http.ServeMux
	pathless string
}

type Value struct {
	Name string
	Type string
	Data []byte
}

func NewCircuit(pathless string) Circuit {
	if pathless == "" {
		pathless = "http://localhost:1000"
	}

	return &circuit{
		router:   http.NewServeMux(),
		pathless: pathless,
	}
}

func (c *circuit) Router() *http.ServeMux {
	return c.router
}

func (c *circuit) Pathless() string {
	return c.pathless
}

func (c *circuit) Handle(data []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Content-Encoding", "gzip")
		w.Write(data)
	}
}

// Reader reads files from a path (file or directory) and returns the directory name if a directory is provided
func (c *circuit) Reader(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return ""
	}

	if !info.IsDir() {
		c.Read(path, "")
		return ""
	}

	dirName := filepath.Base(path)
	filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		c.Read(p, dirName)
		return nil
	})
	return dirName
}

// Read reads a single file and returns a Value struct with its content, content type, and filename without extension
func (c *circuit) Read(path, prefix string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	base := filepath.Base(path)
	name := base[:len(base)-len(filepath.Ext(base))]
	content := mime.TypeByExtension(filepath.Ext(path))
	if content == "" {
		content = http.DetectContentType(data)
	}

	v := &Value{
		Name: name,
		Type: content,
		Data: c.Compress(data),
	}
	c.addRoute(prefix, v)
}

func (c *circuit) addRoute(prefix string, v *Value) {
	var path string
	if v.Type == "application/json" {
		if prefix != "" {
			path = "/" + prefix
		} else {
			path = "/" + v.Name
		}
	} else {
		if prefix != "" {
			path = "/" + prefix + "/" + v.Name
		} else {
			path = "/" + v.Name
		}
	}

	c.Router().HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", v.Type)
		w.Header().Set("Content-Encoding", "gzip")
		w.Write(v.Data)
	})
}

func (c *circuit) ToBytes(input string) []byte {
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

func (c *circuit) Compress(data []byte) []byte {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	w.Write(data)
	w.Close()
	return buf.Bytes()
}
