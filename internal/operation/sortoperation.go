package operation

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fevzian/phototool/internal/operation/extractor"
	"github.com/fevzian/phototool/internal/validator"
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

		ext := filepath.Ext(info.Name())
		_, ok := EXTENSIONS[strings.ToUpper(ext)]
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
		defer log.Println("-----------------------------------")

		log.Printf("Filename: '%s'", file.Name())
		dateTime, err := parseDateTime(file)
		if err != nil {
			log.Printf("Error parsing date for file %s: %v\n", path, err)
			return nil
		}
		year, month := parseYearMonth(dateTime)

		// Create the destination directory if it does not exist
		destDir := filepath.Join(cmdParams.DestDir, strconv.Itoa(year), month)
		err = os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			log.Printf("Error creating directory %s: %v\n", destDir, err)
			return nil
		}

		// Move the file to the destination directory
		//var newFileName string
		//if time is not available
		// if !hasTimePart(dateTime) {
		// 	newFileName = fmt.Sprintf("%s_%s", dateTime.Format("2006-01-02"), strconv.FormatInt(time.Now().UnixMicro(), 10))
		// } else {
		// 	newFileName = dateTime.Format("2006-01-02_150405")
		// }
		destPath := filepath.Join(destDir, fmt.Sprintf("%s%s", dateTime.Format("2006-01-02_150405"), ext))

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

		removeEmptyDirs(cmdParams.SrcDir)
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking the path %s: %v", cmdParams.SrcDir, err)
	}

	log.Println("Processed", count, "files")

	return nil
}

func parseDateTime(file *os.File) (dateTime *time.Time, err error) {
	extractors := make([]extractor.FileDateTimeExtractor, 2)
	extractors[0] = &extractor.ExifDateTimeExtractor{}
	extractors[1] = &extractor.FilenameDateTimeExtractor{}

	for _, ext := range extractors {
		if dateTime, err := ext.ExtractDateTime(file); err == nil {
			return dateTime, nil
		}
	}
	return nil, fmt.Errorf("date could not be extracted from '%s'", file.Name())
}

func parseYearMonth(dateTime *time.Time) (year int, month string) {
	return dateTime.Year(), dateTime.Month().String()[:3]
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

func removeEmptyDirs(root string) error {
	isEmpty, err := isDirEmpty(root)
	if err != nil {
		return err
	}
	if isEmpty {
		err = os.Remove(root)
		if err != nil {
			return err
		}
		log.Println("Removed empty directory", root)
		return nil
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			err = removeEmptyDirs(filepath.Join(root, entry.Name()))
			if err != nil {
				return err
			}
		}
	}

	// Check again if the directory is empty after removing subdirectories
	isEmpty, err = isDirEmpty(root)
	if err != nil {
		return err
	}
	if isEmpty {
		err = os.Remove(root)
		if err != nil {
			return err
		}
		log.Println("Removed empty directory", root)
	}

	return nil
}

func isDirEmpty(dir string) (bool, error) {
	f, err := os.Open(dir)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
