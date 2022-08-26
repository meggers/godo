package godo

import (
	"bufio"
	"os"

	util "github.com/meggers/godo/internal"
)

func ReadTodoFile(fileName string) []TodoItem {
	file, err := os.OpenFile(fileName, os.O_CREATE&os.O_RDONLY, 0)
	defer closeFile(file)
	util.CheckError(err, "failed to open todo file")

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)
	var todoItems []TodoItem

	lineNumber := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			todoItem := NewTodoItem(lineNumber, scanner.Text())
			todoItems = append(todoItems, todoItem)
			lineNumber++
		}
	}

	return todoItems
}

func WriteTodoFile(fileName string, todoItems []TodoItem) {
	file, err := os.Create(fileName)
	defer closeFile(file)
	util.CheckError(err, "failed to open todo file")

	writer := bufio.NewWriter(file)

	for _, todoItem := range todoItems {
		_, err := writer.WriteString(todoItem.String() + "\n")
		util.CheckError(err, "failed to write item to file")
	}

	writer.Flush()
}

func closeFile(f *os.File) {
	err := f.Close()
	util.CheckError(err, "failed to close file")
}
