package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

const prefix = "├───"
const endPrefix = "└───"
const separator = "│   "
const emptySeparator = "    "
const ending = "\n"

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, isPrintFiles bool) error {
	return printDir(out, path, isPrintFiles, "")
}

func printDir(out io.Writer, path string, isPrintFiles bool, parentPrefix string) error {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return errors.New("cannot read dir" + path)
	}

	printEntries := []os.DirEntry{}
	for _, entry := range dirEntries {
		if entry.IsDir() || isPrintFiles {
			printEntries = append(printEntries, entry)
		}
	}

	countPrintEntries := len(printEntries)
	for i, entry := range printEntries {
		isLast := countPrintEntries-1 == i

		err := printEntry(out, entry, parentPrefix, isLast)
		if err != nil {
			return errors.New("cannot write: " + entry.Name())
		}

		if entry.IsDir() {
			subDirPath := path + "/" + entry.Name()
			subDirSeparator := parentPrefix + separator
			if isLast {
				subDirSeparator = parentPrefix + emptySeparator
			}

			err = printDir(out, subDirPath, isPrintFiles, subDirSeparator)
			if err != nil {
				return errors.New("canot parse dir: " + subDirPath)
			}
		}
	}

	return nil
}

func printEntry(out io.Writer, entry os.DirEntry, parentSeparator string, isLast bool) error {
	var line = entry.Name()

	entryInfo,err := entry.Info()
	if err != nil {
		return err
	}

	if entryInfo.Size() == 0 {
		line += " (empty)"
	} else {
		line += fmt.Sprintf("(%db)", entryInfo.Size())
	}

	line +=  ending

	if isLast {
		line = parentSeparator + endPrefix + line
	} else {
		line = parentSeparator +prefix + line
	}

	_, err = out.Write([]byte(line))

	return err
}