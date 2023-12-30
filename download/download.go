package download

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xonoxitron/tenebra/types"
)

// Item represents the structure of each object in the JSON array
type Item struct {
	Name        string    `json:"name"`
	ProgramURL  string    `json:"program_url"`
	URL         string    `json:"URL"`
	Count       int       `json:"count"`
	Change      int       `json:"change"`
	IsNew       bool      `json:"is_new"`
	Platform    string    `json:"platform"`
	Bounty      bool      `json:"bounty"`
	LastUpdated time.Time `json:"last_updated"`
}

// FetchAndParseJSON fetches and parses JSON data from the specified URL
func FetchAndParseJSON(url string, filePath string) ([]types.Item, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch JSON from %s: %s", url, resp.Status)
	}

	var items []types.Item // Use the correct type here (not download.Item)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&items)
	if err != nil {
		return nil, err
	}

	// Write the downloaded content to the local file
	jsonContent, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(filePath, jsonContent, 0644); err != nil {
		return nil, err
	}

	return items, nil
}

// DownloadAndUnzip downloads and unzips the file
func DownloadAndUnzip(url, outputDir string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download %s: %s", url, resp.Status)
	}

	// Extract the filename from the URL
	tokens := strings.Split(url, "/")
	filename := tokens[len(tokens)-1]

	// Create the output file
	filePath := filepath.Join(outputDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the downloaded content to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	// Unzip the file
	if err := unzip(filePath, outputDir); err != nil {
		return err
	}

	logrus.Infof("Successfully processed: %s", filename)
	return nil
}

// EnumerateAndRemoveZIPFiles enumerates and removes .zip files in the directory
func EnumerateAndRemoveZIPFiles(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".zip" {
			continue
		}

		filePath := filepath.Join(dir, entry.Name())
		if err := os.Remove(filePath); err != nil {
			return err
		}

		logrus.Infof("Removed .zip file: %s", entry.Name())
	}

	return nil
}

// ParallelCheckAndRemoveEmptyTXTFiles checks and removes empty .txt files in parallel
func ParallelCheckAndRemoveEmptyTXTFiles(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".txt" {
			continue
		}

		wg.Add(1)
		go func(entry fs.DirEntry) {
			defer wg.Done()

			filePath := filepath.Join(dir, entry.Name())

			// Locking to prevent race conditions when accessing shared resources
			mu.Lock()
			defer mu.Unlock()

			if isEmpty, err := isFileEmpty(filePath); err != nil {
				logrus.Errorf("Error checking if file is empty: %v", err)
			} else if isEmpty {
				if err := os.Remove(filePath); err != nil {
					logrus.Errorf("Error removing empty .txt file: %v", err)
				} else {
					logrus.Infof("Removed empty .txt file: %s", entry.Name())
				}
			}
		}(entry)
	}

	// Wait for all checks and removals to complete
	wg.Wait()

	return nil
}

// isFileEmpty checks if a file is empty
func isFileEmpty(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return false, err
	}

	return stat.Size() == 0, nil
}

// ParallelMergeURLFiles merges URL files in parallel
func ParallelMergeURLFiles(dir, outputFile string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var mu sync.Mutex
	var mergedURLs []string
	var wg sync.WaitGroup

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".txt" {
			continue
		}

		wg.Add(1)

		// Process each file in a separate goroutine
		go func(entry fs.DirEntry) {
			defer wg.Done()

			filePath := filepath.Join(dir, entry.Name())

			// Read and process each line in the file
			urls, err := readAndFilterURLs(filePath)
			if err != nil {
				logrus.Errorf("Error reading and filtering URLs from %s: %v", entry.Name(), err)
				return
			}

			// Locking to prevent race conditions when accessing shared resources
			mu.Lock()
			defer mu.Unlock()

			mergedURLs = append(mergedURLs, urls...)
		}(entry)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Write the merged URLs to the output file
	if err := writeMergedURLs(outputFile, mergedURLs); err != nil {
		return err
	}

	logrus.Info("Merged URLs written to:", outputFile)
	return nil
}

// readAndFilterURLs reads and filters URLs from a file
func readAndFilterURLs(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var filteredURLs []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		url := scanner.Text()

		// Skip URLs starting with a wildcard
		if !strings.HasPrefix(url, "*") {
			filteredURLs = append(filteredURLs, url)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return filteredURLs, nil
}

// writeMergedURLs writes merged URLs to the output file
func writeMergedURLs(outputFile string, urls []string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, url := range urls {
		_, err := fmt.Fprintln(file, url)
		if err != nil {
			return err
		}
	}

	return nil
}

// unzip extracts files from a zip archive
func unzip(src, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, file.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if file.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.Create(fpath)
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}
