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
	w.Resize(fyne.NewSize(900, 700))
	w.CenterOnScreen()
	w.SetMaster()

	// регистрируем падение скрипта
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("!!!		Пааааника %c%c%c		!!!", 128561, 128561, 128561)
			os.Exit(1)
		}
	}()

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

	// NOTE:
	// прога для работы с сиситемой Zettelkasten
	// будет сохранять файлы в определенном виде, ну и читать их
	// Открыть: выбранный из списка файл в редакторе с возможностью сохранения
	// Создать: ввод текста и доп. данных, сохранение в формате, запись в список тегов и тд

	// Поиск попробоавать организовать так:
	// возможность поиска из тех что уже найдены или вернуться к поиску по всему каталогу (галочкой)

	// сортировака по имени и дате (теме?)

	// теги - есть файл со списком тегов, желательно в каком то виде вывести этот файл, чтобы
	// можно было создавать новые теги (то же темы и др). а так же можно сравнивать с этим список
	// когда создается файл и указываются теги.
	// в поиске можно использовать tagSelectEntry с выбором всех имеющихся, но при выводе файла  нужно отображать все

	// ?решить проблемку \r \n (либо разбивку на строки делать както иначе, либо зачищать от символов, добавлять свой перевод строки в конце)
	// убрать txt из названия (добавить расширение как настройку (md, txt др)? - или лишнее усложнение?)
	// добавить обновление списка файлов (по кнопке?) после создания файла хотя бы.. лучше в потоке

	// теги то с большой то с маленькой буквы - привести к общему виду
	// перевести доп файлы на .csv?

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
	img.Resize(fyne.NewSize(66, 90))
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

func mainForm() (box *fyne.Container) {
	var list *widget.List
	statusLabel := widget.NewLabel("Тут что-нибудь отладочное...")
	selectedDir := "C:\\Users\\Totoro\\Dropbox\\Zettelkasten"

	// кнопки
	openButton := widget.NewButton("Открыть", func() {
		text := getText(selectedFile.filePath)
		data := fileRead(selectedFile.filePath)
		if data.filePath != "" { //todo: ???
			textEditor(data, text)
		}
	})
	createButton := widget.NewButton("Создать", func() {
		var data fileType
		data.date = time.Now()
		data.filePath = filepath.Join(selectedDir, "new")
		textEditor(data, "")
	})

	files, err := os.ReadDir(selectedDir)
	if err != nil {
		fmt.Printf("Ошибка: рабочая папка не открыта\n") // TODO: как обрабатывать ошибки
		statusLabel.Text = "Ошибка: рабочая папка не открыта"
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

		// note: использовать если нужен выбор разных файлов + нужно считывать папку дополнительно
		// dialog.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
		var dir fyne.ListableURI
		d := storage.NewFileURI(selectedDir)
		dir, _ = storage.ListerForURI(d)

		dialog.SetLocation(dir)
		dialog.SetConfirmText("Выбрать")
		dialog.SetDismissText("Закрыть")
		dialog.Show()
	})
	dirBox := container.NewBorder(nil, nil, dirLabel, dirButton)

	// поиск
	const (
		topic = "Тема" //переименовать
		tag   = "Тег"
	)
	searchLabel := widget.NewLabel("подсказка как искать")
	searchSelect := widget.NewSelect([]string{topic, tag}, func(value string) { // todo: дата
		if value == topic {
			searchLabel.SetText("Введите слова, которые должна содержать тема")
		} else if value == tag {
			searchLabel.SetText("Введите один тег без знака #")
		}
	})
	searchSelect.SetSelected(topic)
	searchEntry := widget.NewEntry()
	searchEntry.TextStyle.Monospace = true
	searchButton := widget.NewButton("  Поиск  ", nil)
	clearButton := widget.NewButton("Очистить", func() {
		searchEntry.SetText("")
	})
	check := widget.NewCheck("Поиск по всей папке", func(b bool) {
	})
	searchButtonBox := container.NewBorder(nil, nil, check, searchButton)
	searchBox := container.NewVBox(widget.NewLabel(""), container.NewBorder(nil, nil, searchSelect, clearButton, searchEntry), searchLabel, searchButtonBox)

	topBottom := container.NewVBox(dirBox, searchBox)

	// краткое отображение карточки
	fileNameEntry := widget.NewEntry()
	fileNameEntry.TextStyle.Monospace = true
	topicEntry := widget.NewEntry()
	topicEntry.TextStyle.Monospace = true
	tagEntry := widget.NewEntry()
	tagEntry.TextStyle.Monospace = true
	dateEntry := widget.NewEntry()
	dateEntry.TextStyle.Monospace = true

	entryBox := container.NewVBox(
		newlabel(""),
		container.NewBorder(nil, nil, newlabel("Имя:  "), nil, fileNameEntry),
		container.NewBorder(nil, nil, newlabel("Тема: "), nil, topicEntry),
		container.NewBorder(nil, nil, newlabel("Теги: "), nil, tagEntry),
		container.NewBorder(nil, nil, newlabel("Дата: "), nil, dateEntry),
		container.NewHBox(createButton, layout.NewSpacer(), openButton),
	)

	// список
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
		for _, tag := range selectedFile.tags {
			tags += tag + " "
		}
		tagEntry.SetText(tags)
		fileNameEntry.SetText(files[id].Name())
		topicEntry.SetText(selectedFile.topic)
		dateEntry.SetText(selectedFile.date.Format("02.01.2006 15:04"))
	}

	panelBox := container.NewBorder(topBottom, nil, nil, nil, entryBox)
	split := container.NewHSplit(list, panelBox)
	box = container.NewBorder(nil, statusLabel, nil, nil, split)

	return
}
