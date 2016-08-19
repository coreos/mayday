package mayday

import ()

func Run(t Tar, tarables []Tarable) error {

	for _, tb := range tarables {
		t.Add(tb)
		t.MaybeMakeLink(tb.Link(), tb.Name())
	}

	return nil
}
