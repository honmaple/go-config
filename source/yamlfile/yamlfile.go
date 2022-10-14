package yamlfile

import (
	"fmt"
	"io"
	"os"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

type yamlfile struct {
	reader io.Reader
	files  []string
}

func cleanMap(m interface{}) interface{} {
	switch t := m.(type) {
	case []interface{}:
		tmp := make([]interface{}, len(t))
		for i := range t {
			tmp[i] = cleanMap(t[i])
		}
		return tmp
	case map[interface{}]interface{}:
		tmp := make(map[string]interface{})
		for k, v := range t {
			tmp[fmt.Sprintf("%v", k)] = cleanMap(v)
		}
		return tmp
	}
	return m
}

func (s *yamlfile) Read() (map[string]interface{}, error) {
	var result map[string]interface{}

	for _, file := range s.files {
		fh, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer fh.Close()

		b, err := io.ReadAll(fh)
		if err != nil {
			return nil, err
		}
		if len(b) == 0 {
			continue
		}
		var tmp map[string]interface{}

		if err := yaml.Unmarshal(b, &tmp); err != nil {
			return nil, err
		}
		if err := mergo.Map(&result, tmp); err != nil {
			return nil, err
		}
	}
	if s.reader != nil {
		var tmp map[string]interface{}
		if err := yaml.NewDecoder(s.reader).Decode(&tmp); err != nil {
			return nil, err
		}
		if err := mergo.Map(&result, tmp); err != nil {
			return nil, err
		}
	}
	for k, v := range result {
		result[k] = cleanMap(v)
	}
	return result, nil
}

func New(files ...string) *yamlfile {
	return &yamlfile{files: files}
}

func NewWithReader(r io.Reader) *yamlfile {
	return &yamlfile{reader: r}
}
