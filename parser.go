package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type ztcBasicsType struct {
	filePath     string
	id           string // todo: а зачем переводить в числа?
	title        string
	tags         []string
	sourceNumber []string // может ли быть несколько источников?
	source       []string
	bindNumbers  []string
	bind         []string
	data         time.Time
	// quotation    string
	// comment      string
}

// type fileType struct { // ztcElementsType ztcBasicsType
// 	filePath     string // полное имя файла с путем и расширением файла
// 	topic        string
// 	tags         []string
// 	links        []string
// 	bindingFiles []string
// 	date         time.Time
// }

var selectedFile ztcBasicsType

func fileRead(filePath string) (ztc ztcBasicsType, err error) {

	ztc.filePath = filePath
	ztc.title = getTopicFromFile(filePath)
	ztc.id = getCardIdFromFile(filePath)
	ztc.tags = getTagsFromFile(filePath)
	ztc.sourceNumber, ztc.source = getSourceFromFile(filePath)
	ztc.bindNumbers, ztc.bind = getBindFromFile(filePath)
	ztc.data = getDataFromFile(filePath)
	// ztc.quotation = getQuotationFromFile(filePath)
	// ztc.comment = getCommentFromFile(filePath)

	return
}

func getTopicFromFile(filePath string) (s string) {
	const tagTopic = "title"
	var err error

	s, err = getElementFromFile(filePath, tagTopic)
	if err != nil {
		fmt.Println(err) // todo: куда выводить ошибки и нудо ли?
	}
	return
}

func getCardIdFromFile(filePath string) (s string) {
	const tagId = "id"
	var err error

	s, err = getElementFromFile(filePath, tagId)
	if err != nil {
		fmt.Println(err)
	}
	sl := getAllNumbers(s)
	fmt.Printf("%q", sl)
	return sl[0]
}

func getTagsFromFile(filePath string) (sl []string) {
	const tagTags = "tags"
	var err error

	s, err := getElementFromFile(filePath, tagTags)
	if err != nil {
		fmt.Println(err)
	}
	sl = strings.Split(s, " ")
	fmt.Printf("%q", sl)
	return
}

func getSourceFromFile(filePath string) (numbers, names []string) {
	const tagSource = "source"

	s, err := getElementFromFile(filePath, tagSource)
	if err != nil {
		fmt.Println(err)
	}
	numbers = getAllNumbers(s[:strings.Index(s, "[[")])
	names = removeSquareBrackets(s)
	fmt.Printf("%q %q", numbers, names)
	return
}

func getBindFromFile(filePath string) (numbers, names []string) {
	const tagBinds = "bind"
	s, err := getElementFromFile(filePath, tagBinds)
	if err != nil {
		fmt.Println(err)
	}
	// _связное:_ 7, 1 [[7 - Изучение языков]] [[1 - Смысл]]
	names = removeSquareBrackets(s)                     //
	numbers = getAllNumbers(s[:strings.Index(s, "[[")]) //_связное:_ 7, 1
	fmt.Printf("%q %q", numbers, names)
	return
}

func getDataFromFile(filePath string) (t time.Time) {
	const tagData = "date"

	s, err := getElementFromFile(filePath, tagData)
	if err != nil {
		fmt.Println(err)
	}
	s = strings.TrimSpace(s)
	t, err = time.Parse("2006-01-02 15:04", s) // todo: если будет хоть один лишний пробел по бокам то ошибка
	if err != nil {
		fmt.Println("error") // todo: вопрос куда ошибки и в какой форме
	}
	fmt.Println(t)

	return
}

func getTextFromFile(filePath string) (s string) {
	const tagText = "text"
	var err error

	s, err = getElementFromFile(filePath, tagText)
	if err != nil {
		fmt.Println("error")
	}

	return
}

func getQuotationFromFile(filePath string) (s string) {
	const tagQuotation = "quotation"
	var err error

	s, err = getElementFromFile(filePath, tagQuotation)
	if err != nil {
		fmt.Println("error")
	}

	return
}

func getCommentFromFile(filePath string) (s string) {
	const tagComment = "comment"
	var err error

	s, err = getElementFromFile(filePath, tagComment)
	if err != nil {
		fmt.Println("error")
	}

	return
}

// выбрать все числа
func getAllNumbers(s string) []string {
	var temp []rune

	for _, r := range s[:] {
		switch r {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', ',':
			temp = append(temp, r)
		}
	}
	str := string(temp)

	return strings.Split(str, ",")
}

// вынуть то что в квадратных скобках
func removeSquareBrackets(s string) (ss []string) {

	f := func() (elem string) {
		sub := "[["
		i := strings.Index(s, sub)
		s = s[i+2:]

		sub = "]]"
		i = strings.Index(s, sub)
		elem = s[:i]
		s = s[i+2:] // обрезать строку

		return elem
	}

	for strings.Count(s, "[[") > 0 {
		ss = append(ss, f())
	}

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
	s = strings.TrimSpace(stemp)

	return
}

/*
func mTrimPrefix(s, prefix string) string {

	if strings.HasPrefix(s, prefix) {
		s = strings.TrimPrefix(s, prefix)
	}
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "\r")
	return s
}*/

/*
getText - прочитать файл,
вернуть только текст заметки
*/
/*
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
}*/

/*
textEditor - открыть окно с текстом выбранного файла
Сохранить изменения
Закрыть без сохранения
*/
/*
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
}*/

/*
сохраняет данные из структуры в формате ztc
*/
/*func saveFile(data fileType, text string) error {
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
}*/
