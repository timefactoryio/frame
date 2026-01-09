package zero

import (
	"encoding/json"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"
)

func (c *circuit) ReaderV2(path string, compress bool) string {
	dirName, values, order := c.loadFiles(path)
	c.sortFiles(values, order)
	c.value[dirName] = values
	jsonData, _ := json.Marshal(values)
	if compress {
		jsonData = c.Compress(jsonData)
	}
	c.registerRoute(dirName, jsonData, compress)
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

func (c *circuit) registerRoute(dirName string, jsonData []byte, compress bool) {
	c.Router().HandleFunc("/"+dirName, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if compress {
			w.Header().Set("Content-Encoding", "gzip")
		}
		w.Write(jsonData)
	})
}
