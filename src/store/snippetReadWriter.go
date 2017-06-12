package store

import (
	"strings"
)

// NewSnippetReadWriter creates a SnippetReadWriter which is a simpler interface than
// Doc to create and read local documents and is specifically for Snippets.
func NewSnippetReadWriter(file File) SnippetReadWriter {
	return &snippetReadWriter{file: file}
}

type snippetReadWriter struct {
	file File
}

const subDir = "snippets"

func makePath(uri string) string {
	return strings.Replace(strings.TrimPrefix(uri, "/"), "/", "-", -1)
}

// Write writes the given string content to disk and includes it in a 'holding' directory.
// This directory is to allow for the compilation of snippets which require it and the
// subsequent execution of the binary.
func (sr *snippetReadWriter) Write(uri string, content string) (string, error) {
	return sr.file.Write(subDir, makePath(uri), content, true)
}

// Read reads the content represented by the given uri from disk.
func (sr *snippetReadWriter) Read(uri string) (string, error) {
	return sr.file.Read(subDir, makePath(uri), true, 0)
}

// RmDir removes the 'holding' directory for the given uri.
func (sr *snippetReadWriter) RmDir(uri string) error {
	return sr.file.RmDir(subDir, makePath(uri))
}
