package output

import (
	"os"
	"path/filepath"
	"sync"
)

type FileWriter struct {
	mu   sync.Mutex
	file *os.File
}

func NewFileWriter(filename string) (*FileWriter, error) {
	dir := filepath.Dir(filename)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return &FileWriter{
		file: f,
	}, nil
}

func (w *FileWriter) WriteLine(line string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	_, err := w.file.WriteString(line + "\n")
	return err
}

func (w *FileWriter) Close() error {
	return w.file.Close()
}
