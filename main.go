package main

import (
	"github.com/meggers/godo/internal/godo"
)

func main() {
	config := godo.NewConfig()
	godo.NewApplication(*config).Run()
}
