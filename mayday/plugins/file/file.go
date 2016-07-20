package file

import (
	"archive/tar"
	"bytes"
	"io"
	"io/ioutil"
	"log"
)

type MaydayFile struct {
	name    string        // the name of the file on the filesystem
	file    io.ReadCloser // the file handler. Not populated until read
	header  *tar.Header   // a Header containing file path, file size, etc
	content *bytes.Buffer // contents of the file, copied in by .Content()
	link    string        // a link to make in the root of the tarball
}

func New(c io.ReadCloser, h *tar.Header, n string, l string) *MaydayFile {
	f := new(MaydayFile)
	f.name = n
	f.header = h
	f.link = l
	f.file = c

	return f
}

func (f *MaydayFile) Content() *bytes.Buffer {
	if f.content == nil {
		log.Printf("Collecting file: %q\n", f.name)
		fbytes, err := ioutil.ReadAll(f.file)
		if err != nil {
			log.Printf("error reading file: %s", err)
		}
		f.content = bytes.NewBuffer(fbytes)
	}
	return f.content
}

func (f *MaydayFile) Header() *tar.Header {
	if f.content == nil {
		f.Content()
	}
	f.header.Size = int64(f.content.Len())
	return f.header
}

func (f *MaydayFile) Name() string {
	return f.name
}

func (f *MaydayFile) Link() string {
	return f.link
}

func (f *MaydayFile) Close() error {
	return f.file.Close()
}
