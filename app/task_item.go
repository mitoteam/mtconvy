package app

import "fmt"

type TaskItem struct {
	Name string
	Path string
}

func (i TaskItem) SelectStreams() {
	fmt.Printf("*** FILE: %s", i.Name)

	FfGetStreamList(i.Path)
}
