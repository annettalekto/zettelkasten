package main

import (
	"fmt"
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
	var list *widget.List
	// selectedDir := "C:\\Users\\nesterovaaa\\Dropbox\\Zettelkasten"
	selectedDir := "C:\\Users\\Totoro\\Dropbox\\Zettelkasten"

	//todo добавить кнопку открывающую список тем и тегов. В какой формат переделать?

	openButton := widget.NewButton("Открыть", func() {
		text := getText(selectedFile.filePath)
		data := fileRead(selectedFile.filePath)
		if data.filePath != "" { // todo не открывать если не выбран файл
			textEditor(data, text)
		}
	})
	createButton := widget.NewButton("Создать", func() {
		var data fileType
		textEditor(data, "")
	})
	bottomBox := container.NewHBox(layout.NewSpacer(), createButton, openButton)

	files, err := os.ReadDir(selectedDir)
	if err != nil {
		fmt.Printf("Ошибка: рабочая папка не открыта\n")
	}
	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
	}
	dirLabel := widget.NewLabel(selectedDir)

	dirButton := widget.NewButton("Каталог", func() {
		dialog := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
			if r != nil && err == nil {
				fmt.Println(r.URI())
				selectedDir = filepath.Dir(r.URI().Path())
				dirLabel.SetText(selectedDir)
				files, _ = os.ReadDir(selectedDir)
				list.Refresh()
			}
		},
			fyne.CurrentApp().Driver().AllWindows()[0],
		)

		// dialog.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
		var dir fyne.ListableURI
		d := storage.NewFileURI(selectedDir)
		dir, _ = storage.ListerForURI(d)

		dialog.SetLocation(dir)
		dialog.SetConfirmText("Выбрать") // вместо ok
		dialog.SetDismissText("Закрыть") // вместо Cancel
		dialog.Show()
	})
	dirBox := container.NewBorder(nil, nil, dirLabel, dirButton)
	// topBox := container.NewVBox(dirBox)

	fileNameEntry := widget.NewEntry()
	fileNameEntry.TextStyle.Monospace = true
	// fileNameSearchButton := widget.NewButton("Поиск", nil)
	topicEntry := widget.NewEntry()
	topicEntry.TextStyle.Monospace = true
	// topicSearchButton := widget.NewButton("Поиск", nil)
	tagEntry := widget.NewEntry()
	tagEntry.TextStyle.Monospace = true
	// tagSearchButton := widget.NewButton("Поиск", nil) // + возможность выбора
	dateEntry := widget.NewEntry()
	dateEntry.TextStyle.Monospace = true
	// dateSearchButton := widget.NewButton("Поиск", nil) // заменить на элемент с датой
	searchButton := widget.NewButton("Поиск", nil)
	clearButton := widget.NewButton("Очистить", nil)

	searchBox := container.NewVBox(
		label(""),
		container.NewBorder(nil, nil, label("Имя:  "), nil, fileNameEntry),
		container.NewBorder(nil, nil, label("Тема: "), nil, topicEntry),
		container.NewBorder(nil, nil, label("Теги: "), nil, tagEntry),
		container.NewBorder(nil, nil, label("Дата: "), nil, dateEntry),
		container.NewBorder(nil, nil, nil, container.NewHBox(clearButton, searchButton)),
	)
	// searchBox2 := container.NewBorder(nil, nil, nil, container.NewHBox(clearButton, searchButton))
	// searchBox := container.NewVBox(searchBox1, searchBox2)

	list = widget.NewList(
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
