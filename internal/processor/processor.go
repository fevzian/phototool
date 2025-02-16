package processor

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fevzian/photosort/internal/cli"
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

func Validate(cmdParams *cli.CmdParams) error {
	var validators = []cli.Validator{
		&cli.SrcDirValidator{DirPath: cmdParams.SrcDir},
		&cli.DestDirValidator{DirPath: cmdParams.DestDir},
		&cli.CopyValidator{IsCopy: cmdParams.IsCopy},
		&cli.GroupByValidator{GroupBy: cmdParams.GroupBy},
	}

	for _, v := range validators {
		if err := v.Validate(); err != nil {
			log.Fatalf("Error validating: %v", err)
		}
	}

	return nil
}

func Process(cmdParams *cli.CmdParams) error {
	log.Println("Processing...", cmdParams)

	files, err := os.ReadDir(cmdParams.SrcDir)
	if err != nil {
		log.Fatalf("Error reading source directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileExt := strings.ToUpper(filepath.Ext(file.Name()))
		if !EXTENSIONS[fileExt] {
			continue
		}

		filePath := filepath.Join(cmdParams.SrcDir, file.Name())
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatalf("Error opening file: %v\n", err)
		}

		year, month, err := fetchYearMonth(file)
		if err != nil {
			log.Printf("Error parsing date for file %s: %v\n", filePath, err)
			continue
		}

		// Create the destination directory if it does not exist
		destDir := filepath.Join(cmdParams.DestDir, strconv.Itoa(year), month)
		err = os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			log.Printf("Error creating directory %s: %v\n", destDir, err)
			continue
		}

		// Move the file to the destination directory
		destPath := filepath.Join(destDir, filepath.Base(filePath))
		err = os.Rename(filePath, destPath)
		if err != nil {
			log.Printf("Error moving file %s to %s: %v\n", filePath, destPath, err)
			continue
		}
		log.Println("Moved file", filePath, "to", destPath)
	}
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

func parseTime(dateTime string) (year int, month string, err error) {
	parsedDateTime, err := time.Parse("2006:01:02 15:04:05", dateTime)
	if err != nil {
		return 0, "", err
	}
	return parsedDateTime.Year(), parsedDateTime.Month().String()[:3], nil
}
