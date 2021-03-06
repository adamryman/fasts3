package util

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"syscall"

	"gopkg.in/alecthomas/kingpin.v2"
)

type s3List []string

// Set overrides kingping's Set method to validate value for s3 URIs
func (s *s3List) Set(value string) error {
	hasMatch, err := regexp.MatchString("^s3://", value)
	if err != nil {
		return err
	}
	if !hasMatch {
		return fmt.Errorf("%s not a valid S3 uri, Please enter a valid S3 uri. Ex: s3://mary/had/a/little/lamb\n", value)
	} else {
		*s = append(*s, value)
		return nil
	}
}

func (s *s3List) String() string {
	return ""
}

// IsCumulative specifies S3List as a cumulative argument
func (s *s3List) IsCumulative() bool {
	return true
}

// S3List creates a new S3List kingpin setting
func S3List(s kingpin.Settings) (target *[]string) {
	target = new([]string)
	s.SetValue((*s3List)(target))
	return
}

func GetNumFileDescriptors() (uint64, error) {
	var rLimit syscall.Rlimit

	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	return rLimit.Cur, err
}

// getReaderByExt is a factory for reader based on the extension of the key
func GetReaderByExt(reader io.ReadCloser, key string) (io.ReadCloser, error) {
	ext := path.Ext(key)
	if ext == ".gz" || ext == ".gzip" {
		gzReader, err := gzip.NewReader(reader)
		if err != nil {
			return reader, nil
		}
		return gzReader, nil
	} else {
		return reader, nil
	}
}

// createPathIfNotExists takes a path and creates
// it if it doesn't exist
func CreatePathIfNotExists(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	} else {
		return nil
	}
	return nil
}
