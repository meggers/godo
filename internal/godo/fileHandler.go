package godo

import (
	"bufio"
	"fmt"
	"os"
)

func ReadTodoFile(fileName string) []TodoItem {
	file, err := os.Open(fileName)
	defer closeFile(file)
	check(err)

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
	check(err)

	writer := bufio.NewWriter(file)

	for _, todoItem := range todoItems {
		_, err := writer.WriteString(todoItem.String() + "\n")
		check(err)
	}

	writer.Flush()
}

func closeFile(f *os.File) {
	err := f.Close()
	check(err)
}

func check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
