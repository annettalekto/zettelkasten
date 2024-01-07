package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type viewType struct {
	Text *widget.Entry
	Date *widget.Label
	Name *widget.Entry
}

var viewForm viewType

func (v *viewType) initForm() *fyne.Container {

	// номер карты в шапке?
	v.Date = newFormatLabel(fmt.Sprintf("%v", selectedFile.data))
	d := container.NewBorder(nil, nil, nil, v.Date)
	v.Name = newFormatEntry()
	v.Name.SetText("Имя файла")
	top := container.NewVBox(d, v.Name)

	v.Text = widget.NewMultiLineEntry()
	v.Text.TextStyle.Monospace = true
	v.Text.Wrapping = fyne.TextWrapWord
	v.Text.SetText("Текст")

	bottom := container.NewBorder(nil, nil, nil, widget.NewButton("Редакт.", nil))

	return container.NewBorder(top, bottom, nil, nil, v.Text)
}

type addInfoType struct {
	Name   *widget.Entry
	Tegs   *widget.Entry
	Binds  *widget.Entry
	Source *widget.Entry
}

var addInfoForm addInfoType

func (a *addInfoType) initForm() *fyne.Container {

	a.Name = newFormatEntry()
	a.Name.SetText("Имя файла")

	a.Tegs = widget.NewMultiLineEntry()
	a.Tegs.SetText("#тег1\n#тег2")
	a.Binds = widget.NewMultiLineEntry()
	a.Binds.SetText("Связное")
	a.Source = widget.NewMultiLineEntry()
	a.Source.SetText("Источники")
	box := container.NewGridWithColumns(1, a.Tegs, a.Binds, a.Source)

	return container.NewBorder(a.Name, nil, nil, nil, box)
}

type sourceInfoType struct {
	Source    *widget.Entry
	Quotation *widget.Entry
}

var sourceInfoForm sourceInfoType

func (s *sourceInfoType) initForm() *fyne.Container {
	s.Source = widget.NewMultiLineEntry()
	s.Source.SetText("Источник")

	s.Quotation = widget.NewMultiLineEntry()
	s.Quotation.SetText("Цитата")

	return container.NewGridWithColumns(1, s.Source, s.Quotation) // todo: разделить
}

type commentType struct {
	Comment *widget.Entry
}

var commentForm commentType

func (c *commentType) initForm() *fyne.Container {
	c.Comment = widget.NewMultiLineEntry()
	c.Comment.SetText("Комментарий")

	return container.NewGridWithColumns(1, c.Comment)
}

func (z *ztcBasicsType) refreshTabs() {
	viewForm.Date.SetText(fmt.Sprintf("%v", z.data.Format("2006-01-02 15:04"))) //d.Format("2006-01-02 15:04")
	viewForm.Name.SetText(z.title)
	viewForm.Text.SetText(getTextFromFile(z.filePath))
	addInfoForm.Name.SetText(z.title)
	addInfoForm.Tegs.SetText(fmt.Sprintf("%q", z.tags))
	addInfoForm.Binds.SetText(fmt.Sprintf("%q", z.bind))
	addInfoForm.Source.SetText(fmt.Sprintf("%q", z.source))
	sourceInfoForm.Source.SetText(fmt.Sprintf("%q", z.source))
	sourceInfoForm.Quotation.SetText(getQuotationFromFile(z.filePath))
	commentForm.Comment.SetText(getCommentFromFile(z.filePath))
}
