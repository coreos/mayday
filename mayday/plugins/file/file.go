package file

import (
	"archive/tar"
	"bytes"
	"io"
	"log"
)

type MaydayFile struct {
	name    string        // the name of the file on the filesystem
	file    io.Reader     // the file handler. Not populated until read
	header  *tar.Header   // a Header containing file path, file size, etc
	content *bytes.Buffer // contents of the file, copied in by .Content()
	link    string        // a link to make in the root of the tarball
}

func New(c io.Reader, h *tar.Header, n string, l string) *MaydayFile {
	f := new(MaydayFile)
	f.name = n
	f.header = h
	f.link = l
	f.file = c

	return f
}

func (f MaydayFile) Content() *bytes.Buffer {
	log.Printf("Collecting file: %q\n", f.name)
	f.content = new(bytes.Buffer)
	f.content.ReadFrom(f.content)
	return f.content
}

func (f MaydayFile) Header() *tar.Header {
	if f.content == nil {
		f.Content()
	}
	f.header.Size = int64(f.content.Len())
	return f.header
}

func (f MaydayFile) Name() string {
	return f.name
}

func (f MaydayFile) Link() string {
	return f.link
}
