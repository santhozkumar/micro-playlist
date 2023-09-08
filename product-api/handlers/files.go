package handlers

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/santhozkumar/micro-playlist/product-api/files"
)

type File struct {
	l     *slog.Logger
	store files.Storage
}

func NewFiles(s files.Storage, l *slog.Logger) *File {
	return &File{l: l, store: s}
}


func (f *File) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
    r.Use(middleware.RequestID)

	r.Get("/{id:[0-9]+}/{filename:[a-zA-Z]+.[a-z]{3}}", f.GetFile)
	// r.With(GzipMiddleWare).Get("/{id:[0-9]+}/{filename:[a-zA-Z]+.[a-z]{3}}", f.GetFile)
	r.Post("/{id:[0-9]+}/{filename:[a-zA-Z]+.[a-z]{3}}", f.FileUpload)

	r.Post("/multi", f.FileMultiPart)
	return r
}

func (f *File) GetFile(w http.ResponseWriter, r *http.Request) {
    requestID := r.Context().Value(middleware.RequestIDKey)
    f.l.Log(r.Context(), slog.LevelInfo, "Getting file for", "reqeustID", requestID)

    root, err := filepath.Abs(files.UploadPath)
    if err != nil {
        http.Error(w, "Client Error", http.StatusInternalServerError)
    }
    fs := http.StripPrefix("/file/", http.FileServer(http.Dir(root)))
    fmt.Println("file sent")
    fs.ServeHTTP(w, r)
    fmt.Println("file sent successfully")
}

func (f *File) FileUpload(w http.ResponseWriter, r *http.Request) {
	f.l.Log(r.Context(), slog.LevelInfo, "Handle POST File Upload ")
	id := chi.URLParam(r, "id")
	filename := chi.URLParam(r, "filename")

	if filename == "" || id == "" {
		f.InvalidURI(r.URL.String(), w)
	}
	err := f.fileSave(id, filename, w, r.Body)
	if err != nil {
		f.l.Log(r.Context(), slog.LevelInfo, "failed to upload", err)
		http.Error(w, "File Upload Failed", http.StatusBadRequest)
        return
	}
}

func (f *File) FileMultiPart(w http.ResponseWriter, r *http.Request) {
    f.l.Log(r.Context(), slog.LevelInfo, "Handling Multipart post file upload")
    err := r.ParseMultipartForm(128 * 1024)
    if err != nil {
        f.l.Log(context.Background(), slog.LevelInfo, "No multipart form data")
        http.Error(w, "No multipart form data", http.StatusBadRequest)
        return
    }

    id := r.FormValue("id")
    if _, err = strconv.Atoi(id); err != nil {
        f.l.Log(context.Background(), slog.LevelInfo, "id expected to be int")
        http.Error(w, "id expected to be int", http.StatusBadRequest)
        return
    }

    file, mh, err := r.FormFile("file")
    
    if err != nil {
        f.l.Log(context.Background(), slog.LevelInfo, "File error or not found")
        http.Error(w, "File error or not found", http.StatusBadRequest)
        return
    }
    f.fileSave(id, mh.Filename, w, file)
}


func (f *File) InvalidURI(url string, w http.ResponseWriter) {
	f.l.Log(context.Background(), slog.LevelInfo, fmt.Sprintf("Invalid URL String : %s", url))
	http.Error(w, fmt.Sprintf("Invalid URL String : %s", url), http.StatusBadRequest)
}

func (f *File) fileSave(id string, filename string, w http.ResponseWriter, r io.Reader) error {
	local_filename := filepath.Join(id, filename)
	err := f.store.Save(local_filename, r)
	if err != nil {
        return err
	}
	return nil
}


