package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const commandList = "create,sort"

func main() {
	commands := strings.Split(commandList, ",")
	config := newConfig()

	todoItems := readTodoFile(config.TodoFile)

	app := tview.NewApplication()
	inputField := tview.NewInputField().
		SetLabel("> ").
		SetFieldWidth(0)
	list := tview.NewList().
		SetWrapAround(false).
		SetSelectedFocusOnly(true)
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'e':
				selectedItem := list.GetCurrentItem()
				text, _ := list.GetItemText(selectedItem)
				inputField.SetText(text)
				app.SetFocus(inputField)
			case 'x':
				selectedItem := list.GetCurrentItem()
				todoItems[selectedItem].Complete = true
				list.SetItemText(selectedItem, todoItems[selectedItem].String(), "")
				writeTodoFile(config.TodoFile, todoItems)
			case 'd':
				selectedItem := list.GetCurrentItem()
				list.RemoveItem(selectedItem)
				todoItems = remove(todoItems, selectedItem)
				writeTodoFile(config.TodoFile, todoItems)
			}
		case tcell.KeyDown:
			if list.GetCurrentItem() == list.GetItemCount()-1 {
				app.SetFocus(inputField)
			}
		}

		return event
	})

	grid := tview.NewGrid().
		SetRows(0, 1).
		SetBorders(true).
		AddItem(list, 0, 0, 1, 1, 0, 0, true).
		AddItem(inputField, 1, 0, 1, 1, 0, 0, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			currentFocus := app.GetFocus()
			if currentFocus == inputField {
				app.SetFocus(list)
			} else {
				app.SetFocus(inputField)
			}
		case tcell.KeyEsc:
			app.Stop()
			return nil
		}

		return event
	})

	inputField.SetFieldBackgroundColor(tcell.ColorBlack)
	inputField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			inputText := inputField.GetText()
			newItem := newTodoItem(len(todoItems), inputText)
			todoItems = append(todoItems, newItem)
			writeTodoFile(config.TodoFile, todoItems)
			list.AddItem(newItem.String(), "", 0, nil)
			list.SetCurrentItem(list.GetItemCount() - 1)
			inputField.SetText("")
		default:
			inputField.SetText("")
		}
	})

	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			list.SetCurrentItem(list.GetItemCount() - 1)
			app.SetFocus(list)
		}

		return event
	})

	inputField.SetAutocompleteFunc(func(currentText string) (entries []string) {
		if len(currentText) == 0 {
			return
		}
		for _, word := range commands {
			if strings.HasPrefix(strings.ToLower(word), strings.ToLower(currentText)) {
				entries = append(entries, word)
			}
		}
		return
	})

	for _, item := range todoItems {
		list.AddItem(item.String(), "", 0, nil)
	}

	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}
}

func remove(slice []todoItem, s int) []todoItem {
	return append(slice[:s], slice[s+1:]...)
}
