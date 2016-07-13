package rkt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	multi_app_output = `UUID					APP	IMAGE NAME					IMAGE ID		STATE	CREATED					STARTED					NETWORKS
3edcaaea-acb4-497f-85fa-55014a780568	etcd	coreos.com/etcd:v2.3.7				sha512-7d28419b27d5	running	2016-07-12 13:59:10.606 -0700 PDT	2016-07-12 13:59:10.757 -0700 PDT	default:ip4=172.16.28.2
					nginx	registry-1.docker.io/library/nginx:latest	sha512-a119eb68f973
f62c02b1-a514-4f92-b965-610eea629a60	etcd	coreos.com/etcd:v2.3.7				sha512-7d28419b27d5	exited	2016-07-12 13:10:39.777 -0700 PDT	2016-07-12 13:10:39.871 -0700 PDT	default:ip4=172.16.28.5
`

	empty = `UUID					APP	IMAGE NAME					IMAGE ID		STATE	CREATED					STARTED					NETWORKS
`
)

func TestProcessMultiAppRkt(t *testing.T) {
	pods := ProcessRktOutput(multi_app_output)

	assert.Equal(t, len(pods), 2)

	assert.Equal(t, pods[0].Uuid, "3edcaaea-acb4-497f-85fa-55014a780568")
	assert.Equal(t, pods[0].State, "running")
	assert.Equal(t, pods[0].Network, "default:ip4=172.16.28.2")
	assert.Equal(t, len(pods[0].Apps), 2)

	assert.Equal(t, pods[0].Apps[0].Name, "etcd")
	assert.Equal(t, pods[0].Apps[0].ImageName, "coreos.com/etcd:v2.3.7")
	assert.Equal(t, pods[0].Apps[0].ImageId, "sha512-7d28419b27d5")

	assert.Equal(t, pods[0].Apps[1].Name, "nginx")
	assert.Equal(t, pods[0].Apps[1].ImageName, "registry-1.docker.io/library/nginx:latest")
	assert.Equal(t, pods[0].Apps[1].ImageId, "sha512-a119eb68f973")

	assert.Equal(t, pods[1].Uuid, "f62c02b1-a514-4f92-b965-610eea629a60")
	assert.Equal(t, pods[1].State, "exited")
	assert.Equal(t, pods[1].Network, "default:ip4=172.16.28.5")
	assert.Equal(t, len(pods[1].Apps), 1)
}

func TestProcessEmptyRkt(t *testing.T) {
	pods := ProcessRktOutput(empty)
	assert.Equal(t, len(pods), 0)
}
