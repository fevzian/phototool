package operation

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fevzian/phototool/internal/validator"
	"github.com/rwcarlsen/goexif/exif"
)

type nothing struct{}

var nothing_value = nothing{}

var EXTENSIONS = map[string]nothing{
	".JPG":  nothing_value,
	".JPEG": nothing_value,
	".PNG":  nothing_value,
	".GIF":  nothing_value,
	".TIFF": nothing_value,
	".BMP":  nothing_value,
	".HEIC": nothing_value,
	".MP4":  nothing_value,
	".MP":   nothing_value,
	".MOV":  nothing_value,
	".AVI":  nothing_value,
	".MKV":  nothing_value,
	".WEBM": nothing_value,
	".FLV":  nothing_value,
	".WMV":  nothing_value,
	".MPG":  nothing_value,
	".MPEG": nothing_value,
	".3GP":  nothing_value,
	".3G2":  nothing_value,
}

type SortParams struct {
	SrcDir  string
	DestDir string
}

type SortOperation struct {
	SourceDir string
	DestDir   string
}

type FileDateTimeExtractor interface {
	extractDateTime(file *os.File) (value *time.Time, err error)
}

type ExifDateTimeExtractor struct {
}

type FilenameDateTimeExtractor struct {
}

func (extractor *FilenameDateTimeExtractor) extractDateTime(file *os.File) (value *time.Time, err error) {
	ext := filepath.Ext(file.Name())
	name := filepath.Base(strings.TrimSuffix(file.Name(), ext))

	log.Printf("Filename: '%s'", name)
	parts := splitAny(name, "_-")

	for _, part := range parts {
		parsed, err := time.Parse("20060102", part)
		if err == nil {
			return &parsed, nil
		} else {
			log.Printf("Error parsing date from part='%s': %v", part, err)
		}
	}
	return nil, errors.New("not Implemented yet")
}

func (ext *ExifDateTimeExtractor) extractDateTime(file *os.File) (value *time.Time, err error) {
	exifData, err := exif.Decode(file)
	if err != nil {
		log.Printf("Error decoding exif data for file '%s': %v", file.Name(), err)
		return nil, err
	}

	dateTaken, err := exifData.Get(exif.DateTimeOriginal)
	if err != nil {
		log.Printf("Error getting '%s' tag: %v", exif.DateTimeOriginal, err)
		return nil, err
	}

	dTimeVal, err := dateTaken.StringVal()
	if err != nil {
		log.Printf("Error getting string value for '%s' tag: %v", exif.DateTimeOriginal, err)
		return nil, err
	}
	parsed, err := time.Parse("2006:01:02 15:04:05", dTimeVal)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
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
	return validator.Validate([]validator.Validator{
		&validator.SrcDirValidator{DirPath: cmdParams.SrcDir},
		&validator.DestDirValidator{DirPath: cmdParams.DestDir},
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
		_, ok := EXTENSIONS[fileExt]
		if !ok {
			log.Println("Not supported file type: ", info.Name())
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
	extractors := make([]FileDateTimeExtractor, 2)
	extractors[0] = &ExifDateTimeExtractor{}
	extractors[1] = &FilenameDateTimeExtractor{}

	for _, ext := range extractors {
		if parsedTime, err := ext.extractDateTime(file); err == nil {
			return parsedTime.Year(), parsedTime.Month().String()[:3], nil
		}
	}
	return -1, "", fmt.Errorf("date could not be extracted from '%s'", file.Name())
}

func addSuffix(filePath string, count int) string {
	ext := filepath.Ext(filePath)
	filePathWithoutExt := strings.TrimSuffix(filePath, ext)
	newPath := fmt.Sprintf("%s_%d%s", filePathWithoutExt, count, ext)
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		return newPath
	}
	return addSuffix(filePath, count+1)
}

func splitAny(s string, seps string) []string {
	splitter := func(r rune) bool {
		return strings.ContainsRune(seps, r)
	}
	return strings.FieldsFunc(s, splitter)
}
