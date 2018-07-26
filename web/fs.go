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

// ListFiles returns all files in FS
func (fs *BundledFS) ListFiles() (files []macaron.TemplateFile) {
	for _, filename := range fs.List() {
		files = append(files, &BundledFile{fs: fs, FileName: filename})
	}
	return files
}

// Get returns the content of filename
func (fs *BundledFS) Get(filename string) (io.Reader, error) {
	return bytes.NewReader(fs.Bytes(filename)), nil
}

// BundledFile represents a file in a BundledFS
type BundledFile struct {
	fs       *BundledFS
	FileName string
}

// Name represents the name of the file
func (b *BundledFile) Name() string {
	return strings.TrimSuffix(b.FileName, path.Ext(b.FileName))
}

// Data returns the content of file
func (b *BundledFile) Data() []byte {
	return b.fs.Bytes(b.FileName)
}

// Ext returns the file extension
func (b *BundledFile) Ext() string {
	return path.Ext(b.FileName)
}
