package tarable

import (
	"archive/tar"
	"bytes"
	"io"
	"time"
)

type Tarable interface {
	Content() io.Reader
	Header() *tar.Header
	Name() string // full path of file in archive
	Link() string // short link to file in archive
}

// the default implementation of Header()
func Header(content io.Reader, name string) *tar.Header {

	header := new(tar.Header)
	header.Name = name
	header.Mode = 0666

	// virtual files, like those in /proc, report a size of 0 to stat().
	// this means the header in the tarfile reports a size of 0 for the file.
	// to avoid this, we copy the file into a buffer, and use that to get the
	// number of bytes to copy.

	buf := new(bytes.Buffer)
	buf.ReadFrom(content)
	header.Size = int64(buf.Len())

	header.ModTime = time.Now()

	return header
}
