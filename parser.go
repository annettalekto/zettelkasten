package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type fileType struct {
	fileName    string
	topic       string // всегда одна?
	tag         []string
	link        string
	bindingFile []string
	date        time.Time
}

/*
todo
+конвертировать все файлы - избавиться от BOM \ufeff
+день недели убрать из даты
+link
теги то с большой то с маленькой буквы

перевести доп файлы на .csv
NewEntryWithData
*/

// var base []fileType

func fileRead(fileName string) (f fileType) {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		// todo add status line
		return
	}
	text := strings.Split(string(bytes), "\n")

	// for i, line := range text { //отладка
	// 	fmt.Printf("%d: %s\n", i, string(line))
	// }

	f.fileName = fileName

	for _, line := range text {
		if strings.Contains(line, "topic:") {
			f.topic = strings.TrimPrefix(line, "topic: ")
		}

		if strings.Contains(line, "#") {
			slice := strings.Split(line, "#")
			for _, s := range slice {
				if s != "" && s != "\r" {
					s = strings.TrimSuffix(s, " ")
					f.tag = append(f.tag, s)
				}
			}
		}

		if strings.Contains(line, "link:") {
			s := strings.TrimSuffix(line, "\r")
			f.link = strings.TrimPrefix(s, "link: ")
		}

		if strings.Contains(line, "[") { // todo парс строк
			slice := strings.Split(line, "[")
			for _, s := range slice {
				s = strings.TrimSuffix(s, "\r")
				s = strings.TrimSuffix(s, "]")
				if s != "" {
					f.bindingFile = append(f.bindingFile, s)
				}
			}
		}

		if strings.Contains(line, "data:") {
			// если есть дол. символы Parse не работает
			s := strings.TrimPrefix(line, "data: ")

			s = strings.TrimSuffix(s, "\r")
			s = strings.TrimSpace(s)
			d, err := time.Parse("2006.01.02 15:04", s)
			fmt.Println("v", d.Weekday())
			fmt.Println(d, err)
			f.date = d
		}
	}

	// base = append(base, val)
	return
}

func getText(fileName string) (fileText string) {

	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return
	}
	text := strings.Split(string(bytes), "\n")

	copy := false
	for _, line := range text {

		if copy {
			if strings.Contains(line, "_____") {
				break
			}
			fileText += line + "\r\n"
		}
		if strings.Contains(line, "_____") {
			copy = true
		}
	}

	return
}

func label(labelName string) *widget.Label {
	l := widget.NewLabel(labelName)
	l.TextStyle.Monospace = true
	return l
}

func textEditor(data fileType, text string) { // текст должен сохранять форматирование
	// в разных окнах расположить остальные элементы

	w := fyne.CurrentApp().NewWindow("Типо текстовый редактор")
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(800, 600))

	fileNameEntry := widget.NewEntry()
	fileNameEntry.TextStyle.Monospace = true
	topicEntry := widget.NewEntry()
	topicEntry.TextStyle.Monospace = true
	tagEntry := widget.NewEntry()
	tagEntry.TextStyle.Monospace = true
	dateEntry := widget.NewEntry()
	dateEntry.TextStyle.Monospace = true

	searchBox := container.NewVBox(
		container.NewBorder(nil, nil, label("Имя:  "), nil, fileNameEntry),
		container.NewBorder(nil, nil, label("Тема: "), nil, topicEntry),
		container.NewBorder(nil, nil, label("Теги: "), nil, tagEntry),
		container.NewBorder(nil, nil, label("Дата: "), nil, dateEntry),
	)

	fileNameEntry.SetText(data.fileName) // только имя todo
	topicEntry.SetText(data.topic)
	tagEntry.SetText(data.tag[0]) // todo все, с #, с пробелом?
	dateEntry.SetText(data.date.Format("02.01.2006 15:04"))

	textEntry := widget.NewMultiLineEntry()
	textEntry.TextStyle.Monospace = true
	textEntry.Wrapping = fyne.TextWrapBreak
	textEntry.SetText(text)

	saveButton := widget.NewButton("Сохранить", func() {
	})

	box := container.NewBorder(searchBox, container.NewBorder(nil, nil, nil, saveButton), nil, nil, textEntry)
	w.SetContent(box)
	w.Show() // ShowAndRun -- panic!
}
