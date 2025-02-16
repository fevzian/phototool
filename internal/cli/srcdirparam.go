package cli

import (
	"fmt"
	"os"
	"strings"
)

type SrcDirValidator struct {
	DirPath string
}

func (s *SrcDirValidator) Validate() error {
	if strings.TrimSpace(s.DirPath) == "" {
		return fmt.Errorf("source directory is not specified")
	}

	dir, err := os.Stat(s.DirPath)
	if err != nil {
		return fmt.Errorf("source directory '%s' does not exist: %v", s.DirPath, err)
	}

	if !dir.IsDir() {
		return fmt.Errorf("source directory '%s' is not a directory", s.DirPath)
	}

	return nil
}
