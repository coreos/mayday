package mayday_test

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"strings"
	"testing"

	"github.com/coreos/mayday/mayday"
	"github.com/stretchr/testify/assert"
)

type TestTarable struct{}

func (tt *TestTarable) Run() error { return nil }

func (tt *TestTarable) Content() io.Reader {
	return strings.NewReader("test_content")
}

func (tt *TestTarable) Header() *tar.Header {
	var h tar.Header
	h.Typeflag = tar.TypeReg
	h.Name = "test"
	return &h
}

func (tt *TestTarable) Name() string { return "test_name" }
func (tt *TestTarable) Link() string { return "" }

func TestContentsAdded(t *testing.T) {
	buf := new(bytes.Buffer)
	var tf mayday.Tar
	tf.Init(buf)

	var testtar *TestTarable

	err := tf.Add(testtar)
	assert.Nil(t, err)
	tf.Close()

	gr, err := gzip.NewReader(buf)
	tr := tar.NewReader(gr)

	newbuf := new(bytes.Buffer)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		assert.Equal(t, int(hdr.Typeflag), int(tar.TypeReg))
		newbuf.ReadFrom(tr)
	}
	assert.Equal(t, newbuf.String(), "test_content")
}

func TestLinkAdded(t *testing.T) {
	buf := new(bytes.Buffer)
	var tf mayday.Tar
	tf.Init(buf)

	tf.MaybeMakeLink("short_path", "annoyingly/long/nested/path")
	tf.Close()

	gr, err := gzip.NewReader(buf)
	if err != nil {
		panic(err)
	}

	tr := tar.NewReader(gr)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		assert.Equal(t, int(hdr.Typeflag), int(tar.TypeSymlink))
		assert.Equal(t, hdr.Name, "short_path")
		assert.Equal(t, hdr.Linkname, "annoyingly/long/nested/path")
	}

}
