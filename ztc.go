package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Zettelkasten")
	w.Resize(fyne.NewSize(700, 500))

	// прога для работы с сиситемой Zettelkasten
	// будет сохранять файлы в определенном виде, ну и читать их
	// тест git

	w.SetContent(form())
	w.ShowAndRun()
}

func face() *fyne.Container {
	serchEntry := widget.NewEntry()
	btnOpen := widget.NewButton("Открыть", nil)
	btnSave := widget.NewButton("Сохранить", nil)

	btnBox := container.NewHBox(btnOpen, btnSave)
	box := container.NewVBox(serchEntry, btnBox)
	return box
}

func form() *fyne.Container {
	titleEntry := widget.NewEntry()
	titleEntry.PlaceHolder = "Заголовок"
	// titleEntry.OnChanged
	hash := widget.NewEntry()
	hash.PlaceHolder = "#хеш"
	boxTop := container.NewHBox(titleEntry, hash)

	textEntry := widget.NewEntry()
	boxEntry := container.NewMax(textEntry)

	link := widget.NewEntry()
	link.PlaceHolder = "Ссылка"

	btnSave := widget.NewButton("Сохранить", nil)
	boxBottom := container.NewHBox(link, btnSave)

	box := container.NewHBox(boxTop, boxEntry, boxBottom)
	return box
}
