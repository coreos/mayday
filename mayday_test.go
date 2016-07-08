package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	confStr = `{
  "files": [
    {
      "name": "/proc/vmstat"
    }, {
      "name": "/proc/meminfo",
      "link": "meminfo"
    }
  ],
  "commands": [
    {
      "args": ["hostname"]
    },
    {
      "args": ["lsof", "-b", "-M", "-n", "-l"],
      "link": "lsof"
    }
  ]
}
`
)

func TestConfigParse(t *testing.T) {
	files, commands, err := readConfig(confStr)

	assert.Nil(t, err)

	command0 := Command{Args: []string{"hostname"}}
	assert.EqualValues(t, commands[0], command0)

	command1 := Command{Args: []string{"lsof", "-b", "-M", "-n", "-l"}, Link: "lsof"}
	assert.EqualValues(t, commands[1], command1)

	assert.EqualValues(t, files[0], File{Name: "/proc/vmstat"})
	assert.EqualValues(t, files[1], File{Name: "/proc/meminfo", Link: "meminfo"})
}
