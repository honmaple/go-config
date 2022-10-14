package env

import (
	"os"

	"strings"

	"github.com/honmaple/go-config/encoder"
	"github.com/imdario/mergo"
	"strconv"
)

const defaultPrefix = "APP_"

type source struct {
	prefix  string
	encoder encoder.Encoder
}

func (s *source) Read() (map[string]interface{}, error) {
	var result map[string]interface{}

	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, s.prefix) {
			continue
		}
		pair := strings.SplitN(strings.TrimPrefix(env, s.prefix), "=", 2)
		keys := strings.Split(strings.ToLower(pair[0]), "_")
		value := pair[1]

		tmp := make(map[string]interface{})
		for i := len(keys) - 1; i >= 0; i-- {
			k := keys[i]
			if i == len(keys)-1 {
				if intValue, err := strconv.Atoi(value); err == nil {
					tmp[k] = intValue
				} else if boolValue, err := strconv.ParseBool(value); err == nil {
					tmp[k] = boolValue
				} else {
					tmp[k] = value
				}
				continue
			}
			tmp = map[string]interface{}{k: tmp}
		}
		if err := mergo.Map(&result, tmp); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func New(prefix string) *source {
	if prefix == "" {
		prefix = defaultPrefix
	}
	return &source{prefix: prefix}
}
