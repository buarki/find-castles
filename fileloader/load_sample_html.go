package fileloader

import (
	"fmt"
	"os"
)

func LoadHTMLFile(filePath string) ([]byte, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load HTML file at [%s], %v", filePath, err)
	}
	return b, nil
}
