package logging

import (
	"fmt"
	"io"
	"os"
)

func NewMultiWriter(writers ...io.Writer) io.Writer {
	return io.MultiWriter(writers...)
}

func NewConsoleWriter() io.Writer {
	return os.Stdout
}

func NewFileWriter(path string) (io.Writer, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("error while opening file: %s", err.Error())
	}
	return file, nil
}
