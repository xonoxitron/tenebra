package main

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/xonoxitron/tenebra/api"
	"github.com/xonoxitron/tenebra/download"
	"github.com/xonoxitron/tenebra/timemachine"
	"github.com/xonoxitron/tenebra/types"
)

const jsonURL = "https://chaos-data.projectdiscovery.io/index.json"
const outputDir = "./output"           // Replace with the desired output directory
const mergedOutputFile = "tenebra.txt" // Name of the merged output file

type Item = types.Item

// CreateOutputFolder checks if the output folder exists, creates it if not
func CreateOutputFolder(outputDir string) error {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		return os.MkdirAll(outputDir, 0755)
	}
	return nil
}

func main() {
	// Check if the output folder exists, create it if not
	if err := CreateOutputFolder(outputDir); err != nil {
		logrus.Fatal("Error creating output folder:", err)
	}

	// Update local JSON and check if it has been updated
	json, process, err := timemachine.FetchJSON()
	if err != nil {
		logrus.Fatal("Error fetching input JSON:", err)
	}

	if process {
		// Create a wait group to wait for all downloads to complete
		var wg sync.WaitGroup

		// Loop through objects and start a goroutine for each download
		for _, item := range json {
			wg.Add(1)
			go func(item Item) {
				defer wg.Done()
				// Download and unzip the file
				err := download.DownloadAndUnzip(item.URL, outputDir)
				if err != nil {
					logrus.Errorf("Error processing %s: %v", item.Name, err)
				}
			}(item)
		}

		// Wait for all downloads to complete
		wg.Wait()

		// Post-processing steps
		if err := download.EnumerateAndRemoveZIPFiles(outputDir); err != nil {
			logrus.Errorf("Error removing .zip files: %v", err)
		}

		if err := download.ParallelCheckAndRemoveEmptyTXTFiles(outputDir); err != nil {
			logrus.Errorf("Error removing empty .txt files: %v", err)
		}

		if err := download.ParallelMergeURLFiles(outputDir, mergedOutputFile); err != nil {
			logrus.Errorf("Error merging URL files: %v", err)
		}

		logrus.Info("Script execution completed.")
	}

	// Start the API
	api.StartAPI()
}
