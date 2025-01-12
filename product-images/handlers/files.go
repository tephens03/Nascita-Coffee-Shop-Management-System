package handlers

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/sgbaotran/Nascita-coffee-shop/product-images/files"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

type Files struct {
	log     hclog.Logger
	storage files.Storage
}

func NewFile(log hclog.Logger, localStorage files.Storage) *Files {
	return &Files{log, localStorage}
}

func (f *Files) UploadFileREST(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fn := vars["filename"]

	if id == "" || fn == "" {
		f.invalidURI(r.URL.String(), rw)
		return
	}

	f.saveFile(id, fn, rw, r.Body)

}
func (f *Files) UploadFileMultipart(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1024 * 128)
	if err != nil {
		http.Error(rw, "Expected multipart form data", http.StatusInternalServerError)
		return
	}
	id, idErr := strconv.Atoi(r.FormValue("id"))
	if idErr != nil {
		http.Error(rw, "Invalid Id", http.StatusInternalServerError)
		return
	}
	f.log.Info("Processing for ID", id)
	file, file_header, fileErr := r.FormFile("file")
	if fileErr != nil {
		http.Error(rw, "Invalid File", http.StatusInternalServerError)
		return
	}
	f.saveFile(r.FormValue("id"), file_header.Filename, rw, file)

}

func (f *Files) invalidURI(uri string, rw http.ResponseWriter) {
	f.log.Error(uri)
	http.Error(rw, "Invalid file path should be in format /[id]/[filename] :(", http.StatusBadRequest)
}

func (f *Files) saveFile(id, filename string, rw http.ResponseWriter, r io.ReadCloser) {
	f.log.Info("Trying to save file for Product #", id, ", file ", filename)
	fp := filepath.Join(id, filename)
	f.log.Info(fp)

	err := f.storage.Save(fp, r)
	if err != nil {
		f.log.Error("Unable to save file", "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}
}
