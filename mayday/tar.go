package mayday

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

type Tarable interface {
	Content() io.Reader
	Header() *tar.Header
	Name() string // full path of file in archive
	Link() string // short link to file in archive
}

type Tar struct {
	gzw *gzip.Writer
	tw  *tar.Writer
}

func (t *Tar) Init(w io.Writer) error {
	t.gzw = gzip.NewWriter(w)
	t.tw = tar.NewWriter(t.gzw)
	return nil
}

func (t *Tar) Add(tb Tarable) error {

	// virtual files, like those in /proc, report a size of 0 to stat().
	// this means the header in the tarfile reports a size of 0 for the file.
	// to avoid this, we copy the file into a buffer, and use that to get the
	// number of bytes to copy.

	buf := new(bytes.Buffer)
	buf.ReadFrom(tb.Content())
	header := tb.Header()
	header.Size = int64(buf.Len())
	header.Name = strings.TrimPrefix(header.Name, "/")

	var err error

	if err = t.tw.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(t.tw, buf)

	if err != nil {
		return fmt.Errorf("could not copy file: %v", err)
	}
	return nil
}

func (t *Tar) MaybeMakeLink(src string, dst string) error {
	if src == "" {
		return nil
	}

	var header tar.Header
	header.Name = src
	header.Linkname = strings.TrimPrefix(dst, "/")
	header.Typeflag = tar.TypeSymlink
	header.ModTime = time.Now()

	log.Printf("Creating link: %q -> %q", src, dst)
	if err := t.tw.WriteHeader(&header); err != nil {
		return err
	}

	return nil
}

func (t *Tar) Close() error {
	t.tw.Flush()
	t.gzw.Flush()

	if err := t.tw.Close(); err != nil {
		log.Fatalf("error closing zipfile: %v", err)
	}
	if err := t.gzw.Close(); err != nil {
		log.Fatalf("error closing zipfile: %v", err)
	}
	return nil
}
