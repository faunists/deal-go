package processors

import (
	"encoding/json"
	"io/ioutil"

	"github.com/faunists/deal/entities"
)

func ReadContractFile(filePath string) (entities.Contract, error) {
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return entities.Contract{}, err
	}

	rawContract := entities.Contract{}
	if err = json.Unmarshal(jsonData, &rawContract); err != nil {
		return entities.Contract{}, err
	}

	return rawContract, nil
}
