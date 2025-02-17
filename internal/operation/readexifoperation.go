package operation

import (
	"log"
	"os"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

type Walker struct{}

type ReadExifOperation struct {
	FilePath string
}

// Execute method for ReadExifOperation
func (r *ReadExifOperation) Execute() error {
	file, err := os.Open(r.FilePath)
	if err != nil {
		log.Fatalf("failed to open image file: %v", err)
	}
	defer file.Close()

	// Decode EXIF data
	x, err := exif.Decode(file)
	if err != nil {
		log.Fatalf("failed to decode EXIF data: %v", err)
	}

	// Print EXIF data
	x.Walk(&Walker{})

	return nil
}

// GetName method for ReadExifOperation
func (r *ReadExifOperation) GetName() string {
	return "exif"
}

func (w *Walker) Walk(name exif.FieldName, tag *tiff.Tag) error {
	data, _ := tag.MarshalJSON()
	log.Printf("    %v: %v\n", name, string(data))
	return nil
}
