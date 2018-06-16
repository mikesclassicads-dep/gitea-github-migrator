package web

import (
	"bytes"
	"io"
	"path"
	"strings"

	"github.com/gobuffalo/packr"
	"gopkg.in/macaron.v1"
)

// BundledFS implements ServeFileSystem for packr.Box
type BundledFS struct {
	packr.Box
}

// Exists returns true if filepath exists
func (fs *BundledFS) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		return fs.Has(p)
	}
	return false
}

func (fs *BundledFS) ListFiles() (files []macaron.TemplateFile) {
	for _, filename := range fs.List() {
		files = append(files, &BundledFile{fs: fs, FileName: filename})
	}
	return files
}

func (fs *BundledFS) Get(filename string) (io.Reader, error) {
	return bytes.NewReader(fs.Bytes(filename)), nil
}

type BundledFile struct {
	fs       *BundledFS
	FileName string
}

func (b *BundledFile) Name() string {
	return strings.TrimSuffix(b.FileName, path.Ext(b.FileName))
}

func (b *BundledFile) Data() []byte {
	return b.fs.Bytes(b.FileName)
}

func (b *BundledFile) Ext() string {
	return path.Ext(b.FileName)
}
