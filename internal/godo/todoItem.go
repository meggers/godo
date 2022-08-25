package godo

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	dateFormat     = "2006-01-02"
	completedToken = "x"
)

type TodoItem struct {
	ID             int
	Complete       bool
	Priority       string
	CompletionDate time.Time
	CreationDate   time.Time
	Description    string
}

func NewTodoItem(id int, line string) TodoItem {
	tokens := strings.Fields(line)
	item := TodoItem{
		ID: id,
	}

	var descriptionStartIndex int
	for i, token := range tokens {
		// Parse Completion
		if i == 0 && token == completedToken {
			item.Complete = true
			continue
		}

		// Parse Priority
		if i == 0 || (i == 1 && item.Complete) {
			matched, _ := regexp.MatchString(`\([A-Z]\)`, token)

			if matched {
				item.Priority = strings.Trim(token, "()")
				continue
			}
		}

		// Parse Completion and Creation dates
		parsedDate, err := time.Parse(dateFormat, token)
		if err == nil {
			if item.Complete && item.CompletionDate.IsZero() {
				item.CompletionDate = parsedDate
				continue
			} else if item.CreationDate.IsZero() {
				item.CreationDate = parsedDate
				continue
			}
		}

		descriptionStartIndex = i
		break
	}

	descriptionTokens := tokens[descriptionStartIndex:]
	item.Description = strings.Join(descriptionTokens, " ")

	return item
}

func (item TodoItem) String() string {
	var tokens []string

	if item.Complete {
		tokens = append(tokens, completedToken)
	}

	if item.Priority != "" {
		tokens = append(tokens, fmt.Sprintf("(%v)", item.Priority))
	}

	if !item.CompletionDate.IsZero() {
		tokens = append(tokens, item.CompletionDate.Format(dateFormat))
	}

	if !item.CreationDate.IsZero() {
		tokens = append(tokens, item.CreationDate.Format(dateFormat))
	}

	tokens = append(tokens, item.Description)
	return strings.Join(tokens, " ")
}
