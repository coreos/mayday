package rkt

import (
	"bytes"
	"github.com/coreos/mayday/mayday"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPodTarable(t *testing.T) {
	assert.Implements(t, (*mayday.Tarable)(nil), new(AppTarable))

	pod := Pod{
		Uuid:    "abc123",
		State:   "running",
		Created: "2016-07-12 13:59:10.606 -0700 PDT",
		Started: "2016-07-12 13:59:10.757 -0700 PDT",
		Network: "default:ip4=172.16.28.2",
	}

	pf := NewPodTarable(pod)

	content := new(bytes.Buffer)
	content.ReadFrom(pf.Content())

	res := "STATE\tCREATED\tSTARTED\tNETWORKS\n" +
		"running\t" + pod.Created + "\t" + pod.Started + "\t" + pod.Network + "\n"

	assert.Equal(t, content.String(), res)
}

func TestPodHeader(t *testing.T) {

	pod := Pod{
		Uuid:    "abc123",
		State:   "running",
		Created: "2016-07-12 13:59:10.606 -0700 PDT",
		Started: "2016-07-12 13:59:10.757 -0700 PDT",
		Network: "default:ip4=172.16.28.2",
	}

	pf := NewPodTarable(pod)

	assert.Equal(t, pf.Name(), "abc123")

	hdr := pf.Header()
	assert.Equal(t, hdr.Name, "/containers/abc123/pod")
	assert.Equal(t, hdr.Size, int64(pf.content.Len()))
}
