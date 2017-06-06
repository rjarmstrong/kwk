package app

import (
	"github.com/kwk-super-snippets/cli/src/store"
	//"strings"
	"strings"
)

func NewSnippetReadWriter(file store.File) store.SnippetReadWriter {
	return &snippetReadWriter{file: file}
}

type snippetReadWriter struct {
    file store.File
}

const subDir = "snippets"

func makePath(uri string) string {
	return strings.Replace(strings.TrimPrefix(uri, "/"), "/", "-", -1)
}

func (sr *snippetReadWriter) Write(uri string, content string) (string, error) {
	return sr.file.Write(subDir, makePath(uri), content, true)
}

func (sr *snippetReadWriter) Read(uri string) (string, error) {
	return sr.file.Read(subDir, makePath(uri), true, 0)
}

func (sr *snippetReadWriter) RmDir(uri string) error {
	return sr.file.RmDir(subDir, makePath(uri))
}
