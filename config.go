package config

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/honmaple/go-config/encoder"
	"github.com/honmaple/go-config/encoder/yaml"
	"github.com/honmaple/go-config/source"
	"github.com/honmaple/go-config/source/env"
	"github.com/honmaple/go-config/source/jsonfile"
	"github.com/honmaple/go-config/source/yamlfile"
	"github.com/imdario/mergo"
)

type (
	Config interface {
		Load() error
		Get(string) interface{}
		Set(string, interface{})
		Sub(string) Config
		Data() map[string]interface{}
		SetData(map[string]interface{})

		Scan(string, interface{}) error
		Bytes(string, ...[]byte) ([]byte, error)
		Int(string, ...int) int
		IntSlice(string, ...[]int) []int
		Float64(string, ...float64) float64
		Bool(string, ...bool) bool
		String(string, ...string) string
		StringSlice(string, ...[]string) []string
		StringMap(string, ...map[string]interface{}) map[string]interface{}
		StringMapStringSlice(string, ...map[string][]string) map[string][]string
		Duration(string, ...time.Duration) time.Duration
		Interface(string, ...interface{}) interface{}
	}
	Options struct {
		sources   []source.Source
		encoder   encoder.Encoder
		separator byte
		cache     bool
	}
	Option func(*Options)
)

type config struct {
	opts  *Options
	body  []byte
	data  map[string]interface{}
	cache map[string]interface{}
}

func (s *config) lookup(key string, data interface{}) (interface{}, bool) {
	if key == "" {
		return data, true
	}

	switch v := data.(type) {
	case map[string]interface{}:
		if n, ok := v[key]; ok {
			return n, true
		}
		if key == "#" {
			return len(v), true
		}
		for i := len(key) - 1; i > 0; i-- {
			if key[i] != s.opts.separator {
				continue
			}
			if n, ok := v[key[:i]]; ok {
				return s.lookup(key[i+1:], n)
			}
		}
		return nil, false
	case []interface{}:
		if key == "#" {
			return len(v), true
		}
		i := 0
		for ; i < len(key) && key[i] != s.opts.separator; i++ {
		}
		idx, err := strconv.Atoi(key[:i])
		if err != nil || idx >= len(v) {
			return nil, false
		}
		if i < len(key) {
			return s.lookup(key[i+1:], v[idx])
		}
		return v[idx], true
	case string:
		if key == "#" {
			return len(v), true
		}
		i := 0
		for ; i < len(key) && key[i] != s.opts.separator; i++ {
		}
		idx, err := strconv.Atoi(key[:i])
		if err != nil || idx >= len(v) {
			return nil, false
		}
		if i < len(key) {
			return s.lookup(key[i+1:], v[idx:idx+1])
		}
		return v[idx : idx+1], true
	default:
		return nil, false
	}
}

func (s *config) Load() error {
	for _, so := range s.opts.sources {
		data, err := so.Read()
		if err != nil {
			return err
		}

		if s.data == nil {
			s.data = data
			continue
		}
		if err := mergo.Merge(&s.data, data); err != nil {
			return err
		}
	}
	return nil
}

func (s *config) checkGet(key string) (value interface{}, ok bool) {
	if value, ok = s.data[key]; ok {
		return
	}
	for i := len(key) - 1; i > 0; i-- {
		if key[i] != s.opts.separator {
			continue
		}
		if v, exists := s.data[key[:i]]; exists {
			value, ok = s.lookup(key[i+1:], v)
			if ok {
				return
			}
		}
	}
	return nil, false
}

func (s *config) Get(key string) interface{} {
	if s.opts.cache {
		if value, ok := s.cache[key]; ok {
			return value
		}
	}
	value, ok := s.checkGet(key)
	if !ok {
		return nil
	}
	if s.opts.cache {
		s.cache[key] = value
	}
	return value
}

func (s *config) Set(key string, value interface{}) {
	switch value.(type) {

	}
}

func (s *config) Data() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range s.data {
		keys := strings.Split(k, string(s.opts.separator))
		if len(keys) == 1 {
			rv := result[k]
			if rv == nil {
				result[k] = v
			} else {
				mergo.Map(&result, map[string]interface{}{k: v})
				// result[k] = rv
				// fmt.Println(result[k], v)
			}
			continue
		}
		tmp := make(map[string]interface{})
		for i := len(keys) - 1; i > 0; i-- {
			k := keys[i]
			if i == len(keys)-1 {
				tmp[k] = v
				continue
			}
			tmp = map[string]interface{}{k: tmp}
		}
		k = keys[0]
		rv := result[k]
		if rv == nil {
			result[k] = tmp
			fmt.Println(tmp, k)
		} else {
			mergo.Map(&rv, tmp)
			result[k] = rv
		}
		// mergo.Merge(&rv, tmp)
	}
	return result
}

func (s *config) SetData(data map[string]interface{}) {
	s.data = data
}

func WithEncoder(e encoder.Encoder) Option {
	return func(o *Options) {
		o.encoder = e
	}
}

func WithSources(s ...source.Source) Option {
	return func(o *Options) {
		o.sources = append(o.sources, s...)
	}
}

func WithSeparator(sep byte) Option {
	return func(o *Options) {
		o.separator = sep
	}
}

func WithEnv(prefix string) Option {
	return func(o *Options) {
		o.sources = append(o.sources, env.New(prefix))
	}
}

func WithCache() Option {
	return func(o *Options) {
		o.cache = true
	}
}

func New(opts ...Option) Config {
	o := new(Options)
	o.sources = make([]source.Source, 0)
	for _, opt := range opts {
		opt(o)
	}
	if o.encoder == nil {
		o.encoder = yaml.NewEncoder()
	}
	if len(o.sources) == 0 {
		o.sources = append(o.sources, jsonfile.New("config.json"))
	}
	if o.separator == 0 {
		o.separator = '.'
	}
	conf := &config{opts: o}
	if o.cache {
		conf.cache = make(map[string]interface{})
	}
	return conf
}

func LoadFile(path string) (Config, error) {
	var opt Option
	switch filepath.Ext(path) {
	case ".yaml", ".yml":
		opt = WithSources(yamlfile.New(path))
	case ".json":
		opt = WithSources(jsonfile.New(path))
	}
	conf := New(opt)
	return conf, conf.Load()
}
