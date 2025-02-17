package validator

import (
	"fmt"
	"os"
	"strings"
)

type Validator interface {
	Validate() error
}

func Validate(validators []Validator) error {
	for _, v := range validators {
		if err := v.Validate(); err != nil {
			return err
		}
	}
	return nil
}

/*
Source Directory path validator
*/
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

/*
Destination Directory path validator
*/
type DestDirValidator struct {
	DirPath string
}

func (s *DestDirValidator) Validate() error {
	if strings.TrimSpace(s.DirPath) == "" {
		return fmt.Errorf("destination directory is not specified")
	}

	dir, err := os.Stat(s.DirPath)
	if err != nil {
		return fmt.Errorf("destination directory '%s' does not exist: %v", s.DirPath, err)
	}

	if !dir.IsDir() {
		return fmt.Errorf("destination directory '%s' is not a directory", s.DirPath)
	}

	return nil
}

/*
File path validator
*/
type FilePathValidator struct {
	FilePath string
}

func (f *FilePathValidator) Validate() error {
	if strings.TrimSpace(f.FilePath) == "" {
		return fmt.Errorf("file path is not specified")
	}

	_, err := os.Stat(f.FilePath)
	if err != nil {
		return fmt.Errorf("file '%s' does not exist: %v", f.FilePath, err)
	}
	return nil
}
