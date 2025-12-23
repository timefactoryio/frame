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

type Zero interface {
	Forge
	Router() *http.ServeMux
	Output() []*Value
	Input(path string) string
	ToBytes(input string) []byte
	ToString(input string) string
	Pathless() string
	Api() string
}

type zero struct {
	Forge
	router   *http.ServeMux
	values   []*Value
	pathless string
	api      string
}

type Value struct {
	Name string
	Type string
	Data []byte
}

func NewZero(pathless, apiUrl string) Zero {
	return &zero{
		Forge:    NewForge().(*forge),
		router:   http.NewServeMux(),
		values:   []*Value{},
		pathless: pathless,
		api:      apiUrl,
	}
}

func (z *zero) Pathless() string {
	return z.pathless
}

func (z *zero) Api() string {
	return z.api
}

func (z *zero) Router() *http.ServeMux {
	return z.router
}

func (z *zero) Output() []*Value {
	return z.values
}

func (z *zero) Input(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return ""
	}

	if !info.IsDir() {
		if val := z.Read(path); val != nil {
			z.values = append(z.values, val)
		}
		return ""
	}

	dirName := filepath.Base(path)

	filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if val := z.Read(p); val != nil {
			z.values = append(z.values, val)
		}
		return nil
	})

	return dirName
}

// Read reads a single file and returns a Value struct with its content, content type, and filename without extension
func (z *zero) Read(path string) *Value {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	base := filepath.Base(path)
	name := base[:len(base)-len(filepath.Ext(base))]
	content := mime.TypeByExtension(filepath.Ext(path))
	if content == "" {
		content = http.DetectContentType(data)
	}

	// Compress 10KB+
	if len(data) >= 10240 {
		data = z.Compress(data)
	}

	return &Value{
		Name: name,
		Type: content,
		Data: data,
	}
}

func (z *zero) Compress(data []byte) []byte {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	w.Write(data)
	w.Close()
	return buf.Bytes()
}

func (z *zero) ToBytes(input string) []byte {
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

func (z *zero) ToString(input string) string {
	b := z.ToBytes(input)
	if b == nil {
		return ""
	}
	return string(b)
}
