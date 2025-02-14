package cli

import (
	"fmt"
	"os"
)

type DestDirValidator struct {
	DirPath string
}

func (s *DestDirValidator) Validate() error {
	dir, err := os.Stat(s.DirPath)
	if err != nil {
		return fmt.Errorf("destination directory '%s' does not exist: %v", s.DirPath, err)
	}

	if !dir.IsDir() {
		return fmt.Errorf("destination directory '%s' is not a directory", s.DirPath)
	}

	return nil
}
