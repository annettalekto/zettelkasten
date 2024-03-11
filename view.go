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

type elmFormType struct {
	id           *numericalEntry
	FileName     *widget.Entry
	Title        *widget.Entry
	Tags         *widget.Entry
	BindNumbers  *numericalEntry
	Binds        *widget.Entry
	SourceNumber *numericalEntry
	Source       *widget.Entry
	Text         *widget.Entry
	Date         *widget.Entry
	Quotation    *widget.Entry
	Comment      *widget.Entry
}

var view elmFormType

func (e *elmFormType) createCard() {
	file, err := os.Create(filepath.Join(gFilePath, e.FileName.Text+".md"))
	if err != nil {
		fmt.Println("Ошибка записи файла")
		return
	}
	defer file.Close()

	file.WriteString("<!-- title -->\n" + "#### " + e.Title.Text + "\n<!-- /title -->\n\n")
	file.WriteString("<!-- tags -->\n" + e.Tags.Text + "\n<!-- /tags -->\n\n")
	file.WriteString("<!-- text -->\n" + e.Text.Text + "\n<!-- /text -->\n\n")

	source := strings.Split(e.Source.Text, "\n")
	file.WriteString("<!-- source -->\n" + "_источник:_ " + e.SourceNumber.Text + "\n")
	for _, s := range source {
		file.WriteString("[[" + s + "]]\n")
	}
	file.WriteString("\n<!-- /source -->\n\n")

	binds := strings.Split(e.Binds.Text, "\n")
	file.WriteString("<!-- bind-->\n" + "_связное:_  " + e.BindNumbers.Text + "\n")
	for _, b := range binds {
		file.WriteString("[[" + b + "]]\n")
	}
	file.WriteString("\n<!-- /bind -->\n\n")

	file.WriteString("<!-- id -->" + " _номер:_ " + e.id.Text + " <!-- /id -->\n")
	file.WriteString("<!-- date --> " + e.Date.Text + " <!-- /date -->\n")
	file.WriteString("\n___\n\n")

	file.WriteString("<!-- quotation -->\n" + e.Quotation.Text + "\n<!-- /quotation -->\n\n")
	file.WriteString("<!-- comment -->\n" + e.Comment.Text + "\n<!-- /comment -->\n")

	// ztc.id = e.id.Text
	// ztc.title =
	// ztc.filePath = e.FileName.Text
	// ztc.tags = strings.Split(e.Tags.Text, "\n")
	// ztc.bind = strings.Split(e.Binds.Text, "\n")
	// ztc.bindNumbers = strings.Split(e.BindNumbers.Text, ", ")
	// ztc.source = strings.Split(e.Source.Text, "\n")
	// ztc.sourceNumber = strings.Split(e.SourceNumber.Text, ", ")
	// ztc.data, err = time.Parse("2006-01-02 15:04", e.Date.Text)
	// if err != nil {
	// 	fmt.Println("ошибочка")
	// }
	// text
	// quotation
	// comment
}

func (e *elmFormType) nameCardForm() *fyne.Container {
	e.id = newNumericalEntry()
	ent1 := container.NewBorder(nil, nil, widget.NewLabel("Номер карты:"), nil, e.id)
	e.id.Entry.OnChanged = func(s string) {
		e.FileName.SetText(fmt.Sprintf("%s - %s", s, e.Title.Text))
	}
	e.Title = newFormatEntry()
	title := container.NewBorder(nil, nil, widget.NewLabel("Название:"), nil, e.Title)
	e.Title.OnChanged = func(s string) {
		e.FileName.SetText(fmt.Sprintf("%s - %s", e.id.Text, s))
	}
	e.FileName = widget.NewEntry() // без папки
	ent2 := container.NewBorder(nil, nil, widget.NewLabel("Имя файла:"), nil, e.FileName)
	top := container.NewVBox(ent1, title, ent2)

	e.Text = newText()

	e.Date = newFormatEntry()
	e.Date.SetText(fmt.Sprintf("%v", time.Now().Format("2006-01-02 15:04")))

	CreateButton := widget.NewButton("Создать карточку", func() {
		e.createCard() // true?
	})
	bottom := container.NewBorder(nil, nil, nil, CreateButton, e.Date)

	return container.NewBorder(top, bottom, nil, nil, e.Text)
}

func (e *elmFormType) textForm() *fyne.Container {

	e.Date = newFormatEntry()
	e.Date.SetText(fmt.Sprintf("%v", selectedFile.data))
	e.Title = newFormatEntry()
	e.Text = newText()
	// e.Text.SetText("<Текст>")

	return container.NewBorder(
		e.Title,
		e.Date,
		nil,
		nil,
		e.Text)
}

func (e *elmFormType) addInfoForm() *fyne.Container {

	e.BindNumbers = newNumericalEntry()
	e.Tags = widget.NewMultiLineEntry()
	e.Binds = widget.NewMultiLineEntry()

	binsNumber := container.NewBorder(nil, nil, widget.NewLabel("Номера связных карт:"), nil, e.BindNumbers)

	ent1 := container.NewBorder(nil, nil, nil, nil,
		container.NewBorder(widget.NewLabel("Теги:"), nil, nil, nil, e.Tags))
	e.Tags.SetText("#")
	e.Tags.OnChanged = func(s string) {
		if s[len(s)-1] == '\n' {
			e.Tags.SetText(s + "#")
		}
	}

	ent2 := container.NewBorder(nil, nil, nil, nil,
		container.NewBorder(widget.NewLabel("Связное:"), nil, nil, nil, e.Binds))

	return container.NewBorder(binsNumber, nil, nil, nil, container.NewGridWithColumns(1, ent2, ent1))
}

func (e *elmFormType) sourceForm() *fyne.Container {
	e.SourceNumber = newNumericalEntry()
	e.Source = newText()
	e.Quotation = newText()

	sourceNumber := container.NewBorder(nil, nil, widget.NewLabel("Номер карты источника:"), nil, e.SourceNumber)
	// todo: найти по номеру карту и подставить имя в поле ниже
	ent1 := container.NewBorder(nil, nil, nil, nil,
		container.NewBorder(widget.NewLabel("Источники:"), nil, nil, nil, e.Source))

	ent2 := container.NewBorder(nil, nil, nil, nil,
		container.NewBorder(widget.NewLabel("Цитата:"), nil, nil, nil, e.Quotation))

	return container.NewBorder(sourceNumber, nil, nil, nil, container.NewGridWithColumns(1, ent1, ent2))
}

func (e *elmFormType) commentForm() *fyne.Container {
	e.Comment = newText()

	return container.NewGridWithColumns(1,
		container.NewBorder(widget.NewLabel("Комментарий:"), nil, nil, nil, e.Comment))
}

func (e *elmFormType) refreshTabs(z ztcBasicsType) {

	e.Date.SetText(fmt.Sprintf("%v", z.data.Format("2006-01-02 15:04")))
	e.Title.SetText(z.title)
	e.Tags.SetText(sliceInColumn(z.tags))
	e.Binds.SetText(sliceInColumn(z.bind))
	e.BindNumbers.SetText(sliceInString(z.bindNumbers))
	e.SourceNumber.SetText(sliceInString(z.sourceNumber))
	e.Source.SetText(sliceInColumn(z.source))
	e.Text.SetText(getTextFromFile(z.filePath))
	e.Quotation.SetText(getQuotationFromFile(z.filePath))
	e.Comment.SetText(getCommentFromFile(z.filePath))
}

/*
func (e *elmFormType) editForm() { // открывать только для текущего файла
	w := fyne.CurrentApp().NewWindow("0")
	w.Resize(fyne.NewSize(800, 600))
	w.CenterOnScreen()

	// todo: оформить по типу шаблонного файла
	tagBox := container.NewBorder(widget.NewLabel("Теги:"), nil, nil, nil, e.Tags)
	dindsBox := container.NewBorder(widget.NewLabel("Связное:"), nil, nil, nil, e.Binds)
	// номер и название источников отдельно? как удобнее считывать и составлять итоговый файл
	sourceBox := container.NewBorder(widget.NewLabel("Источники:"), nil, nil, nil, e.Source)
	quotationBox := container.NewBorder(widget.NewLabel("Цитата:"), nil, nil, nil, e.Quotation)
	commentBox := container.NewBorder(widget.NewLabel("Комментарий:"), nil, nil, nil, e.Comment)

	bottom := container.NewGridWithColumns(1, tagBox, dindsBox, sourceBox, quotationBox, commentBox)
	box := container.NewBorder(e.Name, bottom, nil, nil, e.Text)

	w.SetContent(box)
	w.Show()
}
*/
