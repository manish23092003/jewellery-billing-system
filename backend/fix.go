package main

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	err := filepath.WalkDir("internal", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".go" {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		content = bytes.ReplaceAll(content, []byte("\"jewellery-billing/pkg/response\""), []byte("\"jewellery-billing/internal/apiresponse\""))
		content = bytes.ReplaceAll(content, []byte("response.Success"), []byte("apiresponse.Success"))
		content = bytes.ReplaceAll(content, []byte("response.Error"), []byte("apiresponse.Error"))
		content = bytes.ReplaceAll(content, []byte("response.BadRequest"), []byte("apiresponse.BadRequest"))
		content = bytes.ReplaceAll(content, []byte("response.Unauthorized"), []byte("apiresponse.Unauthorized"))
		content = bytes.ReplaceAll(content, []byte("response.Forbidden"), []byte("apiresponse.Forbidden"))
		content = bytes.ReplaceAll(content, []byte("response.NotFound"), []byte("apiresponse.NotFound"))
		content = bytes.ReplaceAll(content, []byte("response.InternalError"), []byte("apiresponse.InternalError"))
		content = bytes.ReplaceAll(content, []byte("response.Meta"), []byte("apiresponse.Meta"))

		return os.WriteFile(path, content, 0644)
	})
	if err != nil {
		panic(err)
	}
}


