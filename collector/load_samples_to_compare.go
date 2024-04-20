package collector

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/buarki/find-castles/castle"
)

func loadJSONToCompare(filePath string) ([]castle.Model, error) {
	var castles []castle.Model
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open portuguese json, got %v", err)
	}
	err = json.Unmarshal(b, &castles)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON, got %v", err)
	}
	return castles, nil
}
