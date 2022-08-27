package godo

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Application struct {
	Config *Config
	View   *tview.Application
}

func NewApplication(config Config) Application {
	itemStore := NewItemStore(config)

	app := tview.NewApplication()

	titleView := tview.NewTextView()
	titleView.
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorGreen)

	fmt.Fprint(titleView, ` _____       _     
|   __|___ _| |___ 
|  |  | . | . | . |
|_____|___|___|___|
`)

	hotkeyView := tview.NewTextView().
		SetDynamicColors(true)

	fmt.Fprintln(hotkeyView, "[yellow](e)[white]   edit")
	fmt.Fprintln(hotkeyView, "[yellow](x)[white]   complete")
	fmt.Fprintln(hotkeyView, "[yellow](d)[white]   delete")
	fmt.Fprintln(hotkeyView, "[yellow](n)[white]   jump to new")
	fmt.Fprintln(hotkeyView, "[yellow](esc)[white] quit")

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
				if list.GetItemCount() == 0 {
					app.SetFocus(inputField)
				}
			case 'd':
				selectedItem := list.GetCurrentItem()
				itemStore.RemoveItem(selectedItem)
				list.RemoveItem(selectedItem)
				if list.GetItemCount() == 0 {
					app.SetFocus(inputField)
				}
			case 'n':
				app.SetFocus(inputField)
			}
		case tcell.KeyDown:
			if list.GetCurrentItem() == list.GetItemCount()-1 {
				app.SetFocus(inputField)
			}
		}

		return event
	})

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

	for _, item := range itemStore.items {
		list.AddItem(item.String(), "", 0, nil)
	}

	grid := tview.NewGrid().
		SetRows(5, 0, 1).
		SetColumns(0, 20).
		SetBorders(true).
		AddItem(list, 0, 0, 2, 5, 0, 0, false).
		AddItem(titleView, 0, 5, 1, 1, 5, 20, false).
		AddItem(hotkeyView, 1, 5, 1, 1, 0, 0, false).
		AddItem(inputField, 2, 0, 1, 6, 0, 0, true)

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
