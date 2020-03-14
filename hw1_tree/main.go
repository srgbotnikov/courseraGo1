package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

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

func dirTree(out io.Writer, filePath string, printFiles bool) error {
	err := printDirTree("", out, filePath, printFiles)
	if err != nil {
		return err
	}
	return nil
}

func printDirTree(preString string, out io.Writer, filePath string, printFiles bool) error {
	files, err := createSortedArr(filePath, printFiles)
	if err != nil {
		fmt.Errorf("error createSortedArr %s", err.Error())
		return err
	}

	cntFiles := len(files)
	for i, f := range files {
		if f.IsDir() {
			var tmpPreStr string
			if cntFiles > i+1 {
				fmt.Fprintf(out, preString+"├───"+"%s\n", f.Name())
				tmpPreStr = preString + "│\t"
			} else {
				fmt.Fprintf(out, preString+"└───"+"%s\n", f.Name())
				tmpPreStr = preString + "\t"
			}
			newFilePath := filepath.Join(filePath, f.Name())
			printDirTree(tmpPreStr, out, newFilePath, printFiles)
		} else if printFiles {
			if f.Size() > 0 {
				if cntFiles > i+1 {
					fmt.Fprintf(out, preString+"├───%s (%vb)\n", f.Name(), f.Size())
				} else {
					fmt.Fprintf(out, preString+"└───%s (%vb)\n", f.Name(), f.Size())
				}
			} else {
				if cntFiles > i+1 {
					fmt.Fprintf(out, preString+"├───%s (empty)\n", f.Name())
				} else {
					fmt.Fprintf(out, preString+"└───%s (empty)\n", f.Name())
				}
			}
		}
	}
	return nil
}

func createSortedArr(filePath string, printFiles bool) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		fmt.Errorf("error read dir %s", err.Error())
		return nil, err
	}

	var filesMap map[string]os.FileInfo = map[string]os.FileInfo{}
	var keysFileName []string

	for _, f := range files {

		if printFiles {
			keysFileName = append(keysFileName, f.Name())
			filesMap[f.Name()] = f
		} else {
			if f.IsDir() {
				keysFileName = append(keysFileName, f.Name())
				filesMap[f.Name()] = f
			}
		}

	}
	sort.Strings(keysFileName)

	var sortedFilesArr []os.FileInfo = []os.FileInfo{}
	for _, k := range keysFileName {
		sortedFilesArr = append(sortedFilesArr, filesMap[k])
	}

	return sortedFilesArr, nil
}
