package fx

// type Fx interface {
// 	AddFile(filePath string, prefix string) error
// 	AddPath(dir string) string
// 	AddRoute(path string, data []byte, contentType string)
// 	Router() *http.ServeMux
// }

// type fx struct {
// 	mux      *http.ServeMux
// 	pathless string
// 	api      string
// }

// func NewFx(pathless, apiUrl string) Fx {
// 	f := &fx{
// 		mux:      http.NewServeMux(),
// 		pathless: pathless,
// 		api:      apiUrl,
// 	}
// 	return f
// }

// func (f *fx) Router() *http.ServeMux {
// 	return f.mux
// }

// func (f *fx) AddRoute(path string, data []byte, contentType string) {
// 	compressed := f.Compress(data)
// 	f.Router().HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", contentType)
// 		w.Header().Set("Content-Encoding", "gzip")
// 		w.Write(compressed)
// 	})
// }

// func (f *fx) Compress(data []byte) []byte {
// 	var buf bytes.Buffer
// 	gzipWriter := gzip.NewWriter(&buf)
// 	gzipWriter.Write(data)
// 	gzipWriter.Close()
// 	return buf.Bytes()
// }

// func (f *fx) Serve() {
// 	go func() {
// 		http.ListenAndServe(":1001", f.cors(f.Router()))
// 	}()
// }

// func (f *fx) cors(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Access-Control-Allow-Origin", f.pathless)
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
// 		if r.Method == "OPTIONS" {
// 			w.WriteHeader(http.StatusOK)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }

// // Add a single file to the frame with a prefix path
// func (f *fx) AddFile(filePath string, prefix string) error {
// 	fileData, err := os.ReadFile(filePath)
// 	if err != nil {
// 		return err
// 	}

// 	base := filepath.Base(filePath)
// 	name := base[:len(base)-len(filepath.Ext(base))]
// 	contentType := f.getType(base, fileData)
// 	routePath := "/" + strings.Trim(prefix, "/") + "/" + name

// 	f.AddRoute(routePath, fileData, contentType)
// 	return nil
// }

// // Walk directory and load files into memory, determine Content-Type based on file extension, register routes as /<dirname>/<file without extension>
// func (f *fx) AddPath(dir string) string {
// 	prefix := filepath.Base(dir)
// 	var orderData []byte
// 	var orderContentType string

// 	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
// 		if err != nil || info.IsDir() {
// 			return err
// 		}

// 		base := filepath.Base(path)
// 		name := base[:len(base)-len(filepath.Ext(base))]
// 		contentType := f.getType(base, nil)

// 		if base == "sort.json" {
// 			orderData, _ = os.ReadFile(path)
// 			orderContentType = contentType
// 		} else {
// 			fileData, err := os.ReadFile(path)
// 			if err != nil {
// 				return err
// 			}
// 			contentType = f.getType(base, fileData)
// 			routePath := "/" + prefix + "/" + name
// 			f.AddRoute(routePath, fileData, contentType)
// 		}
// 		return nil
// 	})

// 	if orderData != nil {
// 		routePath := "/" + prefix
// 		f.AddRoute(routePath, orderData, orderContentType)
// 	}
// 	return prefix
// }

// func (f *fx) getType(filename string, data []byte) string {
// 	contentType := mime.TypeByExtension(filepath.Ext(filename))
// 	if contentType == "" {
// 		contentType = http.DetectContentType(data)
// 	}
// 	return contentType
// }
