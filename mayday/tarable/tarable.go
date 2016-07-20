package tarable

import (
	"archive/tar"
	"bytes"
	"time"
)

type Tarable interface {
	Content() *bytes.Buffer
	Header() *tar.Header
	Name() string // full path of file in archive
	Link() string // short link to file in archive
}

// the default implementation of Header()
func Header(content *bytes.Buffer, name string) *tar.Header {

	header := new(tar.Header)
	header.Name = name
	header.Mode = 0666

	header.Size = int64(content.Len())
	header.ModTime = time.Now()

	return header
}
