package operation

import (
	"errors"
	"log"

	exif "github.com/barasher/go-exiftool"
)

type Walker struct{}

type ReadExifOperation struct {
	FilePath string
}

// Execute method for ReadExifOperation
func (r *ReadExifOperation) Execute() error {
	et, err := exif.NewExiftool()
	if err != nil {
		log.Printf("Error when intializing: %v\n", err)
		return errors.New("no metadata found")
	}
	defer et.Close()

	fileInfos := et.ExtractMetadata(r.FilePath)

	for _, fileInfo := range fileInfos {
		if fileInfo.Err != nil {
			log.Printf("Error concerning %v: %v\n", fileInfo.File, fileInfo.Err)
			continue
		}

		for k, v := range fileInfo.Fields {
			log.Printf("[%v] %v\n", k, v)
		}
	}
	return nil
}

// GetName method for ReadExifOperation
func (r *ReadExifOperation) GetName() string {
	return "exif"
}
