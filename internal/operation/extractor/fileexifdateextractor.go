package extractor

import (
	"log"
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

type ExifDateTimeExtractor struct {
}

func (ext *ExifDateTimeExtractor) ExtractDateTime(file *os.File) (value *time.Time, err error) {
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
	log.Printf("Parsed date from exif: %v", parsed)
	return &parsed, nil
}
