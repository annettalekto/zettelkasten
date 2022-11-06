package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type fileType struct {
	fileName    string
	topic       string // всегда одна?
	tag         []string
	link        string
	bindingFile []string
	data        time.Time
}

/*
todo
перевести доп файлы на .csv
link
день недели убрать из даты
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

	for i, line := range text { //отладка
		fmt.Printf("%d: %s\n", i, string(line))
	}

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
			s := strings.TrimPrefix(line, "data: ")
			s = strings.TrimSuffix(s, "\r")
			d, err := time.Parse("2006.01.02 15:04", s)
			fmt.Println("v", d.Weekday())
			fmt.Println(d, err)
			f.data = d
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

func textEditor(text string) { // текст должен сохранять форматирование
	// в разных окнах расположить остальные элементы

	w := fyne.CurrentApp().NewWindow("Типо текстовый редактор")
	w.Resize(fyne.NewSize(800, 600))

	textEntry := widget.NewMultiLineEntry()
	textEntry.Wrapping = fyne.TextWrapBreak
	textEntry.SetText(text)

	w.SetContent(textEntry)
	w.ShowAndRun()
}
