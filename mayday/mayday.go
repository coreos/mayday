package mayday

import (
	"github.com/coreos/mayday/mayday/tar"
	"github.com/coreos/mayday/mayday/tarable"
)

func Run(t tar.Tar, tarables []tarable.Tarable) error {

	for _, tb := range tarables {
		t.Add(tb)
		t.MaybeMakeLink(tb.Link(), tb.Name())
	}

	return nil
}
