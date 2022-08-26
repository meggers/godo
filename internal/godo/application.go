package godo

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Application struct {
	Config *Config
	View   *tview.Application
}

const commandList = "create,sort"

func NewApplication(config Config) Application {
	commands := strings.Split(commandList, ",")
	itemStore := NewItemStore(config)

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
				itemStore.ArchiveItem(selectedItem)
				list.RemoveItem(selectedItem)
			case 'd':
				selectedItem := list.GetCurrentItem()
				itemStore.RemoveItem(selectedItem)
				list.RemoveItem(selectedItem)
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
			newItem := NewTodoItem(inputText)
			itemStore.AddItem(newItem)
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

	for _, item := range itemStore.items {
		list.AddItem(item.String(), "", 0, nil)
	}

	app.SetRoot(grid, true)

	return Application{
		Config: &config,
		View:   app,
	}
}

func (app Application) Run() {
	if err := app.View.Run(); err != nil {
		panic(err)
	}
}

func remove(slice []TodoItem, s int) []TodoItem {
	return append(slice[:s], slice[s+1:]...)
}
