package cli

import (
	"fmt"
	"os"
)

type SrcDirValidator struct {
	DirPath string
}

func (s *SrcDirValidator) Validate() error {
	dir, err := os.Stat(s.DirPath)
	if err != nil {
		return fmt.Errorf("source directory '%s' does not exist: %v", s.DirPath, err)
	}

	if !dir.IsDir() {
		return fmt.Errorf("source directory '%s' is not a directory", s.DirPath)
	}

	return nil
}
