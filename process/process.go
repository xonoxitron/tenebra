package process

import (
	"sync"

	"github.com/xonoxitron/tenebra/download"

	"github.com/sirupsen/logrus"
)

// ProcessItems downloads and processes items concurrently
func ProcessItems(items []download.Item, outputDir string) {
	var wg sync.WaitGroup

	// Loop through objects and start a goroutine for each download
	for _, item := range items {
		wg.Add(1)
		go func(item download.Item) {
			defer wg.Done()
			// Download and process the item
			DownloadAndProcessItem(item, outputDir)
		}(item)
	}

	// Wait for all downloads to complete
	wg.Wait()
}

// DownloadAndProcessItem downloads and processes a single item
func DownloadAndProcessItem(item download.Item, outputDir string) {
	// Download and unzip the file
	err := download.DownloadAndUnzip(item.URL, outputDir)
	if err != nil {
		logrus.Errorf("Error processing %s: %v", item.Name, err)
	}
}
