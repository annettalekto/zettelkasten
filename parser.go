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
	topic       string // всегда одна?
	tag         []string
	link        string
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
	if err != nil {
		// FIXME: add status line
		return
	}

	text := strings.Split(string(bytes), "\n")

	f.filePath = filePath

	for _, line := range text {
		if strings.Contains(line, "topic:") {
			f.topic = strings.TrimPrefix(line, "topic: ")
		}

		if strings.Contains(line, "#") { // fixme: # может быть использована для форматирования. переделать на метку tag
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

		if strings.Contains(line, "[") { // TODO: парс строк
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
		tags += "#" + tag + " "
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

	searchBox := container.NewVBox( // TODO: переименовать
		container.NewBorder(nil, nil, newlabel("Имя:  "), nil, fileNameEntry),
		container.NewBorder(nil, nil, newlabel("Тема: "), nil, topicEntry),
		container.NewBorder(nil, nil, newlabel("Теги: "), nil, tagEntry), //  fixme: tagEntry
		container.NewBorder(nil, nil, newlabel("Дата: "), nil, dateEntry),
	)

	fileNameEntry.SetText(filepath.Base(data.filePath))
	topicEntry.SetText(data.topic)
	dateEntry.SetText(data.date.Format("02.01.2006 15:04")) // TODO:  dateEntry ->type DatePicker
	tagEntry.OnChanged = func(s string) {
		fmt.Println(s)
	}
	textEntry := widget.NewMultiLineEntry()
	textEntry.TextStyle.Monospace = true
	textEntry.Wrapping = fyne.TextWrapBreak
	textEntry.SetText(text)

	saveButton := widget.NewButton("Сохранить", func() {
		var d fileType
		d.filePath = fileNameEntry.Text
		d.topic = topicEntry.Text
		// tags
		sl := strings.Split(tagEntry.Text, "#") // debug: проверить тут
		for _, s := range sl {
			if s != "" && s != "\r" {
				s = strings.TrimSuffix(s, " ")
				d.tag = append(d.tag, s)
			}
		}
		// d.date =  dateEntry ->type DatePicker сначала переделать элемент
		// сохранить в папку файл в соответствии с паттерном
		// теги сохранить в общий файл?
		// data добавить в слайс, обновить список файлов слева?
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

	box := container.NewBorder(searchBox, container.NewBorder(nil, statusLabel, nil, btn), nil, nil, textEntry)
	w.SetContent(box)
	w.Show() // ShowAndRun -- panic!
}

func saveFile(data fileType, dir, text string) {

}
