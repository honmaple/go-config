package config

import (
	"fmt"
	"strings"
	"testing"

	"github.com/honmaple/go-config/source/yamlfile"
	"github.com/imdario/mergo"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

const yamlFile = `
test:
  string: "string"
  string_slice:
	- testSlice0
	- testSlice1
  int: 100
  int_slice:
	- 100
	- 101
  float: 0.618
  bool: true
test1:
  - string: "string0"
  - string: "string1"
test2:
  map0:
	string: map0
  map1:
	string: map1
`

type (
	yamlTestStruct struct {
		String      string   `yaml:"string"`
		StringSlice []string `yaml:"string_slice"`
		Int         int      `yaml:"int"`
		IntSlice    []int    `yaml:"int_slice"`
		Float       float64  `yaml:"float"`
		Bool        bool     `yaml:"bool"`
	}
	yamlStruct struct {
		Test  yamlTestStruct
		Test1 []yamlTestStruct
		Test2 map[string]yamlTestStruct
	}
)

type Test2 struct {
	Name string
	Type int
}

func TestConfig(t *testing.T) {
	content := strings.ReplaceAll(yamlFile, "\t", "    ")
	conf := New(WithEnv("APP_"), WithSources(yamlfile.NewWithReader(strings.NewReader(content))))
	if err := conf.Load(); err != nil {
		fmt.Println(err)
	}

	var data yamlStruct
	assert.Nil(t, yaml.Unmarshal([]byte(content), &data))

	result := map[string]interface{}{
		"test.string": data.Test.String,
		"test.int":    data.Test.Int,
		"test.float":  data.Test.Float,
		"test.bool":   data.Test.Bool,
	}

	for k, v := range result {
		assert.Equal(t, conf.Interface(k), v)
	}
	assert.Equal(t, conf.IntSlice("test.int_slice"), data.Test.IntSlice)
	assert.Equal(t, conf.StringSlice("test.string_slice"), data.Test.StringSlice)
	assert.Equal(t, conf.String("test1.0.string"), data.Test1[0].String)
	assert.Equal(t, conf.String("test1.1.string"), data.Test1[1].String)
}

func Te1stMerge(t *testing.T) {
	d1 := map[string]interface{}{
		"server": map[string]interface{}{
			"host": "127",
			"port": 80,
		},
	}
	d2 := map[string]interface{}{
		"server": map[interface{}]interface{}{
			"host": "128",
		},
	}
	fmt.Println(mergo.Merge(&d1, d2))
	fmt.Println(d1)
}
