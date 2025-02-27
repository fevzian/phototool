package extractor

import (
	"os"
	"time"
)

type FileDateTimeExtractor interface {
	ExtractDateTime(file *os.File) (value *time.Time, err error)
}
