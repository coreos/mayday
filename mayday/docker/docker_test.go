package docker

import (
	"encoding/json"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

const (
	dcString            = `{"Config": {"Safe": "abc", "Env": ["POSTGRES_PASSWORD=mysecretpassword", "FOO=bar"]}}`
	dcStringScrubbedEnv = `{"Config": {"Safe": "abc", "Env": [
			"POSTGRES_PASSWORD=scrubbed by mayday",
			"FOO=scrubbed by mayday"]}}`
	dcUuid = "do-re-mi-abc-123"
)

func TestStruct(t *testing.T) {
	dc := New(strings.NewReader(dcString), dcUuid)

	assert.Equal(t, dc.Name(), dcUuid)
	assert.Equal(t, dc.Link(), "") // docker containers don't have links

	assert.EqualValues(t, dc.file, strings.NewReader(dcString))
}

func TestHeader(t *testing.T) {
	dc := New(strings.NewReader(dcString), dcUuid)
	h := dc.Header()

	assert.Equal(t, h.Name, "/docker/"+dcUuid)
}

func TestRepeatedContentCalling(t *testing.T) {
	// Content() should return the same pointer on multiple calls -- there's
	// no need to reparse json, re-read files, etc.
	dc := New(strings.NewReader(dcString), dcUuid)
	call1 := dc.Content()
	call2 := dc.Content()
	assert.Equal(t, &call1, &call2)
}

func TestGetLogsDanger(t *testing.T) {
	viper.Set("danger", true)
	defer viper.Set("danger", false)
	var containers []*DockerContainer
	c1 := New(strings.NewReader(dcString), dcUuid)
	c2 := New(strings.NewReader("content2"), "xyz")
	containers = append(containers, &c1)
	containers = append(containers, &c2)

	logs := getLogs(containers)
	assert.Equal(t, len(logs), 2)

	assert.Equal(t, logs[0].Output, "/docker/do-re-mi-abc-123.log")
	assert.Equal(t, logs[1].Output, "/docker/xyz.log")

	assert.Equal(t, logs[0].Args(), []string{"docker", "logs", "do-re-mi-abc-123"})
	assert.Equal(t, logs[1].Args(), []string{"docker", "logs", "xyz"})
}

func TestGetLogsSafe(t *testing.T) {
	viper.Set("danger", false)
	var containers []*DockerContainer
	c1 := New(strings.NewReader(dcString), dcUuid)
	c2 := New(strings.NewReader("content2"), "xyz")
	containers = append(containers, &c1)
	containers = append(containers, &c2)

	logs := getLogs(containers)
	assert.Equal(t, len(logs), 0)

}

func TestContentSafeMode(t *testing.T) {
	dc := New(strings.NewReader(dcString), dcUuid)
	c := dc.Content()

	var cParsed interface{}
	var dcParsed interface{}

	cbytes, _ := ioutil.ReadAll(c)
	dcbytes := []byte(dcStringScrubbedEnv)

	json.Unmarshal(cbytes, &cParsed)
	// after passing through dc.Content(), the env variables should be scrubbed
	// (as the --danger flag has not been set)
	json.Unmarshal(dcbytes, &dcParsed)
	assert.Equal(t, cParsed, dcParsed)
}

func TestContentDangerMode(t *testing.T) {
	viper.Set("danger", true)
	defer viper.Set("danger", false)
	dc := New(strings.NewReader(dcString), dcUuid)
	c := dc.Content()

	var cParsed map[string]interface{}
	var dcParsed map[string]interface{}

	json.NewDecoder(c).Decode(cParsed)
	// after passing through dc.Content(), the env variables should NOT be scrubbed
	// (as the --danger flag has been set)
	json.Unmarshal([]byte(dcString), dcParsed)

	assert.EqualValues(t, cParsed, dcParsed)
}
