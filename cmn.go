package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// func newFormatLabel(name string) *widget.Label {
// 	l := widget.NewLabel(name)
// 	l.TextStyle.Monospace = true
// 	return l
// }

func newFormatEntry() *widget.Entry {
	en := widget.NewEntry()
	en.TextStyle.Monospace = true
	return en
}

func newText() *widget.Entry {
	t := widget.NewMultiLineEntry()
	t.TextStyle.Monospace = true
	t.Wrapping = fyne.TextWrapWord
	return t
}

func sliceInColumn(sl []string) (text string) {
	for _, s := range sl {
		text += s + "\n"
	}
	return strings.TrimSuffix(text, "\n")
}

func sliceInString(sl []string) (text string) {
	for _, s := range sl {
		text += s + ", "
	}
	return strings.TrimSuffix(text, ", ")
}
