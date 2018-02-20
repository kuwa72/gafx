package main

import (
	"fmt"
	"log"
	"os"

	"github.com/marcusolsson/tui-go"
)

/*
twopane
  scrollarea
    leftpain
  search
  scrollarea
    rightpain
  search
log
*/

/*
type Item struct {
	info    os.File
	display string
}
*/

type Pane struct {
	directory os.File // Directory
	items     []os.FileInfo
	searchKey string
}

func NewPane() (*Pane, error) {
	pane := Pane{}

	var dir string
	var file *os.File
	var files []os.FileInfo
	var err error

	if dir, err = os.Getwd(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	if file, err = os.Open(dir); err != nil {
		log.Fatal(err)
		return nil, err
	}

	if files, err = file.Readdir(0); err != nil {
		log.Fatal(err)
		return nil, err
	}

	pane.directory = *file
	pane.items = files

	return &pane, nil
}

func (p Pane) createWidget() *tui.Box {
	pwdField := tui.NewLabel(p.directory.Name())
	itemList := tui.NewList()

	for _, f := range p.items {
		itemList.AddItems(f.Name())
	}

	searchField := tui.NewTextEdit()
	searchField.SetText(p.searchKey)

	return tui.NewVBox(pwdField, itemList, searchField)
}

var leftItems *Pane
var righttems *Pane

func main() {
	// ignore error
	logout, _ := os.OpenFile("gafx.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	logger := log.New(logout, "gafx: ", log.Lshortfile|log.LstdFlags)
	if _, err := os.Getwd(); err != nil {
		log.Fatal("Can not get pwd")
		log.Fatal(err)
	}
	tui.SetLogger(logger)

	leftItems, _ = NewPane()
	righttems, _ = NewPane()
	fmt.Print(leftItems)

	messages := tui.NewTextEdit()

	topPane := tui.NewHBox(leftItems.createWidget(), righttems.createWidget())
	root := tui.NewVBox(topPane, messages)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
