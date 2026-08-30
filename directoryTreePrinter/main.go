package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	data := map[string][]string{}

	folder, err := os.ReadDir("L1")
	if err != nil {
		fmt.Printf("error while reading input folder: %v", err)
	}

	// fmt.Println(folder)
	// fmt.Printf("%v : %v\n", file.Name(), file.IsDir())

	getTrace(folder, "L1", data)

	fmt.Println("L1")
	printTree(data, "L1", "")

}

func getTrace(folder []os.DirEntry, root string, data map[string][]string) {
	for _, file := range folder {

		fileName := ""
		if strings.Contains(root, "\\") {
			fileNameArr := strings.Split(root, "\\")
			fileName = fileNameArr[len(fileNameArr)-1]
		} else {
			fileName = root
		}

		data[fileName] = append(data[fileName], file.Name())
		// fmt.Println(file.Name())
		if file.IsDir() {
			folder, err := os.ReadDir(filepath.Join(root, file.Name()))
			if err != nil {
				fmt.Printf("error inside getTrace read: %v\n", err)
			}
			getTrace(folder, filepath.Join(root, file.Name()), data)
		}

	}
}

func printTree(tree map[string][]string, node string, prefix string) {
	children := tree[node]

	for i, child := range children {
		isLast := i == len(children)-1

		branch := "├── "
		if isLast {
			branch = "└── "
		}

		fmt.Println(prefix + branch + child)

		childPrefix := prefix + "│   "
		if isLast {
			childPrefix = prefix + "    "
		}

		printTree(tree, child, childPrefix)
	}
}
