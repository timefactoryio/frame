package zero

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Circuit interface {
	Router() *http.ServeMux
	Read(path, prefix string)
	Reader(path string) string
	Pathless() string
	ToBytes(input string) []byte
	ToString(input string) string
	Compress(data []byte) []byte
	Value() map[string][]*Value
}

type circuit struct {
	router   *http.ServeMux
	pathless string
	value    map[string][]*Value
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
		value:    make(map[string][]*Value),
	}
}

func (c *circuit) Value() map[string][]*Value {
	return c.value
}

func (c *circuit) Router() *http.ServeMux {
	return c.router
}

func (c *circuit) Pathless() string {
	return c.pathless
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

func (c *circuit) ToString(input string) string {
	b := c.ToBytes(input)
	if b == nil {
		return ""
	}
	return string(b)
}

func (c *circuit) Compress(data []byte) []byte {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	w.Write(data)
	w.Close()
	return buf.Bytes()
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
		Data: data,
	}

	c.value[prefix] = []*Value{v}

	compressed := c.Compress(v.Data)
	c.Router().HandleFunc("/"+prefix+"/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", v.Type)
		w.Header().Set("Content-Encoding", "gzip")
		w.Write(compressed)
	})
}

func (c *circuit) Reader(path string) string {
	dirName, values, order := c.loadFiles(path)
	c.sortFiles(values, order)
	c.value[dirName] = values
	jsonData, _ := json.Marshal(values)
	c.registerRoute(dirName, c.Compress(jsonData))
	return dirName
}

func (c *circuit) loadFiles(path string) (string, []*Value, []string) {
	dirName := filepath.Base(path)
	var values []*Value
	var order []string

	filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		base := filepath.Base(p)
		if base == "sort.json" {
			if data, err := os.ReadFile(p); err == nil {
				json.Unmarshal(data, &order)
			}
			return nil
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return nil
		}
		ext := filepath.Ext(base)
		ct := mime.TypeByExtension(ext)
		if ct == "" {
			ct = http.DetectContentType(data)
		}
		values = append(values, &Value{Name: base[:len(base)-len(ext)], Type: ct, Data: data})
		return nil
	})

	return dirName, values, order
}

func (c *circuit) sortFiles(values []*Value, order []string) {
	if len(order) == 0 {
		return
	}
	orderMap := make(map[string]int, len(order))
	for i, name := range order {
		orderMap[name] = i
	}
	sort.Slice(values, func(i, j int) bool {
		posI, foundI := orderMap[values[i].Name]
		posJ, foundJ := orderMap[values[j].Name]
		if foundI && foundJ {
			return posI < posJ
		}
		if foundI {
			return true
		}
		if foundJ {
			return false
		}
		return values[i].Name < values[j].Name
	})
}

func (c *circuit) registerRoute(dirName string, jsonData []byte) {
	c.Router().HandleFunc("/"+dirName, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Content-Encoding", "gzip")
		w.Write(jsonData)
	})
}
