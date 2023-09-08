package files

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

var UploadPath string = "./uploads"

type Storage interface {
    Save(string, io.Reader) error
}

type LocalStorage struct {
    dir string 
}

func NewLocalStorage(basePath string) (*LocalStorage, error){
    fullpath, err := filepath.Abs(basePath)
    if err != nil {
        return nil, err
    }
    return &LocalStorage{dir:fullpath}, nil

}

func (l *LocalStorage) Save(filename string, contents io.Reader) error {
    fullPath := filepath.Join(l.dir, filename)

    dir := filepath.Dir(fullPath)
    err := os.MkdirAll(dir, fs.ModePerm)

    if err != nil && !os.IsExist(err) {
        return err
    }

    f, err := os.Create(fullPath)
    if err != nil {
        return err
    }
    defer f.Close()

    _, err = io.Copy(f, contents)
    if err != nil {
        return err
    }
    return nil
}


