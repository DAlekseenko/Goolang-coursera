package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
)

type FileSystemObj struct {
	path   string
	name   string
	level  int
	prefix string
	size   string
}

func getSize(file os.FileInfo) string {
	var fileSize string

	if file.Size() == 0 {
		fileSize = "empty"
	} else {
		fileSize = strconv.FormatInt(file.Size(), 10) + "b"
	}

	return fileSize
}

func GetFsTree(out io.Writer, path string, printFiles bool) error {

	stack := []FileSystemObj{{path: path, name: path, level: 0, prefix: ""}}

	var output []FileSystemObj

	for len(stack) > 0 {
		n := len(stack) - 1 // Верхний элемент
		current := stack[n] // текущий объект
		stack = stack[:n]   // Pop
		dirContent, err := ioutil.ReadDir(current.path)

		if err != nil {
			return err
		}

		for _, file := range dirContent {

			fileName := file.Name()
			size := ""
			if !file.IsDir() {
				size += " (" + getSize(file) + ")"
			}
			path := current.path + string(os.PathSeparator) + fileName
			obj := FileSystemObj{path: path, level: current.level + 1, prefix: "", name: file.Name(), size: size}
			if printFiles && !file.IsDir() {
				output = append(output, obj)
			}
			if !file.IsDir() {
				continue
			}
			output = append(output, obj)
			stack = append(stack, obj)

		}
	}

	sort.SliceStable(output, func(i, j int) bool { return output[i].path < output[j].path })

	for len(output) > 0 {

		var prefix string
		current := output[0] // текущий объект
		output = output[1:]  // Режем для поиска

		currentLevelInPath := false
		for _, objs := range output {
			if objs.level == current.level {
				currentLevelInPath = true
			}
			if current.level > objs.level {
				break
			}
		}

		lessLevelOut := current.level

		for _, objs := range output {
			if lessLevelOut > objs.level {
				lessLevelOut = objs.level
				break
			}
		}

		if lessLevelOut == current.level {
			lessLevelOut = 0
		}

		for i := 0; i < lessLevelOut; i++ {
			prefix += BackSlashChar + "\t"
		}

		for i := 1; i < current.level-lessLevelOut; i++ {
			prefix += "\t"
		}

		if currentLevelInPath {
			prefix += BackChar
		} else {
			prefix += CornerChar
		}
		name := prefix + current.name
		if current.size != "" && printFiles {

			name += current.size
		}

		_, err := fmt.Fprintln(out, name)

		if err != nil {
			return fmt.Errorf(err.Error())
		}
	}

	return nil
}
