package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func isLastElement(idx int, infos []os.FileInfo, printFiles bool) bool {
	if idx == len(infos)-1 {
		return true
	}
	if printFiles {
		return false
	}
	subRange := infos[idx+1:]
	for _, info := range subRange {
		if info.IsDir() {
			return false
		}
	}
	return true
}

func treeRecursion(path string, prefix string, out io.Writer, printFiles bool) {
	fileInfos, _ := ioutil.ReadDir(path)
	for idx, info := range fileInfos {
		if !info.IsDir() && !printFiles {
			continue
		}

		var nextPrefix, outStr string
		if isLastElement(idx, fileInfos, printFiles) {
			outStr = prefix + "└───" + info.Name()
			nextPrefix = prefix + "\t"
		} else {
			outStr = prefix + "├───" + info.Name()
			nextPrefix = prefix + "│" + "\t"
		}
		if info.IsDir() {
			fmt.Fprintln(out, outStr)
			subdir := filepath.Join(path, info.Name())
			treeRecursion(subdir, nextPrefix, out, printFiles)
		} else {
			if info.Size() == 0 {
				outStr += " (empty)"
			} else {
				outStr += " (" + strconv.FormatInt(info.Size(), 10) + "b)"
			}
			fmt.Fprintln(out, outStr)
		}
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {

	treeRecursion(path, "", out, printFiles)

	return nil
}

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
