package filesystem

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// LoadFromJSONFile is used to load test data from a JSON file into an interface for testing purposes.
func LoadFromJSONFile(filename string, dest any) error {
	file, err := LoadBytesFromJSONFile(filename)
	if err != nil {
		return errors.Wrap(err, "failed to LoadBytes")
	}

	err = json.Unmarshal(file, dest)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal file [%s] to dest", filename)
	}

	return nil
}

func LoadBytesFromJSONFile(jsonfile string) ([]byte, error) {
	abspath, err := filepath.Abs(jsonfile)
	if err != nil {
		return nil, errors.Wrap(err, "could not determine absolute path for file")
	}

	file, err := os.ReadFile(abspath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file [%s]", jsonfile)
	}

	return file, nil
}
