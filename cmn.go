package main

import (
	"fyne.io/fyne/v2"
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

// func newFormatLabelAndEntry(name string) *fyne.Container {
// 	return container.NewBorder(nil, nil, newFormatLabel(name), nil, newFormatEntry())
// }

func newText() *widget.Entry {
	t := widget.NewMultiLineEntry()
	t.TextStyle.Monospace = true
	t.Wrapping = fyne.TextWrapWord
	return t
}

func formatSlice(sl []string) (text string) {
	for _, s := range sl {
		text += s + "\n"
	}
	return
}
