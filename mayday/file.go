package mayday

import (
	"archive/tar"
	"io"
)

type File struct {
	Name    string      // the name of the file on the filesystem
	Content io.Reader   // a Reader containing the contents of the file
	Header  *tar.Header // a Header containing file path, file size, etc
	Link    string      // a link to make in the root of the tarball
}
