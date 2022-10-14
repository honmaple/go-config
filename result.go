package config

import (
	// "fmt"
	"reflect"
	"time"

	"github.com/spf13/cast"
)

var (
	DefaultConfig = New()
)

func Sub(key string) Config                            { return DefaultConfig.Sub(key) }
func Set(key string, value interface{})                { DefaultConfig.Set(key, value) }
func Get(key string) interface{}                       { return DefaultConfig.Get(key) }
func Load() error                                      { return DefaultConfig.Load() }
func Scan(key string, v interface{}) error             { return DefaultConfig.Scan(key, v) }
func Bytes(key string, def ...[]byte) ([]byte, error)  { return DefaultConfig.Bytes(key, def...) }
func Int(key string, def ...int) int                   { return DefaultConfig.Int(key, def...) }
func IntSlice(key string, def ...[]int) []int          { return DefaultConfig.IntSlice(key, def...) }
func Float64(key string, def ...float64) float64       { return DefaultConfig.Float64(key, def...) }
func Bool(key string, def ...bool) bool                { return DefaultConfig.Bool(key, def...) }
func String(key string, def ...string) string          { return DefaultConfig.String(key, def...) }
func StringSlice(key string, def ...[]string) []string { return DefaultConfig.StringSlice(key, def...) }
func StringMap(key string, def ...map[string]interface{}) map[string]interface{} {
	return DefaultConfig.StringMap(key, def...)
}
func StringMapStringSlice(key string, def ...map[string][]string) map[string][]string {
	return DefaultConfig.StringMapStringSlice(key, def...)
}
func Interface(key string, def ...interface{}) interface{} {
	return DefaultConfig.Interface(key, def...)
}
func Duration(key string, def ...time.Duration) time.Duration {
	return DefaultConfig.Duration(key, def...)
}

func (s *config) Sub(key string) Config {
	data := s.Get(key)
	if data == nil {
		return nil
	}
	if reflect.TypeOf(data).Kind() == reflect.Map {
		subc := New()
		subc.SetData(cast.ToStringMap(data))
		return subc
	}
	return nil
}

func (s *config) Int(key string, def ...int) (result int) {
	if len(def) > 0 {
		result = def[0]
	}
	value := s.Get(key)
	if value == nil {
		return
	}
	d, err := cast.ToIntE(value)
	if err != nil {
		return
	}
	return d
}

func (s *config) IntSlice(key string, def ...[]int) (result []int) {
	if len(def) > 0 {
		result = def[0]
	}
	value := s.Get(key)
	if value == nil {
		return
	}
	d, err := cast.ToIntSliceE(value)
	if err != nil {
		return
	}
	return d
}

func (s *config) Duration(key string, def ...time.Duration) (result time.Duration) {
	if len(def) > 0 {
		result = def[0]
	}
	value := s.Get(key)
	if value == nil {
		return
	}
	d, err := cast.ToDurationE(value)
	if err != nil {
		return
	}
	return d
}

func (s *config) Float64(key string, def ...float64) (result float64) {
	if len(def) > 0 {
		result = def[0]
	}
	value := s.Get(key)
	if value == nil {
		return
	}
	d, err := cast.ToFloat64E(value)
	if err != nil {
		return
	}
	return d
}

func (s *config) Bool(key string, def ...bool) (result bool) {
	if len(def) > 0 {
		result = def[0]
	}
	value := s.Get(key)
	if value == nil {
		return
	}
	d, err := cast.ToBoolE(value)
	if err != nil {
		return
	}
	return d
}

func (s *config) String(key string, def ...string) (result string) {
	if len(def) > 0 {
		result = def[0]
	}
	value := s.Get(key)
	if value == nil {
		return
	}
	d, err := cast.ToStringE(value)
	if err != nil {
		return
	}
	return d
}

func (s *config) StringSlice(key string, def ...[]string) (result []string) {
	if len(def) > 0 {
		result = def[0]
	}
	value := s.Get(key)
	if value == nil {
		return
	}
	d, err := cast.ToStringSliceE(value)
	if err != nil {
		return
	}
	return d
}

func (s *config) StringMap(key string, def ...map[string]interface{}) (result map[string]interface{}) {
	if len(def) > 0 {
		result = def[0]
	}
	value := s.Get(key)
	if value == nil {
		return
	}
	d, err := cast.ToStringMapE(value)
	if err != nil {
		return
	}
	return d
}

func (s *config) StringMapStringSlice(key string, def ...map[string][]string) (result map[string][]string) {
	if len(def) > 0 {
		result = def[0]
	}
	value := s.Get(key)
	if value == nil {
		return
	}
	d, err := cast.ToStringMapStringSliceE(value)
	if err != nil {
		return
	}
	return d
}

func (s *config) Interface(key string, def ...interface{}) (result interface{}) {
	if len(def) > 0 {
		result = def[0]
	}
	value := s.Get(key)
	if value == nil {
		return
	}
	return value
}

func (s *config) Bytes(key string, def ...[]byte) (result []byte, err error) {
	if len(def) > 0 {
		result = def[0]
	}
	value := s.Get(key)
	if value == nil {
		return
	}
	return s.opts.encoder.Encode(value)
}

func (s *config) Scan(key string, v interface{}) error {
	value := s.data
	if key != "" {
		// v := s.Get(key)
		// if !ok {
		//	return fmt.Errorf("The key '%s' is not exists", key)
		// }
		// value = v
	}
	b, err := s.opts.encoder.Encode(value)
	if err != nil {
		return err
	}
	return s.opts.encoder.Decode(b, v)
}
