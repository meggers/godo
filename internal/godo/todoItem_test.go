package godo

import "testing"

func TestParseDescriptionOnly(t *testing.T) {
	got := NewTodoItem("a description")

	if got.Description != "a description" {
		t.Errorf("Abs(a description) = %v; want a description", got.Description)
	}
}
