package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func newFormatLabel(name string) *widget.Label {
	l := widget.NewLabel(name)
	l.TextStyle.Monospace = true
	return l
}

func newFormatEntry() *widget.Entry {
	en := widget.NewEntry()
	en.TextStyle.Monospace = true
	return en
}

// todo: нужно?
func newFormatLabelAndEntry(name string) *fyne.Container {
	return container.NewBorder(nil, nil, newFormatLabel(name), nil, newFormatEntry())
}
