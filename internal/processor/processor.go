package processor

import (
	"log"

	"github.com/fevzian/photosort/internal/cli"
)

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
	return nil
}
