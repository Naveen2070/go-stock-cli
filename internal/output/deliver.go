package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Naveen2070/go-stock-cli/models"
)

func Deliver(filePath string, selections []models.Selection) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	if err := encoder.Encode(selections); err != nil {
		return fmt.Errorf("error encoding selections: %v", err)
	}
	return nil
}
