package data

import (
	"io/ioutil"
	"os"
)

type Reader interface {
	Read(path string) ([]byte, error)
}
type reader struct{}

func NewReader() Reader {
	return &reader{}
}

func (r *reader) Read(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	body, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return body, nil
}
