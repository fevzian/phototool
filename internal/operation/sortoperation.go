package operation

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fevzian/phototool/internal/cli"
	"github.com/rwcarlsen/goexif/exif"
)

var EXTENSIONS = map[string]bool{
	".JPG":  true,
	".JPEG": true,
	".PNG":  true,
	".GIF":  true,
	".TIFF": true,
	".BMP":  true,
	".HEIC": true,
	".MP4":  true,
	".MOV":  true,
	".AVI":  true,
	".MKV":  true,
	".WEBM": true,
	".FLV":  true,
	".WMV":  true,
	".MPG":  true,
	".MPEG": true,
	".3GP":  true,
	".3G2":  true,
}

type SortParams struct {
	SrcDir  string
	DestDir string
}

type SortOperation struct {
	SourceDir string
	DestDir   string
}

func (s *SortOperation) Execute() error {
	params := &SortParams{
		SrcDir:  s.SourceDir,
		DestDir: s.DestDir,
	}
	err := validate(params)
	if err != nil {
		log.Printf("Validation failed: %v", err)
		return err
	}
	return process(params)
}

func (s *SortOperation) GetName() string {
	return "sort"
}

func validate(cmdParams *SortParams) error {
	return cli.Validate([]cli.Validator{
		&cli.SrcDirValidator{DirPath: cmdParams.SrcDir},
		&cli.DestDirValidator{DirPath: cmdParams.DestDir},
	})
}

func process(cmdParams *SortParams) error {
	log.Println("Processing...", cmdParams)
	var count = 0
	err := filepath.Walk(cmdParams.SrcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fileExt := strings.ToUpper(filepath.Ext(info.Name()))
		if !EXTENSIONS[fileExt] {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			log.Printf("Error opening file: %v\n", err)
			return nil
		}
		defer file.Close()

		year, month, err := fetchYearMonth(file)
		if err != nil {
			log.Printf("Error parsing date for file %s: %v\n", path, err)
			return nil
		}

		// Create the destination directory if it does not exist
		destDir := filepath.Join(cmdParams.DestDir, strconv.Itoa(year), month)
		err = os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			log.Printf("Error creating directory %s: %v\n", destDir, err)
			return nil
		}

		// Move the file to the destination directory
		destPath := filepath.Join(destDir, filepath.Base(path))
		// Check if the file with the same name exists in the destination directory
		if _, err := os.Stat(destPath); err == nil {
			// File exists, add a suffix to the file name
			destPath = addSuffix(destPath, 1)
		}

		err = os.Rename(path, destPath)
		if err != nil {
			log.Printf("Error moving file %s to %s: %v\n", path, destPath, err)
			return nil
		}
		count++
		log.Println("Moved file", path, "to", destPath)

		return nil
	})

	if err != nil {
		log.Fatalf("Error walking the path %s: %v", cmdParams.SrcDir, err)
	}

	log.Println("Processed", count, "files")

	return nil
}

func fetchYearMonth(file *os.File) (year int, month string, err error) {
	exifData, err := exif.Decode(file)
	if err != nil {
		log.Printf("Error decoding exif data for file '%s': %v", file.Name(), err)
		return 0, "", err
	}

	dateTaken, err := exifData.Get(exif.DateTimeOriginal)
	if err != nil {
		log.Printf("Error getting '%s' tag: %v", exif.DateTimeOriginal, err)
		return 0, "", err
	}

	dTimeVal, err := dateTaken.StringVal()
	if err != nil {
		log.Printf("Error getting string value for '%s' tag: %v", exif.DateTimeOriginal, err)
		return 0, "", err
	}

	defer file.Close()
	return parseTime(dTimeVal)
}

func addSuffix(filePath string, count int) string {
	ext := filepath.Ext(filePath)
	name := strings.TrimSuffix(filePath, ext)
	newPath := fmt.Sprintf("%s_%d%s", name, count, ext)
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		return newPath
	}
	return addSuffix(filePath, count+1)
}

func parseTime(dateTime string) (year int, month string, err error) {
	parsedDateTime, err := time.Parse("2006:01:02 15:04:05", dateTime)
	if err != nil {
		return 0, "", err
	}
	return parsedDateTime.Year(), parsedDateTime.Month().String()[:3], nil
}
