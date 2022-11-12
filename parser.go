package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type fileType struct {
	filePath    string
	topic       string // всегда одна?
	tag         []string
	link        string
	bindingFile []string
	date        time.Time
}

var selectedFile fileType

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

func fileRead(filePath string) (f fileType) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		// todo add status line
		return
	}
	text := strings.Split(string(bytes), "\n")

	f.filePath = filePath

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

func getText(filePath string) (fileText string) {

	bytes, err := os.ReadFile(filePath)
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

func textEditor(data fileType, text string) {

	w := fyne.CurrentApp().NewWindow("Типо текстовый редактор")
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(800, 600))

	var tagSlise []string
	for _, s := range data.tag {
		tagSlise = append(tagSlise, "#"+s)
	}
	fileNameEntry := widget.NewEntry()
	fileNameEntry.TextStyle.Monospace = true
	topicEntry := widget.NewEntry()
	topicEntry.TextStyle.Monospace = true
	tagSelectEntry := widget.NewSelectEntry(tagSlise)
	tagSelectEntry.TextStyle.Monospace = true
	if len(tagSlise) > 0 {
		tagSelectEntry.SetText(tagSlise[0])
	}
	dateEntry := widget.NewEntry()
	dateEntry.TextStyle.Monospace = true

	searchBox := container.NewVBox(
		container.NewBorder(nil, nil, label("Имя:  "), nil, fileNameEntry),
		container.NewBorder(nil, nil, label("Тема: "), nil, topicEntry),
		container.NewBorder(nil, nil, label("Теги: "), nil, tagSelectEntry),
		container.NewBorder(nil, nil, label("Дата: "), nil, dateEntry),
	)

	fileNameEntry.SetText(filepath.Base(data.filePath))
	topicEntry.SetText(data.topic)
	dateEntry.SetText(data.date.Format("02.01.2006 15:04"))
	tagSelectEntry.OnChanged = func(s string) {
		fmt.Println(s)
	}
	textEntry := widget.NewMultiLineEntry()
	textEntry.TextStyle.Monospace = true
	textEntry.Wrapping = fyne.TextWrapBreak
	textEntry.SetText(text)

	saveButton := widget.NewButton("Сохранить", func() {
		// считать в data и text
		// сохранить в папку
		// добавить в список файлов слева
	})

	box := container.NewBorder(searchBox, container.NewBorder(nil, nil, nil, saveButton), nil, nil, textEntry)
	w.SetContent(box)
	w.Show() // ShowAndRun -- panic!
}
