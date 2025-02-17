package main

import (
	"flag"
	"log"
	"os"

	"github.com/fevzian/phototool/internal/operation"
)

func main() {
	switch os.Args[1] {
	case "sort":
		params := flag.NewFlagSet("sort", flag.ExitOnError)
		srcDirParam := params.String("src_dir", "", "Source directory")
		destDirParam := params.String("dest_dir", "", "Destination directory")
		params.Parse(os.Args[2:])

		execute(&operation.SortOperation{
			SourceDir: *srcDirParam,
			DestDir:   *destDirParam,
		})
	case "exif":
		params := flag.NewFlagSet("exif", flag.ExitOnError)
		filePathParam := params.String("file", "", "File path")
		params.Parse(os.Args[2:])
		execute(&operation.ReadExifOperation{
			FilePath: *filePathParam,
		})
	default:
		log.Fatalf("Unknown command: %s", os.Args[1])
	}
}

func execute(op operation.Operation) {
	if err := op.Execute(); err != nil {
		log.Fatalf("Operation failed: %v", err)
	}
}
