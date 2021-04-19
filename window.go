package fyneshell

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"strconv"
)

const APP_NAME = "Fyne Shell"
const APP_KEY = "com.github.tk103331.fyneshell"
const APP_DATA = "data"

type Window struct {
	app fyne.App
	win fyne.Window
	tabs *container.DocTabs
	terms []*TermTab
	confs []*Config
}

func (w *Window) AddTermTab(tab *TermTab) {
	w.terms = append(w.terms, tab)
	tabItem := container.TabItem{Text: tab.name, Icon: theme.ComputerIcon(), Content: tab.term}
	w.tabs.Append(&tabItem)
}

func (w *Window) AddConfig(conf *Config) {
	w.confs = append(w.confs, conf)
	w.save()
}

func (w *Window) RemoveConfig(index int) {
	if index < 0 || index > len(w.confs) {
		return
	}
	w.confs = append(w.confs[:index], w.confs[index+1:]...)
	w.save()
}

func (w *Window) Run() {

	w.app = app.NewWithID(APP_KEY)
	w.app.Settings().SetTheme(theme.DarkTheme())

	confs := w.app.Preferences().String(APP_DATA)
	err := json.Unmarshal([]byte(confs), &w.confs)
	if err != nil {
		log.Println(err)
	}

	w.win = w.app.NewWindow(APP_NAME)
	w.win.Resize(fyne.NewSize(800, 600))
	w.initUI()

	w.win.ShowAndRun()
}

func (w *Window) initUI() {
	toolbar := widget.NewToolbar(widget.NewToolbarAction(theme.ComputerIcon(), func() {
		tab,err := newLocalTermTab()
		if err != nil {
			dialog.NewError(err, w.win)
			return
		}
		w.AddTermTab(tab)
	}), widget.NewToolbarAction(theme.DocumentIcon(), func() {
		w.showNewSSHDialog()
	}),widget.NewToolbarSpacer(), widget.NewToolbarAction(theme.InfoIcon(), func() {
		w.showAboutDialog()
	}))

	quickbar := container.NewHBox(widget.NewButton("test", func() {

	}))

	sidebar := widget.NewList(func() int {
		return len(w.confs)
	}, func() fyne.CanvasObject {
		return container.NewHBox(widget.NewLabel(""), layout.NewSpacer(),
			widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), nil),
			widget.NewButtonWithIcon("", theme.DeleteIcon(), nil),
			widget.NewButtonWithIcon("", theme.ComputerIcon(), nil))
	}, func(id widget.ListItemID, object fyne.CanvasObject) {
		box := object.(*fyne.Container)
		label := box.Objects[0].(*widget.Label)
		edit := box.Objects[2].(*widget.Button)
		del := box.Objects[3].(*widget.Button)
		open := box.Objects[4].(*widget.Button)

		conf := w.confs[id]
		label.Text = conf.Name
		edit.OnTapped = func() {
			w.showEditSSHDialog(conf)
		}
		del.OnTapped = func() {
			w.RemoveConfig(id)
		}
		open.OnTapped = func() {
			tab,err := newSSHTermTab(conf)
			if err != nil {
				return
			}
			w.AddTermTab(tab)
		}
	})

	w.tabs = container.NewDocTabs(&container.TabItem{Text: "test", Icon: theme.ComputerIcon(), Content: widget.NewLabel("asdfsadfasdf")})
	//w.createLocalTermTab()
	w.tabs.OnClosed = func(item *container.TabItem) {

	}
	center := container.NewHSplit(sidebar, w.tabs)
	center.Offset = 0.3

	content := container.NewBorder(toolbar, quickbar, nil, nil, center)

	w.win.SetContent(content)
}

func (w *Window) showNewSSHDialog() {
	dlg := w.sshDialog(nil)
	dlg.Show()
}

func (w *Window) showEditSSHDialog(conf *Config) {
	dlg := w.sshDialog(conf)
	dlg.Show()
}

func (w *Window) sshDialog(conf *Config) dialog.Dialog {
	nameEntry := widget.NewEntry()
	hostEntry := widget.NewEntry()
	portEntry := widget.NewEntry()
	userEntry := widget.NewEntry()
	pswdEntry := widget.NewEntry()

	portEntry.Validator = func(s string) error {
		_, err := strconv.Atoi(s)
		return err
	}
	pswdEntry.Password = true

	title := "New SSH config"
	if conf != nil {
		nameEntry.Text = conf.Name
		nameEntry.Disable()
		hostEntry.Text = conf.Host
		portEntry.Text = strconv.Itoa(conf.Port)
		userEntry.Text = conf.User
		pswdEntry.Text = conf.Pswd

		title = "Modify SSH config"
	}


	dlg := dialog.NewForm(title, "OK", "Cancel", []*widget.FormItem{
		widget.NewFormItem("Name", nameEntry),
		widget.NewFormItem("Host", hostEntry),
		widget.NewFormItem("Port", portEntry),
		widget.NewFormItem("Username", userEntry),
		widget.NewFormItem("Password", pswdEntry),
	}, func(b bool) {
		if b {
			port,_ := strconv.Atoi(portEntry.Text)

			c := conf
			if c == nil {
				c = newConfig(nameEntry.Text)
			}
			c.Host = hostEntry.Text
			c.Port = port
			c.User = userEntry.Text
			c.Pswd = pswdEntry.Text

			if conf == nil {
				w.AddConfig(c)
			} else {
				w.save()
			}
		}
	}, w.win)

	dlg.Resize(fyne.NewSize(300, 400))
	return dlg
}

func (w *Window) showAboutDialog() {
	dialog.NewInformation(APP_NAME, "FyneShell is a shell client by fyne, for SSH or local PTY.", w.win).Show()
}

func (w *Window) createLocalTermTab() {
	tab,err := newLocalTermTab()
	if err != nil {
		dialog.NewError(err, w.win)
		return
	}
	w.AddTermTab(tab)
}

func (w *Window) save() {
	data, err := json.Marshal(w.confs)
	if err != nil {
		log.Println(err)
		return
	}
	w.app.Preferences().SetString(APP_DATA, string(data))
}

func Run() {
	(&Window{}).Run()
}
