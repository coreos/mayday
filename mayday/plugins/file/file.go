package file

import (
	"archive/tar"
	"bytes"
	"io"
	"log"
)

type MaydayFile struct {
	name    string        // the name of the file on the filesystem
	content *bytes.Buffer // an in-memory copy of the file, populated by .Content()
	header  *tar.Header   // a Header containing file path, file size, etc
	link    string        // a link to make in the root of the tarball
}

func New(c io.Reader, h *tar.Header, n string, l string) *MaydayFile {
	f := new(MaydayFile)
	f.name = n
	f.header = h
	f.link = l

	buf := new(bytes.Buffer)
	buf.ReadFrom(c)

	f.header.Size = int64(buf.Len())

	f.content = buf

	return f
}

func (f MaydayFile) Content() io.Reader {
	log.Printf("Collecting file: %q\n", f.name)
	return f.content
}

func (f MaydayFile) Header() *tar.Header {
	return f.header
}

func (f MaydayFile) Name() string {
	return f.name
}

func (f MaydayFile) Link() string {
	return f.link
}
