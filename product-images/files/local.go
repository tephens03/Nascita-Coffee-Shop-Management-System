package files

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Local struct {
	maxFileSize int // maximum numbber of bytes for files
	basePath    string
}

func NewLocal(maxFileSize int, basePath string) (*Local, error) {
	p, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}

	return &Local{basePath: p, maxFileSize: maxFileSize}, nil
}

func (l *Local) Save(filename string, contents io.Reader) error {
	fp := l.fullPath(filename)

	dir := filepath.Dir(fp)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Cannot create file directory", err)
	}

	// if the file exists delete it
	_, err = os.Stat(fp)
	if err == nil {
		err = os.Remove(fp)
		if err != nil {
			return fmt.Errorf("Unable to delete file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		// if this is anything other than a not exists error
		return fmt.Errorf("Unable to get file info: %w", err)
	}

	// create a new file at the path
	f, err := os.Create(fp)
	if err != nil {
		return fmt.Errorf("Unable to create file: %w", err)
	}
	defer f.Close()

	// write the contents to the new file
	// ensure that we are not writing greater than max bytes
	_, err = io.Copy(f, contents)
	if err != nil {
		return fmt.Errorf("Unable to write to file: %w", err)
	}

	return nil

}

// returns the absolute path
func (l *Local) fullPath(path string) string {
	// append the given path to the base path
	return filepath.Join(l.basePath, path)
}
