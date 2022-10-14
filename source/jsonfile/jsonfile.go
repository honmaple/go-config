package jsonfile

import (
	"encoding/json"
	"io"
	"os"

	"github.com/imdario/mergo"
)

type jsonfile struct {
	files []string
}

func (s *jsonfile) Read() (map[string]interface{}, error) {
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

		if err := json.Unmarshal(b, &tmp); err != nil {
			return nil, err
		}
		if err := mergo.Map(&result, tmp); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func New(files ...string) *jsonfile {
	return &jsonfile{files}
}
