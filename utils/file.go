package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/xeipuuv/gojsonschema"
)

// Load file into a string
func Load(filename string) (string, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

// LoadJSON loads the file into out. If schema is not "" then it will validate json against schema
func LoadJSON(filename string, out interface{}, schema string) error {
	contents, errLoad := Load(filename)
	if errLoad != nil {
		return errLoad
	}

	if schema != "" {
		if errValidate := validateJSON(contents, schema); errValidate != nil {
			return errValidate
		}
	}

	if errUnmarhsal := json.Unmarshal([]byte(contents), out); errUnmarhsal != nil {
		return errUnmarhsal
	}
	return nil
}

// Save writes content to file
func Save(filename, content string, overwrite bool) error {
	if _, errExists := os.Stat(filename); !overwrite && errExists != nil {
		return errExists
	}
	file, errOpen := os.Create(filename)
	defer file.Close()
	if errOpen != nil {
		return errOpen
	}
	if _, errWrite := file.Write([]byte(content)); errWrite != nil {
		return errWrite
	}
	return nil
}

// SaveJSON writes object as JSON to file
func SaveJSON(filename string, content interface{}, overwrite bool) error {
	outJSON, errJSON := json.MarshalIndent(content, "", "   ")
	if errJSON != nil {
		return errJSON
	}
	if errSave := Save(filename, string(outJSON), overwrite); errSave != nil {
		return errSave
	}
	return nil
}

func validateJSON(contents string, schema string) error {
	doc := gojsonschema.NewStringLoader(contents)
	sch := gojsonschema.NewStringLoader(schema)
	if result, err := gojsonschema.Validate(sch, doc); err != nil {
		return err
	} else if !result.Valid() {
		return fmt.Errorf("JSON file does not conform to provided schema")
	}
	return nil
}
