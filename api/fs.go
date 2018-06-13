package api

import (
	"strings"

	"github.com/gobuffalo/packr"
)

type BundledFS struct {
	packr.Box
}

func (fs *BundledFS) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		return fs.Has(p)
	}
	return false
}
