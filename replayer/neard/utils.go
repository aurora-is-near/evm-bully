package neard

import (
	"os"

	"github.com/buger/jsonparser"
)

type jsonEdit struct {
	Path  []string
	Value string
}

func editJSONFile(filename string, edits []*jsonEdit) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	for _, edit := range edits {
		data, err = jsonparser.Set(data, []byte(edit.Value), edit.Path...)
		if err != nil {
			return err
		}
	}
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return err
	}
	return nil
}
