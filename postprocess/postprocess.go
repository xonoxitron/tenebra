package postprocess

import (
	"github.com/xonoxitron/tenebra/download"

	"github.com/sirupsen/logrus"
)

const mergedOutputFile = "tenebra.txt" // Name of the merged output file

// PerformPostProcess performs post-processing steps
func PerformPostProcess(dir, outputFile string) {
	// Remove .zip files
	if err := download.EnumerateAndRemoveZIPFiles(dir); err != nil {
		logrus.Errorf("Error removing .zip files: %v", err)
	}

	// Remove empty .txt files
	if err := download.ParallelCheckAndRemoveEmptyTXTFiles(dir); err != nil {
		logrus.Errorf("Error removing empty .txt files: %v", err)
	}

	// Merge URL files
	if err := download.ParallelMergeURLFiles(dir, outputFile); err != nil {
		logrus.Errorf("Error merging URL files: %v", err)
	}
}
