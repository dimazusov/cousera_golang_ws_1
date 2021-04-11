package main

import (
	"errors"
	"io"
	"os"
)

type printer struct {
	out       io.Writer
	prefix    string
	endPrefix string
	separator string
}

func NewPrinter(out io.Writer, prefix string, endPrefix string, separator string) *printer {
	return &printer{
		out:       out,
		prefix:    prefix,
		endPrefix: endPrefix,
		separator: separator,
	}
}

func (m printer) PrintDir(path string, isPrintFiles bool, level int) error {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return errors.New("cannot read dir" + path)
	}

	for _, entry := range dirEntries {
		if entry.IsDir() || isPrintFiles {
			line := prefix + entry.Name() + ending

			err := m.printLine(line, level, false)
			if err != nil {
				return errors.New("cannot write: " + line)
			}
		}

		if entry.IsDir() {
			subDirPath := path + "/" + entry.Name()

			err = m.PrintDir(subDirPath, isPrintFiles, level+1)
			if err != nil {
				return errors.New("canot parse dir: " + subDirPath)
			}
		}
	}

	return nil
}

func (m printer) printLine(text string, nestingLevel int, isLast bool) error {
	for i := 0; i < nestingLevel; i++ {
		_, err := m.out.Write([]byte(separator))
		if err != nil {
			return err
		}
	}

	_, err := m.out.Write([]byte(text))

	return err
}
