package main

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"strings"
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

func TestConfigStruct(t *testing.T) {
	// test that the config struct is set up to unmarshal a config file
	viper.SetConfigType("yaml")
	viper.ReadConfig(strings.NewReader(confStr))

	var C Config
	viper.Unmarshal(&C)

	command0 := Command{Args: []string{"hostname"}}
	assert.EqualValues(t, C.Commands[0], command0)

	command1 := Command{Args: []string{"lsof", "-b", "-M", "-n", "-l"}, Link: "lsof"}
	assert.EqualValues(t, C.Commands[1], command1)

	assert.EqualValues(t, C.Files[0], File{Name: "/proc/vmstat"})
	assert.EqualValues(t, C.Files[1], File{Name: "/proc/meminfo", Link: "meminfo"})
}
