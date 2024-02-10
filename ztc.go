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
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var gVersion, gYear, gProgramName string
var gFilePath string

func main() {
	gProgramName = "Zettelkasten"
	gYear = "2023"
	gVersion = "2" // сохр. в файл todo:

	a := app.New()
	w := a.NewWindow(gProgramName)
	w.Resize(fyne.NewSize(600, 450))
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
			fyne.NewMenuItem("Создать новую карточку", nil), // +F2 и  кнопку на редактирование todo:
			fyne.NewMenuItem("Редактировать файл", nil),
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

	go func() { // простите
		time.Sleep(2 * time.Second)
		for _, item := range menu.Items[0].Items {
			if item.Label == "Quit" {
				item.Label = "Выход"
			}
		}
	}()

	w.SetContent(mainForm())
	w.ShowAndRun()
}

func CreateNewCard() {
	// var newFile ztcBasicsType

	w := fyne.CurrentApp().NewWindow(selectedFile.id)
	w.Resize(fyne.NewSize(500, 400))
	w.CenterOnScreen()

	tabs := container.NewAppTabs(
		container.NewTabItem("создание", elmForm.getNameCard()),
		container.NewTabItem("просмотр", elmForm.viewForm()), // todo: не просмотр - текст форм?
		container.NewTabItem("доп.", elmForm.addInfoForm()),
		container.NewTabItem("источник", elmForm.sourceForm()),
		container.NewTabItem("коммент.", elmForm.commentForm()),
	)
	tabs.SetTabLocation(container.TabLocationBottom)

	w.SetContent(tabs)

	w.Show()
}

func ViewCard() {
	w := fyne.CurrentApp().NewWindow(selectedFile.id)
	w.Resize(fyne.NewSize(400, 400))
	w.CenterOnScreen()

	tabs := container.NewAppTabs(
		container.NewTabItem("просмотр", elmForm.viewForm()),
		container.NewTabItem("доп.", elmForm.addInfoForm()),
		container.NewTabItem("источник", elmForm.sourceForm()),
		container.NewTabItem("коммент.", elmForm.commentForm()),
	)
	tabs.SetTabLocation(container.TabLocationBottom)

	w.SetContent(tabs)
	w.Show()
}

func mainForm() (box *fyne.Container) {

	var list *widget.List
	gFilePath = "D:\\ztc test"

	files, err := os.ReadDir(gFilePath)
	if err != nil {
		// 	statusLabel.Text = "Ошибка: рабочая папка не открыта"
		fmt.Println(err)
	}
	/*	dirLabel := widget.NewLabel(gFilePath)

		dirButton := widget.NewButton("Каталог", func() {
			dialog := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
				if r != nil && err == nil {
					fmt.Println(r.URI())
					gFilePath = filepath.Dir(r.URI().Path())
					dirLabel.SetText(gFilePath)
					files, _ = os.ReadDir(gFilePath)
					list.Refresh()
				}
			},
				fyne.CurrentApp().Driver().AllWindows()[0],
			)

			// note: использовать если нужен выбор разных файлов + нужно считывать папку дополнительно
			// dialog.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
			var dir fyne.ListableURI
			d := storage.NewFileURI(gFilePath)
			dir, _ = storage.ListerForURI(d)

			dialog.SetLocation(dir)
			dialog.SetConfirmText("Выбрать")
			dialog.SetDismissText("Закрыть")
			dialog.Show()
		})
		dirBox := container.NewBorder(nil, nil, dirLabel, dirButton)
		topBottom := container.NewVBox(dirBox)*/

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

		selectedFile.filePath = filepath.Join(gFilePath, files[id].Name())
		selectedFile, err = fileRead(selectedFile.filePath)
		ViewCard()

		elmForm.refreshTabs(selectedFile)
	}

	btnCreate := widget.NewButton("Создать", func() {
		CreateNewCard()
	})
	btns := container.NewVBox(btnCreate)
	split := container.NewHSplit(list, btns)
	box = container.NewBorder(nil, nil, nil, nil, split)

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
