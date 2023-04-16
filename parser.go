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

type fileType struct { // ztcElementsType ztcBasicsType
	filePath     string // полное имя файла с путем и расширением файла
	topic        string
	tags         []string
	links        []string
	bindingFiles []string
	date         time.Time // ну текст же
}

var selectedFile fileType

/*
fileRead - чтение файла filePath
достаем из него всю информацию и складываем в структуру
*/
// func fileRead(filePath string) (f fileType) { //todo: переименовать
// 	bytes, err := os.ReadFile(filePath)
// 	if err != nil { // todo: err
// 		return
// 	}

// 	text := strings.Split(string(bytes), "\n")
// 	f.filePath = filePath

// 	// читаем однострочные
// 	for _, line := range text {
// 		if strings.Contains(line, "topic:") {
// 			//line = strings.TrimPrefix(line, "topic:")
// 			// f.topic = strings.TrimPrefix(line, " ")
// 			f.topic = mTrimPrefix(line, "topic:")
// 		}

// 		if strings.Contains(line, "tag:") { // # слитно не используется для форматирования
// 			// line = strings.TrimPrefix(line, "tag:")
// 			// line = strings.TrimPrefix(line, " ")
// 			line = mTrimPrefix(line, "tag:")
// 			slice := strings.Split(line, " ")
// 			for _, s := range slice {
// 				if s != "" {
// 					// s = strings.TrimSuffix(s, " ")
// 					// s = strings.TrimSuffix(s, "\r")
// 					f.tags = append(f.tags, s)
// 				}
// 			}
// 		}

// 		if strings.Contains(line, "date:") {
// 			// если есть доп. символы Parse не работает
// 			s := mTrimPrefix(line, "date:")
// 			// s = strings.TrimPrefix(s, " ")
// 			// s = strings.TrimSuffix(s, "\r")
// 			// s = strings.TrimSpace(s)
// 			d, _ := time.Parse("02.01.2006 15:04", s)
// 			// fmt.Println("v", d.Weekday())
// 			// fmt.Println("v", d.Year())
// 			f.date = d
// 		}
// 	}

// 	// читаем многострочные
// 	copy := false
// 	minLen := 3
// 	for _, line := range text {

// 		if copy {
// 			if strings.Contains(line, "_____") {
// 				copy = false
// 				break
// 			}
// 			lineNext := mTrimPrefix(line, " ")
// 			if len(lineNext) > minLen {
// 				f.links = append(f.links, lineNext)
// 			}
// 		}

// 		if strings.Contains(line, "link:") {
// 			copy = true
// 			line0 := mTrimPrefix(line, "link:")
// 			// line0 = strings.TrimPrefix(line0, " ")
// 			// line0 = strings.TrimSuffix(line0, "\r")
// 			if len(line0) > minLen {
// 				f.links = append(f.links, line0)
// 			}
// 		}
// 	}
// 	for _, line := range text {

// 		if copy {
// 			if strings.Contains(line, "_____") {
// 				copy = false
// 				break
// 			}
// 			lineNext := mTrimPrefix(line, " ")
// 			if len(lineNext) > minLen {
// 				f.bindingFiles = append(f.bindingFiles, lineNext)
// 			}
// 		}

// 		if strings.Contains(line, "bind:") {
// 			copy = true
// 			line0 := mTrimPrefix(line, "bind:")
// 			// line0 = strings.TrimPrefix(line0, " ")
// 			// line0 = strings.TrimSuffix(line0, "\r")
// 			if len(line0) > minLen {
// 				f.bindingFiles = append(f.bindingFiles, line0)
// 			}
// 		}
// 	}

// 	return
// }

//Note: а если читать файл несколько раз? не переломиться... наверное

func getTopic(filePath string) (s string, err error) {
	bytes, _ := os.ReadFile(filePath) // ошибка проверена в ф. выше
	// if err != nil {
	// 	err = fmt.Errorf("file read (%s) error: %s", filePath, err.Error())
	// 	return
	// }

	text := strings.Split(string(bytes), "\r\n")      //note: a new line character in Windows
	text[0], _ = strings.CutPrefix(text[0], "\ufeff") // cut BOM

	copy := false
	stemp := ""
	for _, line := range text {
		if strings.Contains(line, "<topic>") {
			copy = true
		}
		if copy {
			stemp += line
		}
		if strings.Contains(line, "</topic>") {
			copy = false
		}
	}
	if stemp == "" {
		err = fmt.Errorf("topic не найден")
		return
	}
	s, err = cutTags("topic", stemp) // del tags
	return
}

// todo: вынести в отдельный файл
func cutTags(tagName, before string) (s string, err error) {
	ok := false
	s = before // не изменять в случаи ошибки

	before, ok = strings.CutPrefix(before, "<"+tagName+">")
	if !ok {
		err = fmt.Errorf("cut tags prexic error: %s", tagName)
		return
	}
	before, ok = strings.CutSuffix(before, "</"+tagName+">")
	if !ok {
		err = fmt.Errorf("cut tags suffxic error: %s", tagName)
		return
	}
	s = before

	return
}

// todo: переделать на теги
func fileRead(filePath string) (f fileType, err error) { //todo: переименовать
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		err = fmt.Errorf("file read (%s) error: %s", filePath, err.Error())
		fmt.Println(err)
		return
	}
	f.filePath = filePath

	f.topic, err = getTopic(filePath)
	if err != nil {
		fmt.Println(err) // todo: можно в лог писать кстати, а критичные еще в статус строке
	}
	// fmt.Println(f.topic, err) // debug ok

	// todo: как писать в лог правильно, настроить вывод ош в статус лейбл. может какие то уровни ошибок сделать. читать

	//--------------
	text := strings.Split(string(bytes), "\n\r")
	// читаем однострочные
	for _, line := range text {
		if strings.Contains(line, "<topic>") {
			//line = strings.TrimPrefix(line, "topic:")
			// f.topic = strings.TrimPrefix(line, " ")
			f.topic = mTrimPrefix(line, "topic:")
		}

		/*if strings.Contains(line, "tag:") { // # слитно не используется для форматирования
			// line = strings.TrimPrefix(line, "tag:")
			// line = strings.TrimPrefix(line, " ")
			line = mTrimPrefix(line, "tag:")
			slice := strings.Split(line, " ")
			for _, s := range slice {
				if s != "" {
					// s = strings.TrimSuffix(s, " ")
					// s = strings.TrimSuffix(s, "\r")
					f.tags = append(f.tags, s)
				}
			}
		}*/

		/*if strings.Contains(line, "date:") {
			// если есть доп. символы Parse не работает
			s := mTrimPrefix(line, "date:")
			// s = strings.TrimPrefix(s, " ")
			// s = strings.TrimSuffix(s, "\r")
			// s = strings.TrimSpace(s)
			d, _ := time.Parse("02.01.2006 15:04", s)
			// fmt.Println("v", d.Weekday())
			// fmt.Println("v", d.Year())
			f.date = d
		}*/

	}

	// читаем многострочные
	/*copy := false
	minLen := 3
	for _, line := range text {

		if copy {
			if strings.Contains(line, "_____") {
				copy = false
				break
			}
			lineNext := mTrimPrefix(line, " ")
			if len(lineNext) > minLen {
				f.links = append(f.links, lineNext)
			}
		}

		if strings.Contains(line, "link:") {
			copy = true
			line0 := mTrimPrefix(line, "link:")
			// line0 = strings.TrimPrefix(line0, " ")
			// line0 = strings.TrimSuffix(line0, "\r")
			if len(line0) > minLen {
				f.links = append(f.links, line0)
			}
		}
	}*/
	/*for _, line := range text {

		if copy {
			if strings.Contains(line, "_____") {
				copy = false
				break
			}
			lineNext := mTrimPrefix(line, " ")
			if len(lineNext) > minLen {
				f.bindingFiles = append(f.bindingFiles, lineNext)
			}
		}

		if strings.Contains(line, "bind:") {
			copy = true
			line0 := mTrimPrefix(line, "bind:")
			// line0 = strings.TrimPrefix(line0, " ")
			// line0 = strings.TrimSuffix(line0, "\r")
			if len(line0) > minLen {
				f.bindingFiles = append(f.bindingFiles, line0)
			}
		}
	}*/

	return
}

func mTrimPrefix(s, prefix string) string {

	if strings.HasPrefix(s, prefix) {
		s = strings.TrimPrefix(s, prefix)
	}
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "\r")
	return s
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
			// line = strings.TrimSuffix(line, "\n")
			// line = strings.TrimSuffix(line, "\r")
			fileText += line + "\n"
		}
		if strings.Contains(line, "_____") {
			copy = true
		}
	}

	for strings.HasPrefix(fileText, "\r\n") {
		fileText = strings.TrimPrefix(fileText, "\r\n")
	}
	for strings.HasSuffix(fileText, "\r\n") {
		fileText = strings.TrimSuffix(fileText, "\r\n")
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
	for _, tag := range selectedFile.tags {
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
	bindingMEntry.Wrapping = fyne.TextWrapWord
	bindingMEntry.TextStyle.Monospace = true
	linkMEntry := widget.NewMultiLineEntry()
	linkMEntry.Wrapping = fyne.TextWrapWord
	linkMEntry.TextStyle.Monospace = true

	entryBox := container.NewVBox(
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
	textEntry.Wrapping = fyne.TextWrapWord // обязательно перенос по словам
	textEntry.SetText(text)

	binds := ""
	for _, b := range data.bindingFiles {
		binds += b + "\r\n"
	}
	bindingMEntry.SetText(binds)

	links := ""
	for _, b := range data.links {
		links += b + "\r\n"
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
			_, err := time.Parse("02.01.2006 15:04", dateEntry.Text)
			if err != nil {
				statusLabel.SetText("Неверно заполнено поле даты")
				return
			}
		}

		sl := strings.Split(tagEntry.Text, "#")
		for _, s := range sl {
			if s != "" {
				s = mTrimPrefix(s, " ")
				d.tags = append(d.tags, "#"+s)
			}
		}
		sl = strings.Split(bindingMEntry.Text, "\n")
		for _, s := range sl {
			if s != "" {
				s = mTrimPrefix(s, " ")
				d.bindingFiles = append(d.bindingFiles, s)
			}
		}
		sl = strings.Split(linkMEntry.Text, "\n")
		for _, s := range sl {
			if s != "" {
				s = mTrimPrefix(s, " ")
				d.links = append(d.links, s)
			}
		}
		d.filePath = filepath.Join(filepath.Dir(data.filePath), fileNameEntry.Text)
		d.topic = topicEntry.Text
		d.date, _ = time.Parse("02.01.2006 15:04", dateEntry.Text)

		tt := textEntry.Text
		for strings.HasPrefix(tt, "\r\n") {
			tt = strings.TrimPrefix(tt, "\r\n")
		}
		for strings.HasSuffix(tt, "\r\n") {
			tt = strings.TrimSuffix(tt, "\r\n")
		}

		err := saveFile(d, tt)
		if err == nil {
			statusLabel.SetText("Файл сохранён")
		} else {
			statusLabel.SetText("Ошибка сохранения файла: " + err.Error())
		}
	})

	notSaveButton := widget.NewButton("Закрыть без сохранения", func() {
		// todo: задавать вопрос только если были изменения (bool в поле ввода?)
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

	box := container.NewBorder(entryBox, container.NewBorder(nil, bottomBox, nil, btn), nil, nil, textEntry)
	w.SetContent(box)
	w.Show() // ShowAndRun -- panic!
}

/*
сохраняет данные из структуры в формате ztc
*/
func saveFile(data fileType, text string) error {
	sep := "____________________________________________________________\r\n"

	textall := "topic: " + data.topic + "\r\n"
	textall += "tag:"
	for _, tag := range data.tags {
		textall += " " + tag
	}
	textall += "\r\n" + sep
	textall += text + "\r\n"
	textall += sep

	textall += "link:\r\n"
	for _, link := range data.links {
		textall += " " + link + "\r\n"
	}
	textall += sep

	textall += "bind:\r\n"
	for _, bind := range data.bindingFiles {
		textall += " " + bind + "\r\n"
	}
	textall += sep

	textall += "date: "
	textall += data.date.Format("02.01.2006 15:04") + "\r\n"

	err := os.WriteFile(data.filePath, []byte(textall), 0666)

	return err
}
