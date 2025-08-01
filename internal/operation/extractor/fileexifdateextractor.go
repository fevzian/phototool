package extractor

import (
	"errors"
	"log"
	"os"
	"time"

	exif "github.com/barasher/go-exiftool"
)

type ExifDateTimeExtractor struct {
}

func (ext *ExifDateTimeExtractor) ExtractDateTime(file *os.File) (value *time.Time, err error) {
	tool, err := exif.NewExiftool()
	if err != nil {
		log.Printf("Error decoding exif data for file '%s': %v", file.Name(), err)
		return nil, err
	}
	defer tool.Close()

	metaData := tool.ExtractMetadata(file.Name())
	if metaData == nil {
		log.Printf("No metadata found for file '%s'", file.Name())
		return nil, errors.New("no metadata found")
	}

	dTimeVal, err := metaData[0].GetString("DateTimeOriginal")

	// dTimeVal, err := findDateTimeValueByExifAttributes(exifData, []exif.FieldName{
	// 	exif.DateTimeOriginal,
	// 	exif.DateTime,
	// 	exif.DateTimeDigitized,
	// })

	if err != nil {
		log.Printf("No valid date found in exif data for file '%s'", file.Name())
		return nil, err
	}

	parsed, err := time.Parse("2006:01:02 15:04:05", dTimeVal)
	if err != nil {
		return nil, err
	}
	log.Printf("Parsed date from exif: %v", parsed)
	return &parsed, nil
}

/*
func findDateTimeValueByExifAttributes(exif *exif.Exif, attributes []exif.FieldName) (string, error) {
	log.Println("Searching for date time in attributes...")
	for _, attr := range attributes {
		log.Printf("Attribute: %s", attr)
		dateTimeVal, err := exif.Get(attr)
		if err != nil {
			log.Printf("Error getting '%s' tag: %v", attr, err)
			continue
		}
		if dateTimeVal != nil {
			dTimeVal, err := dateTimeVal.StringVal()
			if err != nil {
				log.Printf("Error getting string value for '%s' tag: %v", attr, err)
				continue
			}
			log.Printf("Found date time value: %s", dTimeVal)
			return dTimeVal, nil
		}
	}
	log.Println("No date time attribute found")
	return "", errors.New("No date time attribute found in EXIF data")
}
*/
