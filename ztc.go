package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var gVersion, gYear, gProgramName string

func main() {
	gProgramName = "Zettelkasten"

	a := app.New()
	w := a.NewWindow(gProgramName)
	w.Resize(fyne.NewSize(800, 600))
	w.CenterOnScreen()
	w.SetMaster()

	menu := fyne.NewMainMenu(
		fyne.NewMenu("Файл",
			// a quit item will be appended to our first menu
			fyne.NewMenuItem("Тема", func() { changeTheme(a) }),
			// fyne.NewMenuItem("Выход", func() { a.Quit() }),
		),

		fyne.NewMenu("Справка",
			fyne.NewMenuItem("Посмотреть справку", func() { aboutHelp() }),
			// fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("О программе", func() { abautProgramm() }),
		),
	)
	w.SetMainMenu(menu)

	go func() { // простите
		time.Sleep(1 * time.Second)
		for _, item := range menu.Items[0].Items {
			if item.Label == "Quit" {
				item.Label = "Выход"
			}
		}
	}()

	// прога для работы с сиситемой Zettelkasten
	// будет сохранять файлы в определенном виде, ну и читать их
	// Открыть: выбранный из списка файл в редакторе с возможностью сохранения
	// Создать: ввод текста и доп. данных, сохранение в формате, запись в список тегов и тд

	w.SetContent(mainForm())
	w.ShowAndRun()
}

var currentTheme bool // светлая тема false

func changeTheme(a fyne.App) {
	currentTheme = !currentTheme

	if currentTheme {
		a.Settings().SetTheme(theme.DarkTheme())
	} else {
		a.Settings().SetTheme(theme.LightTheme())
	}
}

func aboutHelp() {
	err := exec.Command("cmd", "/C", ".\\help\\index.html").Run()
	if err != nil {
		fmt.Println("Ошибка открытия файла справки")
	}
}

func abautProgramm() {
	w := fyne.CurrentApp().NewWindow("О программе") // CurrentApp!
	w.Resize(fyne.NewSize(400, 150))
	w.SetFixedSize(true)
	w.CenterOnScreen()

	img := canvas.NewImageFromURI(storage.NewFileURI("ind.png"))
	img.Resize(fyne.NewSize(66, 90)) //без изменений
	img.Move(fyne.NewPos(10, 10))

	l0 := widget.NewLabel(gProgramName)
	l0.Move(fyne.NewPos(80, 10))
	l1 := widget.NewLabel(fmt.Sprintf("Версия %s", gVersion))
	l1.Move(fyne.NewPos(80, 40))
	l2 := widget.NewLabel(fmt.Sprintf("© Anna, %s", gYear))
	l2.Move(fyne.NewPos(80, 70))

	box := container.NewWithoutLayout(img, l0, l1, l2)

	// w.SetContent(fyne.NewContainerWithLayout(layout.NewCenterLayout(), box))
	w.SetContent(box)
	w.Show() // ShowAndRun -- panic!
}

func mainForm() (box *container.Split) {
	selectedDir := "C:\\Users\\nesterovaaa\\Dropbox\\Zettelkasten"
	// selectedDir := "C:\\Users\\Totoro\\Dropbox\\Zettelkasten"

	openButton := widget.NewButton("Открыть", func() {
		text := getText(filepath.Join(selectedDir, selectedFile.fileName))
		data := fileRead(filepath.Join(selectedDir, selectedFile.fileName))
		textEditor(data, text)
	})
	createButton := widget.NewButton("Создать", nil)
	bottomBox := container.NewHBox(layout.NewSpacer(), openButton, createButton)

	files, err := os.ReadDir(selectedDir)
	if err != nil {
		fmt.Printf("Ошибка: рабочая папка не открыта\n")
		log.Fatal(err)
	}
	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
	}
	dirLabel := widget.NewLabel(selectedDir)
	dirButton := widget.NewButton("Каталог", func() {
		// как настроить на каталог?
		// DEBUG если выбрано Закрыть, уходит в panic
		openFile := func(r fyne.URIReadCloser, _ error) {
			fmt.Println(r.URI())
		}
		w := fyne.CurrentApp().Driver().AllWindows()[0]

		dialog := dialog.NewFileOpen(openFile, w)
		// dialog.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
		dialog.Show()
	})
	dirBox := container.NewBorder(nil, nil, dirLabel, dirButton)
	// topBox := container.NewVBox(dirBox)

	fileNameEntry := widget.NewEntry()
	fileNameEntry.TextStyle.Monospace = true
	fileNameSearchButton := widget.NewButton("Поиск", nil)
	topicEntry := widget.NewEntry()
	topicEntry.TextStyle.Monospace = true
	topicSearchButton := widget.NewButton("Поиск", nil)
	tagEntry := widget.NewEntry()
	tagEntry.TextStyle.Monospace = true
	tagSearchButton := widget.NewButton("Поиск", nil) // + возможность выбора
	dateEntry := widget.NewEntry()
	dateEntry.TextStyle.Monospace = true
	dateSearchButton := widget.NewButton("Поиск", nil) // заменить на элемент с датой

	searchBox := container.NewVBox(
		label(""),
		container.NewBorder(nil, nil, label("Имя:  "), fileNameSearchButton, fileNameEntry),
		container.NewBorder(nil, nil, label("Тема: "), topicSearchButton, topicEntry),
		container.NewBorder(nil, nil, label("Теги: "), tagSearchButton, tagEntry),
		container.NewBorder(nil, nil, label("Дата: "), dateSearchButton, dateEntry),
	)

	list := widget.NewList(
		func() int {
			return len(files)
		},
		func() fyne.CanvasObject {
			var style fyne.TextStyle
			style.Monospace = true
			temp := widget.NewLabelWithStyle("temp", fyne.TextAlignLeading, style)
			return temp
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i < len(files) {
				o.(*widget.Label).SetText(files[i].Name())
			}
		})

	list.OnSelected = func(id widget.ListItemID) {

		filePath := filepath.Join(selectedDir, files[id].Name())
		selectedFile = fileRead(filePath)

		tags := ""
		for _, tag := range selectedFile.tag {
			tags += "#" + tag + " "
		}
		tagEntry.SetText(tags)
		fileNameEntry.SetText(files[id].Name())
		topicEntry.SetText(selectedFile.topic)
		dateEntry.SetText(selectedFile.date.Format("02.01.2006 15:04"))
	}

	// panelBox := container.NewBorder(container.NewVBox(dirBox, searchBox), bottomBox, nil, nil)
	// panelBox := container.NewBorder(dirBox, bottomBox, nil, nil, container.NewVBox(searchBox))
	panelBox := container.NewBorder(dirBox, bottomBox, nil, nil, searchBox)

	box = container.NewHSplit(list, panelBox)

	return
}

//--------------------------------------
/*
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
	boxTop := container.NewVBox(titleEntry, hash)

	textEntry := widget.NewEntry()
	boxEntry := container.NewMax(textEntry)

	link := widget.NewEntry()
	link.PlaceHolder = "Ссылка"

	btnSave := widget.NewButton("Сохранить", nil)
	boxBottom := container.NewVBox(link, btnSave)

	box := container.NewHBox(boxTop, boxEntry, boxBottom)
	return box
}
*/
