package mayday

import (
	"archive/tar"
	"io"
	"log"
)

type MaydayFile struct {
	name    string      // the name of the file on the filesystem
	content io.Reader   // a Reader containing the contents of the file
	header  *tar.Header // a Header containing file path, file size, etc
	link    string      // a link to make in the root of the tarball
}

func NewFile(c io.Reader, h *tar.Header, n string, l string) MaydayFile {
	f := new(MaydayFile)
	f.name = n
	f.content = c
	f.header = h
	f.link = l
	return *f
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
