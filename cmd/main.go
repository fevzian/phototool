package main

import (
	"flag"
	"log"

	"github.com/fevzian/photosort/internal/cli"
	"github.com/fevzian/photosort/internal/processor"
)

func main() {
	var sourceDir = flag.String("src_dir", "", "Source directory")
	var destDir = flag.String("dest_dir", "", "Destination directory")
	var isCopy = flag.Bool("copy", false, "Just copy files without sorting")
	var groupBy = flag.String("sort", "YYYY-MM", "Group by date format specified, default is YYYY-MM")

	flag.Parse()

	var cmdParams = &cli.CmdParams{
		SrcDir:  *sourceDir,
		DestDir: *destDir,
		IsCopy:  *isCopy,
		GroupBy: *groupBy,
	}
	if err := processor.Validate(cmdParams); err != nil {
		log.Fatalf("Error validating: %v", err)
	}

	if err := processor.Process(cmdParams); err != nil {
		log.Fatalf("Error sorting photos: %v", err)
	}
}
