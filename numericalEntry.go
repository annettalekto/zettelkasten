package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type numericalEntry struct {
	widget.Entry
}

// создаем функцию конструктор на основе базового типа
func newNumericalEntry() *numericalEntry {
	entry := &numericalEntry{}
	entry.ExtendBaseWidget(entry)
	entry.Wrapping = fyne.TextTruncate
	return entry
}

// переопределяем стандартный метод обработки рун (overriding the TypedRune(rune) method)
// если руна не отвечает условию, просто игнорируем
// иначе делегируем стандартной функции  e.Entry.TypedRune
func (e *numericalEntry) TypedRune(r rune) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', ',':
		// fmt.Printf("Pressed %v\n", r)
		e.Entry.TypedRune(r)
	}
}

// перезаписывает метод TypedShortcut(fyne.Shortcut) чтобы не давал буквы вставлять методом копирования
// иначе ввести можно только цифру (допустимые руны)
// а вставить можно любой текст
func (e *numericalEntry) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	content := paste.Clipboard.Content()
	if _, err := strconv.ParseFloat(content, 64); err == nil {
		e.Entry.TypedShortcut(shortcut)
	}
}
