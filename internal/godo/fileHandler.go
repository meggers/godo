package godo

import (
	"bufio"
	"os"

	util "github.com/meggers/godo/internal"
)

// read all todos
// write all todos
// append to archive
func ReadTodoFile(fileName string) ([]TodoItem, error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDONLY, 0)
	defer closeFile(file)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)
	var todoItems []TodoItem

	lineNumber := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			todoItem := NewTodoItem(scanner.Text())
			todoItems = append(todoItems, todoItem)
			lineNumber++
		}
	}

	return todoItems, nil
}

func WriteTodoFile(fileName string, todoItems []TodoItem) error {
	file, err := os.Create(fileName)
	defer closeFile(file)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)

	for _, todoItem := range todoItems {
		_, err := writer.WriteString(todoItem.String() + "\n")
		util.CheckError(err, "failed to write item to file")
	}

	writer.Flush()
	return nil
}

func WriteArchiveFile(fileName string, todoItems []TodoItem) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	defer closeFile(file)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)

	for _, todoItem := range todoItems {
		_, err := writer.WriteString(todoItem.String() + "\n")
		util.CheckError(err, "failed to write item to file")
	}

	writer.Flush()
	return nil
}

func closeFile(f *os.File) {
	err := f.Close()
	util.CheckError(err, "failed to close file")
}
