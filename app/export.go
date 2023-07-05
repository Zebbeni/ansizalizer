package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Zebbeni/ansizalizer/global"
)

const (
	maxExportJobs = 1000
)

type exportJob struct {
	sourcePath      string
	destinationPath string
}

type MaxExportQueueError struct {
	count int
}

func (r *MaxExportQueueError) Error() string {
	return fmt.Sprintf("%d+ export jobs exceed %d max", r.count, maxExportJobs)
}

// this process may get more complicated if we want to do animated gifs,
// since each gif  will require multiple image exports.
func buildExportQueue(dirPath, destPath string, useSubDirs bool) ([]exportJob, error) {
	// for each image file found in the dirPath, append an exportJob object
	// with the source filepath and its corresponding .ansi destination filepath
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	exportJobs := make([]exportJob, 0, len(entries))
	subDirs := make([]string, 0, len(entries))

	for _, e := range entries {
		sourcePath := filepath.Join(dirPath, e.Name())

		if e.IsDir() {
			subDirs = append(subDirs, sourcePath)
			continue
		}

		ext := filepath.Ext(e.Name())
		if _, ok := global.ImgExtensions[ext]; ok {
			nameWithoutExt := strings.Split(filepath.Base(sourcePath), ".")[0]
			nameWithExt := fmt.Sprintf("%s.ansi", nameWithoutExt)
			destFilePath := filepath.Join(destPath, nameWithExt)
			exportJobs = append(exportJobs, exportJob{
				sourcePath:      sourcePath,
				destinationPath: destFilePath,
			})
		}
	}

	if useSubDirs {
		// call buildExportQueue on each subdirectory in dirPath, creating
		// subdirectories in the destination path to mimic the source directory
		// structure, and providing these subdirectory paths to the build call as well
		for _, subDir := range subDirs {

			subDirName := filepath.Base(subDir)
			subDestPath := filepath.Join(destPath, subDirName)

			var subDirExportJobs []exportJob
			subDirExportJobs, err = buildExportQueue(subDir, subDestPath, true)
			if err != nil {
				return nil, err
			}

			// append resulting exportJob lists to the main list
			exportJobs = append(exportJobs, subDirExportJobs...)
			if len(exportJobs) > maxExportJobs {
				return nil, &MaxExportQueueError{count: len(exportJobs)}
			}

			// skip creating mirrored subdirectories if no files found there
			if len(subDirExportJobs) == 0 {
				continue
			}

			// create the destination folder if it doesn't already exist
			// do this after the recursive call to buildExportQueue. Otherwise,
			// we can hit an infinite loop where our newly created directories
			// get picked up by subsequent buildExportQueue calls, forever.
			if _, err = os.Stat(subDestPath); os.IsNotExist(err) {
				err = os.MkdirAll(subDestPath, os.ModeDir)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return exportJobs, nil
}
