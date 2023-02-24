package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type fileType struct {
	filePath    string
	topic       string
	tag         []string
	link        []string
	bindingFile []string
	date        time.Time
}

var selectedFile fileType

/*
NOTE:
+конвертировать все файлы - избавиться от BOM \ufeff
+день недели убрать из даты
+link
теги то с большой то с маленькой буквы

перевести доп файлы на .csv
NewEntryWithData
*/

// var base []fileType

/*
fileRead - чтение файла filePath
достаем из него всю информацию и складываем в структуру
*/
func fileRead(filePath string) (f fileType) {
	bytes, err := os.ReadFile(filePath)
	if err != nil { // todo: err
		return
	}

	text := strings.Split(string(bytes), "\n")
	f.filePath = filePath

	// читаем однострочные
	for _, line := range text {
		if strings.Contains(line, "topic:") {
			f.topic = strings.TrimPrefix(line, "topic: ")
		}

		if strings.Contains(line, "tag:") { // # слитно не используется для форматирования
			line = strings.TrimPrefix(line, "tag: ")
			slice := strings.Split(line, " ")
			for _, s := range slice {
				if s != "" && s != "\r" {
					s = strings.TrimSuffix(s, " ")
					s = strings.TrimSuffix(s, "\r")
					f.tag = append(f.tag, s)
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

	// читаем многострочные
	copy := false
	minLen := 3
	for _, line := range text {

		if copy {
			if strings.Contains(line, "_____") {
				copy = false
				break
			}
			lineNext := strings.TrimSuffix(line, "\r")
			if len(lineNext) > minLen {
				f.link = append(f.link, lineNext)
			}
		}

		if strings.Contains(line, "link:") {
			copy = true
			line0 := strings.TrimPrefix(line, "link: ")
			line0 = strings.TrimSuffix(line0, "\r")
			if len(line0) > minLen {
				f.link = append(f.link, line0)
			}
		}
	}
	for _, line := range text {

		if copy {
			if strings.Contains(line, "_____") {
				copy = false
				break
			}
			lineNext := strings.TrimSuffix(line, "\r")
			if len(lineNext) > minLen {
				f.bindingFile = append(f.bindingFile, lineNext)
			}
		}

		if strings.Contains(line, "bind:") {
			copy = true
			line0 := strings.TrimPrefix(line, "bind: ")
			line0 = strings.TrimSuffix(line0, "\r")
			if len(line0) > minLen {
				f.bindingFile = append(f.bindingFile, line0)
			}
		}
	}

	return
}

/*
getText - прочитать файл,
вернуть только текст заметки
*/
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

func newlabel(labelName string) *widget.Label {
	l := widget.NewLabel(labelName)
	l.TextStyle.Monospace = true
	return l
}

/*
textEditor - открыть окно с текстом выбранного файла
Сохранить изменения
Закрыть без сохранения
*/
func textEditor(data fileType, text string) {
	statusLabel := widget.NewLabel("Тут что-нибудь отладочное...")

	w := fyne.CurrentApp().NewWindow("Типо текстовый редактор")
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(800, 600))

	tags := ""
	for _, tag := range selectedFile.tag {
		tags += tag + " "
	}

	fileNameEntry := widget.NewEntry()
	fileNameEntry.TextStyle.Monospace = true
	topicEntry := widget.NewEntry()
	topicEntry.TextStyle.Monospace = true
	tagEntry := widget.NewEntry()
	tagEntry.TextStyle.Monospace = true
	tagEntry.SetText(tags)
	dateEntry := widget.NewEntry()
	dateEntry.TextStyle.Monospace = true
	bindingMEntry := widget.NewMultiLineEntry()
	bindingMEntry.TextStyle.Monospace = true
	linkMEntry := widget.NewMultiLineEntry()
	linkMEntry.TextStyle.Monospace = true

	searchBox := container.NewVBox( // TODO: переименовать
		container.NewBorder(nil, nil, newlabel("Имя:     "), nil, fileNameEntry),
		container.NewBorder(nil, nil, newlabel("Тема:    "), nil, topicEntry),
		container.NewBorder(nil, nil, newlabel("Теги:    "), nil, tagEntry),
		container.NewBorder(nil, nil, newlabel("Дата:    "), nil, dateEntry),
		container.NewBorder(nil, nil, newlabel("Связать: "), nil, bindingMEntry),
		container.NewBorder(nil, nil, newlabel("Cсылки:  "), nil, linkMEntry),
	)

	fileNameEntry.SetText(filepath.Base(data.filePath))
	topicEntry.SetText(data.topic)
	dateEntry.SetText(data.date.Format("02.01.2006 15:04"))
	tagEntry.OnChanged = func(s string) {
		fmt.Println(s)
	}
	textEntry := widget.NewMultiLineEntry()
	textEntry.TextStyle.Monospace = true
	textEntry.Wrapping = fyne.TextWrapBreak
	textEntry.SetText(text)

	binds := ""
	for _, b := range data.bindingFile {
		binds += b + "\n\r"
	}
	bindingMEntry.SetText(binds)

	links := ""
	for _, b := range data.link {
		links += b + "\n\r"
	}
	linkMEntry.SetText(links)

	saveButton := widget.NewButton("Сохранить", func() {
		var d fileType
		if fileNameEntry.Text == "" {
			statusLabel.SetText("Введите имя файла") // todo: можно компактнее обрабатывать ощибки?
			return
		}
		if topicEntry.Text == "" {
			statusLabel.SetText("Заполните поле темы")
			return
		}
		if tagEntry.Text == "" {
			statusLabel.SetText("Заполните поле тегов")
			return
		}
		if dateEntry.Text == "" {
			statusLabel.SetText("Заполните поле даты")
			return
		} else {
			_, err := time.Parse("02.01.2006 15:04:05", dateEntry.Text)
			if err != nil {
				statusLabel.SetText("Неверно заполнено поле даты")
				return
			}
		}

		sl := strings.Split(tagEntry.Text, "#")
		for _, s := range sl {
			if s != "" && s != "\r" {
				s = strings.TrimSuffix(s, " ")
				d.tag = append(d.tag, s)
			}
		}
		sl = strings.Split(bindingMEntry.Text, "\n")
		for _, s := range sl {
			if s != "" && s != "\r" {
				s = strings.TrimSuffix(s, " ")
				d.bindingFile = append(d.tag, s)
			}
		}
		d.filePath = fileNameEntry.Text
		d.topic = topicEntry.Text
		d.date, _ = time.Parse("02.01.2006 15:04:05", dateEntry.Text)

		saveFile(d, textEntry.Text)
	})

	notSaveButton := widget.NewButton("Закрыть без сохранения", func() {
		d := dialog.NewConfirm("Вопрос", "Точно не сохранять?", func(b bool) {
			if b {
				w.Close()
			}
		}, w)
		d.SetDismissText("Hет")
		d.SetConfirmText("Да")
		d.Show()
	})
	btn := container.NewHBox(notSaveButton, layout.NewSpacer(), saveButton)
	bottomBox := container.NewVBox(
		statusLabel,
	)

	box := container.NewBorder(searchBox, container.NewBorder(nil, bottomBox, nil, btn), nil, nil, textEntry)
	w.SetContent(box)
	w.Show() // ShowAndRun -- panic!
}

/*
topic: Тема

tag: #tag1 #tag2

____________________________________________________________

Текст идеи своими словами.

____________________________________________________________

link: www.lingvolive.com

____________________________________________________________

bind:
2020 08 01 1019 Изучение языков.txt
____________________________________________________________

data: 2023.02.22 16:44


*/

func saveFile(data fileType, text string) error {
	sep := "____________________________________________________________\n\n"
	textall := "topic: " + data.topic + "\n\n"
	textall += "tag: "
	for _, tag := range data.tag {
		textall += tag
	}
	textall += "\n\n" + sep
	textall += text + "\n\n"
	textall += sep

	textall += "link: "
	for _, link := range data.link {
		textall += link
	}
	textall += "\n\n" + sep

	textall += "bind: "
	for _, bind := range data.bindingFile {
		textall += bind
	}
	textall += "\n\n" + sep

	textall += data.date.Format("02.01.2006 15:04:05")

	err := os.WriteFile(data.filePath, []byte(textall), 0666)
	return err
}
