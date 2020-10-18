package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func getFSLevel(out io.Writer, path string, printFiles bool, prefix string) error {

	// содержимое текущей директории
	dirContent, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err.Error())
	}

	var directories []os.FileInfo

	// выкинуть файлы, если необходимо
	for _, file := range dirContent {
		if !printFiles && !file.IsDir() {
			continue
		}
		directories = append(directories, file)
	}

	for i, file := range directories {
		fileName := file.Name()
		if !file.IsDir() {
			fileName += " (" + getSize(file) + ")"
		}

		if i == len(directories)-1 {
			_, err = fmt.Fprintln(out, prefix+CornerChar+fileName)
		} else {
			_, err = fmt.Fprintln(out, prefix+BackChar+fileName)
		}

		if err != nil {
			return fmt.Errorf(err.Error())
		}

		currentPrefix := prefix

		if file.IsDir() {
			if i != len(directories)-1 {
				currentPrefix += BackSlashChar + "\t"
			} else {
				currentPrefix += "\t"
			}
			err = getFSLevel(out, path+string(os.PathSeparator)+file.Name(), printFiles, currentPrefix)
			if err != nil {
				return fmt.Errorf(err.Error())
			}
		}
	}

	return nil
}

func GetFSTreeRecursive(out io.Writer, path string, printFiles bool) error {
	return getFSLevel(out, path, printFiles, "")
}
