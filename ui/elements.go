package ui

//go:generate fyne bundle --package ui -o bundled.go ../res/gvaw-sq-86.png
//fyne bundle --package ui -o bundled.go -append image2.png

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"log"
)

var (
	syncButton    *widget.Button
	authCode      *widget.Entry
	syncServer    *widget.Entry
	savedGamesDir *widget.Entry
	persSkin      *widget.Check
	skinSize      *widget.Select
	persSize      *widget.Select
	statusMajor   *widget.ProgressBar
	statusMinor   *widget.ProgressBar
	window        fyne.Window
	pStatus       programStatusCode
	fyneApp       fyne.App
)

var addressPort string

func init() {
	pStatus = statusIdle
}

func initElements() *container.AppTabs {

	elements := container.NewAppTabs(
		container.NewTabItemWithIcon("SYNC", theme.ViewRefreshIcon(), initSyncPage()),
		container.NewTabItemWithIcon("HELP", theme.HelpIcon(), initHelpPage()),
	)

	elements.SetTabLocation(container.TabLocationTop)
	return elements
}

func initSyncPage() *fyne.Container {

	gvawlogo := canvas.NewImageFromResource(resourceGvawSq86Png)
	gvawlogo.FillMode = canvas.ImageFillOriginal
	gvawlogo.Resize(fyne.NewSize(86, 86))

	authCode = widget.NewEntry()
	authCode.Validator = nil
	authCode.Resize(fyne.NewSize(90, 0))

	syncServer = widget.NewEntry()
	syncServer.Validator = nil
	syncServer.Resize(fyne.NewSize(90, 0))
	syncServer.Text = "gvaw.web.id"

	savedGamesDir = widget.NewEntry()
	savedGamesDir.Validator = nil
	savedGamesDir.Resize(fyne.NewSize(90, 0))

	persSkin = widget.NewCheck("Enable", func(value bool) {
		log.Println("persSkin set to", value)
	})
	persSkin.SetChecked(true)
	skinSize = widget.NewSelect([]string{"Full", "Half", "Quarter"}, func(value string) {
		log.Println("skinSize set to", value)
	})
	skinSize.SetSelectedIndex(0)
	persSize = widget.NewSelect([]string{"Full", "Half", "Quarter"}, func(value string) {
		log.Println("persSize set to", value)
	})
	persSize.SetSelectedIndex(0)

	syncButton = widget.NewButtonWithIcon("SYNCHRONIZE", theme.ViewRefreshIcon(), sync)

	statusMajor = widget.NewProgressBar()
	statusMajor.TextFormatter = func() string {
		switch pStatus {
		case statusConnected:
			return "Connected"
		case statusSyncing:
			return "Syncing"
		default:
			return "Idle"
		}
	}

	statusMinor = widget.NewProgressBar()
	statusMinor.TextFormatter = func() string {
		switch pStatus {
		case statusSyncing:
			return "Syncing"
		default:
			return ""
		}
	}

	return container.New(layout.NewVBoxLayout(),
		layout.NewSpacer(),
		container.New(layout.NewHBoxLayout(),
			layout.NewSpacer(), gvawlogo, layout.NewSpacer(),
		),
		container.New(layout.NewHBoxLayout(),
			layout.NewSpacer(), canvas.NewText("GVAW SkinSync", color.White), layout.NewSpacer(),
		),
		layout.NewSpacer(),
		container.New(layout.NewHBoxLayout(),
			layout.NewSpacer(),
			container.New(layout.NewVBoxLayout(),
				container.New(layout.NewFormLayout(),
					widget.NewLabel("Auth Code"), authCode,
					widget.NewLabel("Skin Resolution"), skinSize,
					widget.NewLabel("Personal Skins"), persSkin,
					widget.NewLabel("Personal Skin Res"), persSize,
					widget.NewLabel("Saved Games Dir"), savedGamesDir,
					widget.NewLabel("Sync Server"), syncServer,
				),
				syncButton,
			),
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
		container.New(layout.NewVBoxLayout(),
			container.New(layout.NewHBoxLayout(),
				layout.NewSpacer(),
				canvas.NewText("Status:", color.White),
				layout.NewSpacer(),
			),
			statusMajor, statusMinor,
		),
		layout.NewSpacer(),
	)
}

func initHelpPage() *fyne.Container {
	about := widget.NewLabel("GVAW SkinSync v1.00")
	t0 := widget.NewLabel("This open source software is created by Dalva.\nLicensed under GNU Affero-GPLv3.0.\nSource at https://github.com/dalva24/GVAWSkinSync")
	t01 := widget.NewLabel("Powered by modified Alastor massively concurrent\nfile transfer algorithm to ensure peak maximum\nspeed even in Ind*hom* conditions.\nReference source at https://github.com/dalva24/alastor")
	t1 := widget.NewLabel("Auth Code is like IFF Code, renewed daily.\nAsk around to obtain it.")
	t2 := widget.NewLabel("Choose resolution depending on\nyour available space and compute power")
	t3 := widget.NewLabel("Personal skin is personal helmets and stuff.\nYou can either enable or disable it here")
	t4 := widget.NewLabel("Saved Games Dir is automatically detected.\nSimply overwrite if it's wrong.")
	t5 := widget.NewLabel("Sync server should be server.gvaw.web.id\nIt can be edited for future-proofing.")
	return container.New(layout.NewVBoxLayout(), layout.NewSpacer(), about, t0, t01, t1, t2, t3, t4, t5, layout.NewSpacer())
}
