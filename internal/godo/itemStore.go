package godo

import (
	"path/filepath"
	"time"

	util "github.com/meggers/godo/internal"
)

type ItemStore struct {
	Config *Config
	items  []TodoItem
}

func NewItemStore(config Config) ItemStore {
	items, err := ReadTodoFile(config.TodoFile)
	util.CheckError(err, "failed to read todo file")

	return ItemStore{
		Config: &config,
		items:  items,
	}
}

func (itemStore *ItemStore) AddItem(item TodoItem) {
	itemStore.items = append(itemStore.items, item)
	itemStore.save()
}

func (itemStore *ItemStore) RemoveItem(s int) {
	itemStore.items = remove(itemStore.items, s)
	itemStore.save()
}

func (itemStore *ItemStore) ArchiveItem(s int) {
	item := itemStore.items[s]
	item.Complete = true

	if itemStore.Config.EnableArchive {
		archivePath := itemStore.getArchivePath()
		err := WriteArchiveFile(archivePath, []TodoItem{item})
		util.CheckError(err, "failed to write archive file")
	}

	itemStore.RemoveItem(s)
}

func (itemStore *ItemStore) save() {
	err := WriteTodoFile(itemStore.Config.TodoFile, itemStore.items)
	util.CheckError(err, "failed to write todo file")
}

func (itemStore ItemStore) getArchivePath() string {
	dateString := ""
	if itemStore.Config.ArchiveFormat != "" {
		dateString = time.Now().Format(itemStore.Config.ArchiveFormat)
		dateString += "-"
	}

	fileName := dateString + "archive.txt"
	return filepath.Join(itemStore.Config.ArchiveDirectory, fileName)
}
