package mayday_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/coreos/mayday/mayday"
	"github.com/stretchr/testify/assert"
)

type tempReadCloser struct {
	io.Reader
}

func (tempReadCloser) Close() error {
	return nil
}

func TestFile(t *testing.T) {

	file := mayday.File{
		Source: tempReadCloser{bytes.NewBufferString("contents1")},
		Path:   os.TempDir() + "/file1",
	}

	res := file.Collect(workspace)
	assert.Nil(t, res)

	f, err := ioutil.ReadFile(workspace + os.TempDir() + "/file1")
	assert.Nil(t, err)
	assert.Equal(t, string(f), "contents1")
}

func TestFileWithLink(t *testing.T) {
	file := mayday.File{
		Source: tempReadCloser{bytes.NewBufferString("contents2")},
		Path:   os.TempDir() + "/annoyingly_long_name",
		Link:   "short",
	}

	res := file.Collect(workspace)
	assert.Nil(t, res)

	f, err := ioutil.ReadFile(workspace + "/short")
	assert.Nil(t, err)
	assert.Equal(t, string(f), "contents2")
}
