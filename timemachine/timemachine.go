package timemachine

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/xonoxitron/tenebra/download"
	"github.com/xonoxitron/tenebra/types"
)

const jsonURL = "https://chaos-data.projectdiscovery.io/index.json"
const inputJSONFileName = "input.json"

// FetchJSON checks and updates the input.json file
func FetchJSON() ([]types.Item, bool, error) {
	// Check if input.json exists in the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, false, fmt.Errorf("error getting current directory: %v", err)
	}

	inputJSONPath := filepath.Join(currentDir, inputJSONFileName)
	newJSONPath := filepath.Join(currentDir, "new.json")

	// Check if the local input.json file exists
	_, err = os.Stat(inputJSONPath)
	if err != nil {
		// If the file doesn't exist, download the new input.json
		// Use the common FetchAndParseJSON method to fetch the JSON content
		items, err := download.FetchAndParseJSON(jsonURL, inputJSONPath)
		if err != nil {
			logrus.Fatal("Error downloading input.json:", err)
		}

		logrus.Info("Local input.json created.")
		return items, true, nil
	}

	// download the new new.json
	newItems, err := download.FetchAndParseJSON(jsonURL, newJSONPath)
	if err != nil {
		logrus.Fatal("Error downloading new.json:", err)
	}

	inputJSONContent, err := os.ReadFile(inputJSONPath)
	if err != nil {
		logrus.Fatal("Error reading input.json:", err)
	}
	newJSONContent, err := json.Marshal(newItems)
	if err != nil {
		logrus.Fatal("Error reading new.json:", err)

	}
	// Compare the content of the input.json files
	updated, err := compareFileHashes(inputJSONContent, newJSONContent)
	if err != nil {
		logrus.Fatal("Error comparing hashes:", err)

	}

	println(updated)

	if !updated {
		logrus.Info("Updating input.json.")
		// Replace the local input.json with the new one
		// Read the content of the input file
		inputContent, err := os.ReadFile(newJSONPath)
		if err != nil {
			logrus.Fatal(err)
		}

		// Write the content to the output file, overwriting its previous content
		err = os.WriteFile(inputJSONPath, inputContent, 0644)
		if err != nil {
			logrus.Fatal(err)
		}

		// Optional: Print a success message
		println("File overwritten successfully!")

		// Remove the "new.json" file
		err = os.Remove(newJSONPath)
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Info("Local input.json updated.")
		return newItems, true, nil
	}

	logrus.Info("Local input.json is up-to-date.")

	return newItems, false, nil
}

// calculateFileHash calculates the SHA256 hash of a file's content
func calculateFileHash(content []byte) (string, error) {

	hash := sha256.New()
	hash.Write(content)

	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	println(hashString)
	return hashString, nil
}

// compareFileHashes compares the SHA256 hashes of two files
func compareFileHashes(file1, file2 []byte) (bool, error) {
	hash1, err := calculateFileHash(file1)
	if err != nil {
		return false, err
	}

	hash2, err := calculateFileHash(file2)
	if err != nil {
		return false, err
	}

	return hash1 == hash2, nil
}
