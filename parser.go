package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type ztcBasicsType struct {
	filePath     string
	id           int
	title        string
	tags         []string
	sourceNumber int
	source       string
	bindNumbers  []int
	binds        []string
	data         time.Time
	quotation    string
	comment      string
}

type fileType struct { // ztcElementsType ztcBasicsType
	filePath     string // полное имя файла с путем и расширением файла
	topic        string
	tags         []string
	links        []string
	bindingFiles []string
	date         time.Time
}

var selectedFile ztcBasicsType

// todo: переделать на теги
func fileRead2(filePath string) (ztc ztcBasicsType, err error) { //todo: переименовать
	const (
		tagTopic     = "title"
		tagId        = "id"
		tagTags      = "tags"
		tagSource    = "source"
		tagBinds     = "bind"
		tagData      = "date"
		tagQuotation = "quotation"
		tagComment   = "comment"
	)
	temp := ""

	ztc.filePath = filePath

	// заглавие
	temp, err = getElementFromFile(filePath, tagTopic)
	if err != nil {
		fmt.Println(err)
	}
	ztc.title = temp

	// номер карточки
	temp, err = getElementFromFile(filePath, tagId)
	if err != nil {
		fmt.Println(err)
	}
	sub := "_номер:_"
	number := strings.Index(temp, sub)
	temp = temp[number+len(sub):]
	fmt.Println(len(temp))
	temp = strings.TrimSpace(temp)
	fmt.Println(len(temp))

	tempint, err := strconv.Atoi(temp)
	if err != nil {
		fmt.Println(err) // todo: куда err
	}
	ztc.id = tempint

	// теги
	temp, err = getElementFromFile(filePath, tagTags)
	if err != nil {
		fmt.Println(err)
	}
	tempsl := strings.Split(temp, " ")
	fmt.Println(tempsl)
	ztc.tags = tempsl

	// номер и имя карточки источник
	temp, err = getElementFromFile(filePath, tagSource)
	if err != nil {
		fmt.Println(err)
	}
	// _источник:_ 10[[s10 - Фокус. Как сконцентрироваться на главном]]
	sub = "_источник:_"
	number = strings.Index(temp, sub)
	temp = temp[number+len(sub):]
	sub = "[["
	number = strings.Index(temp, sub)
	ztc.source = temp[number:]
	temp = temp[:number]
	temp = strings.TrimSpace(temp)
	tempint, err = strconv.Atoi(temp)
	if err != nil {
		fmt.Println(err)
	}
	ztc.sourceNumber = tempint
	// ztc.source = strings.TrimLeft(ztc.source, "[[") надо выщитывать положение, пробелы
	// ztc.source = strings.TrimRight(ztc.source, "]]")
	fmt.Println(ztc.sourceNumber, ztc.source)

	// temp, err = getElementFromFile(filePath, tagLink)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// // fmt.Println(temp, err) // debug ok
	// f.links = strings.Split(temp, " ")
	// fmt.Println(f.links, err)

	// todo: как разделять метки для связи
	/*temp, err = getElementFromFile(filePath, tagLink)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(temp, err) // debug ok
	f.links = strings.Split(temp, " ")
	fmt.Println(f.links, err)*/

	//--------------

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

func getElementFromFile(filePath, tag string) (s string, err error) {
	var bytes []byte
	bytes, err = os.ReadFile(filePath)
	if err != nil {
		err = fmt.Errorf("getElementFromFile: file read (%s) error: %s", filePath, err.Error())
		return
	}

	lines := strings.Split(string(bytes), "\n")         //note: a new line character in Windows  \r\n
	lines[0], _ = strings.CutPrefix(lines[0], "\ufeff") // cut BOM todo: strings.Split test

	// <!-- title --> #### Риторика <!-- /title -->
	copy := false
	stemp := ""
	for _, line := range lines {
		if strings.Contains(line, "<!--") && strings.Contains(line, "/"+tag) { // закрывающий тэг
			if strings.Count(line, tag) == 2 { // открывающий и закрывающий тег в одной строке
				// искать между тегами одной строки
				sub := "-->"
				number := strings.Index(line, sub)
				line = line[number+len(sub):]
				sub = "<!"
				number = strings.Index(line, sub)
				stemp = line[:number]
				break
			}
			copy = false
			// stemp = strings.TrimSpace(stemp) // последний пробел лишний
			break
		}
		if copy {
			line = strings.TrimSpace(line)
			stemp += line + " "
		}
		if strings.Contains(line, "<!--") && strings.Contains(line, tag) {
			copy = true
		}
	}
	s = stemp

	return
}

// todo: вынести в отдельный файл если еще парачка будет
// сделать ее универсальной, ориетнироваться на < > todo:
//
//	<!-- /title -->
/*func cutTags(tagName, before string) (s string, err error) {
	ok := false
	s = before // не изменять в случае ошибки

	before, ok = strings.CutPrefix(before, "t")
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
}*/

// todo: переделать на теги
func fileRead(filePath string) (f fileType, err error) { //todo: переименовать
	const (
		tagTopic = "topic"
		tagTag   = "tag"
		tagText  = "text"
		tagLink  = "link"
		tagBind  = "bind"
		tagDate  = "date"
	)

	// bytes, err := os.ReadFile(filePath)
	// if err != nil {
	// 	err = fmt.Errorf("file read (%s) error: %w", filePath, err)
	// 	fmt.Println(err)
	// 	return
	// }
	f.filePath = filePath

	// f.topic, err = getElementFromFile(filePath, tagTopic)
	// if err != nil {
	// 	fmt.Println(err)
	// 	// todo: можно в лог писать кстати, а критичные еще в статус строке
	// }
	// fmt.Println(f.topic, err) // debug ok

	temp := ""
	temp, err = getElementFromFile(filePath, tagTag)
	if err != nil {
		fmt.Println(err)
	}
	f.tags = strings.Split(temp, " ")
	fmt.Println(f.tags, err) // debug ok
	fmt.Println(temp, err)   // debug ok

	temp, err = getElementFromFile(filePath, tagLink)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(temp, err) // debug ok
	f.links = strings.Split(temp, " ")
	fmt.Println(f.links, err)

	// todo: как разделять метки для связи
	/*temp, err = getElementFromFile(filePath, tagLink)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(temp, err) // debug ok
	f.links = strings.Split(temp, " ")
	fmt.Println(f.links, err)*/

	//--------------

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
		container.NewBorder(nil, nil, newFormatLabel("Имя:     "), nil, fileNameEntry),
		container.NewBorder(nil, nil, newFormatLabel("Тема:    "), nil, topicEntry),
		container.NewBorder(nil, nil, newFormatLabel("Теги:    "), nil, tagEntry),
		container.NewBorder(nil, nil, newFormatLabel("Дата:    "), nil, dateEntry),
		container.NewBorder(nil, nil, newFormatLabel("Связать: "), nil, bindingMEntry),
		container.NewBorder(nil, nil, newFormatLabel("Cсылки:  "), nil, linkMEntry),
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
