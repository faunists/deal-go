package processors

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/faunists/deal-go/entities"
)

// ReadContractFile reads a JSON File and try to parse it to a entities.Contract object
func ReadContractFile(filePath string) (entities.Contract, error) {
	splitFilePath := strings.Split(filePath, ".")
	extension := splitFilePath[len(splitFilePath)-1]

	var unmarshaler func([]byte, interface{}) error
	switch extension {
	case "json":
		unmarshaler = json.Unmarshal
	case "yaml", "yml":
		unmarshaler = yaml.Unmarshal
	default:
		return entities.Contract{}, errors.New(
			"invalid contract extension, supported formats are json and yaml",
		)
	}

	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return entities.Contract{}, err
	}

	rawContract := entities.Contract{}
	if err = unmarshaler(fileData, &rawContract); err != nil {
		return entities.Contract{}, err
	}

	return rawContract, nil
}
