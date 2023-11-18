package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
	gYear = "2023"
	gVersion = "2" // сохр. в файл todo:

	a := app.New()
	w := a.NewWindow(gProgramName)
	w.Resize(fyne.NewSize(600, 400))
	w.CenterOnScreen()
	w.SetMaster()

	// log в файл позже. каждый день новый файл? Папку создать ДО
	// f, err := os.OpenFile(".\\logfiles\\info.log", os.O_RDWR|os.O_CREATE, 0666)
	// if err != nil {
	// 	err = fmt.Errorf("ERROR: log file %w", err)
	// 	log.Fatal(err)
	// }
	// defer f.Close()
	// пока лог в стд, так удобнее
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	slog.Info("Start program", slog.String("version", runtime.Version()))

	// регистрируем падение скрипта
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("!!!		Пааааника %c%c%c		!!!", 128561, 128561, 128561)
			slog.Error("Panic")
			os.Exit(1)
		}
	}()

	menu := fyne.NewMainMenu(
		fyne.NewMenu("Файл",
			fyne.NewMenuItem("Создать новую карточку", nil),
			fyne.NewMenuItem("Изменить каталог", nil),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Список тегов", nil),
			fyne.NewMenuItem("Список оглавление", nil),
			fyne.NewMenuItem("Список литературы", nil),
		),
		fyne.NewMenu("Вид",
			fyne.NewMenuItem("Тема", func() { changeTheme(a) }),
			fyne.NewMenuItem("Сортировка", nil),
		),

		fyne.NewMenu("Справка",
			fyne.NewMenuItem("Посмотреть справку", func() { aboutHelp() }),
			fyne.NewMenuItem("О программе", func() { aboutProgram() }),
		),
	)
	w.SetMainMenu(menu)

	tabs := container.NewAppTabs(
		container.NewTabItem("вариант 1", mainForm1()),
		container.NewTabItem("вариант 2", mainForm()),
		container.NewTabItem("просмотр", viewForm()),
		container.NewTabItem("доп.", addInfoForm()),
	)
	tabs.SetTabLocation(container.TabLocationBottom)

	go func() { // простите
		time.Sleep(2 * time.Second)
		for _, item := range menu.Items[0].Items {
			if item.Label == "Quit" {
				item.Label = "Выход"
			}
		}
	}()

	w.SetContent(tabs)
	w.ShowAndRun()
}

func mainForm() (box *fyne.Container) {

	var list *widget.List
	statusLabel := widget.NewLabel(" ")
	selectedDir := "D:\\ztc test"

	files, err := os.ReadDir(selectedDir)
	if err != nil {
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
	topBottom := container.NewVBox(dirBox) //?

	text := widget.NewMultiLineEntry()
	text.SetText("text")

	openButton := widget.NewButton("Открыть", func() {
		text := getText(selectedFile.filePath)
		data, _ := fileRead(selectedFile.filePath)
		if data.filePath != "" {
			textEditor(data, text)
		}
	})
	createButton := widget.NewButton("Создать", func() {
		var data fileType
		data.date = time.Now()
		data.filePath = filepath.Join(selectedDir, "new")
		textEditor(data, "")
	})
	topicEntry := widget.NewEntry()
	topicEntry.TextStyle.Monospace = true
	entr := container.NewBorder(nil, nil, newFormatLabel("Topic:"), nil, topicEntry)

	btn := container.NewHBox(createButton, layout.NewSpacer(), openButton)
	bottom := container.NewVBox(entr, btn)

	entryBox := container.NewBorder(nil, bottom, nil, nil, text)

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
		selectedFile, err = fileRead(filePath)

		tags := ""
		for _, tag := range selectedFile.tags {
			tags += tag + " "
		}
		topicEntry.SetText(selectedFile.topic)
	}

	panelBox := container.NewBorder(topBottom, nil, nil, nil, entryBox)
	split := container.NewHSplit(list, panelBox)
	box = container.NewBorder(nil, statusLabel, nil, nil, split)

	return
}

var darkTheme bool

func changeTheme(a fyne.App) {
	darkTheme = !darkTheme

	if darkTheme {
		a.Settings().SetTheme(&forcedVariant{Theme: theme.DefaultTheme(), variant: theme.VariantDark})
	} else {
		a.Settings().SetTheme(&forcedVariant{Theme: theme.DefaultTheme(), variant: theme.VariantLight})
	}
}

func aboutHelp() {
	err := exec.Command("cmd", "/C", ".\\help\\index.html").Run()
	if err != nil {
		fmt.Println("Ошибка открытия файла справки")
	}
}

func aboutProgram() {
	w := fyne.CurrentApp().NewWindow("О программе")
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

	w.SetContent(box)
	w.Show()
}

// ------------------------------------------------------------------------------------------
// old
func mainForm1() (box *fyne.Container) {
	var list *widget.List
	statusLabel := widget.NewLabel("Тут что-нибудь отладочное...")
	// selectedDir := "C:\\Users\\Totoro\\Dropbox\\Zettelkasten"
	selectedDir := "C:\\Users\\Totoro\\Dropbox\\Zet test"

	// кнопки
	openButton := widget.NewButton("Открыть", func() {
		text := getText(selectedFile.filePath)
		data, _ := fileRead(selectedFile.filePath)
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
	searchSelect := widget.NewSelect([]string{topic, tag}, func(value string) { // todo: дата
		if value == topic {
			statusLabel.SetText("Введите слова, которые должна содержать тема")
		} else if value == tag {
			statusLabel.SetText("Введите один тег без знака #")
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
	searchBox := container.NewVBox(widget.NewLabel(""), container.NewBorder(nil, nil, searchSelect, clearButton, searchEntry), searchButtonBox)

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
		newFormatLabel(""),
		container.NewBorder(nil, nil, newFormatLabel("Имя:  "), nil, fileNameEntry),
		container.NewBorder(nil, nil, newFormatLabel("Тема: "), nil, topicEntry),
		container.NewBorder(nil, nil, newFormatLabel("Теги: "), nil, tagEntry),
		container.NewBorder(nil, nil, newFormatLabel("Дата: "), nil, dateEntry),
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
		selectedFile, err = fileRead(filePath)

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
