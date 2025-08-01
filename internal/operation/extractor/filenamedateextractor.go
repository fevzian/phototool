package extractor

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	DATETIME_FORMAT = "20060102_150405"
	DATE_FORMAT     = "20060102"
)

type FilenameDateTimeExtractor struct {
}

type dateTimeFinder interface {
	findDateTime(words []string) (value *time.Time, err error)
}

type backwardsDateTimeFinder struct {
}

type backForwardDateTimeFinder struct {
}

type singleWordDateTimeFinder struct{}

var finders = []dateTimeFinder{
	&backwardsDateTimeFinder{},
	&backForwardDateTimeFinder{},
	&singleWordDateTimeFinder{},
}

func (extractor *FilenameDateTimeExtractor) ExtractDateTime(file *os.File) (value *time.Time, err error) {
	ext := filepath.Ext(file.Name())
	name := filepath.Base(strings.TrimSuffix(file.Name(), ext))
	parts := splitAny(name, "_-")

	for _, finder := range finders {
		dateTime, err := finder.findDateTime(parts)
		if err == nil {
			return dateTime, nil
		}
	}
	return nil, fmt.Errorf("datetime could not be extracted from filename: '%s'", file.Name())
}

// Tries to parse each array string using DATE_FORMAT format
func (finder *singleWordDateTimeFinder) findDateTime(words []string) (value *time.Time, err error) {
	for _, word := range words {
		parsed, err := parseDateTime(word, DATE_FORMAT)
		if err == nil {
			return parsed, nil
		}
	}
	return nil, errors.New("datetime could not be extracted from words")
}

// Tries to parse a concatanation of array string elements joined by '_' char in the following way:
// It cuts the right and left ends of the array one by one until a single word is left
func (finder *backForwardDateTimeFinder) findDateTime(words []string) (value *time.Time, err error) {
	head := 0
	tail := len(words)

	backwards := true
	for head < tail {
		value := strings.Join(words[head:tail], "_")

		lastPart := tail-head == 1
		format := DATETIME_FORMAT
		if lastPart {
			format = DATE_FORMAT
		}

		parsed, err := parseDateTime(value, format)
		if err == nil {
			return parsed, nil
		}

		if backwards {
			tail--
		} else {
			head++
		}
		backwards = !backwards
	}
	return nil, errors.New("datetime could not be extracted from words")
}

// Tries to parse a concatenation of array string elements joined by '_' char in the following way:
// [0:N], [0:N-1], [0:N-2], ...., [0:1]
func (finder *backwardsDateTimeFinder) findDateTime(words []string) (value *time.Time, err error) {
	tail := len(words)

	for tail > 0 {
		value := strings.Join(words[0:tail], "_")
		parsed, err := parseDateTime(value, DATETIME_FORMAT)
		if err == nil {
			return parsed, nil
		}
		tail--
	}
	return nil, errors.New("datetime could not be extracted from words")
}

func parseDateTime(word string, dateTimeFormat string) (value *time.Time, err error) {
	log.Printf("Trying to parse date from part='%s', format='%s'", word, dateTimeFormat)

	parsed, err := time.Parse(dateTimeFormat, word)
	if err != nil {
		log.Printf("Error parsing date from part='%s', format='%s':%v", word, dateTimeFormat, err)
		return nil, err
	} else {
		log.Printf("Parsed date from  part='%s', format='%s'", word, dateTimeFormat)
		return &parsed, nil
	}
}

func splitAny(s string, seps string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		return strings.ContainsRune(seps, r)
	})
}
