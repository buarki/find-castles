package fileloader

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/buarki/find-castles/castle"
)

func LoadCastlesAsJSONList(filePath string) ([]castle.Model, error) {
	var castles []castle.Model
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to castles JSON file at [%s], got %v", filePath, err)
	}
	err = json.Unmarshal(b, &castles)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal castle JSON file [%s], got %v", filePath, err)
	}
	return castles, nil
}
