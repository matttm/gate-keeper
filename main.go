package main

import (
	"log"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Selections struct {
	gate string
	pos  string
	year int
}

func main() {
	config := GetEnvironment()
	InitializeDatabase(
		config.Credentials.User,
		config.Credentials.Pass,
		config.Credentials.Host,
		config.Credentials.Port,
		config.GateConfig.Dbname,
	)
	selections := &Selections{}
	myApp := app.New()
	myWindow := myApp.NewWindow("Gate Keeper")

	// create labels
	yearLabel := widget.NewLabel("Select a year")
	gateLabel := widget.NewLabel("Select a gate")
	posLabel := widget.NewLabel("Position relative to gate")
	// create select fields
	var yearOptionsSelect *widget.Select
	var gateOptionsSelect *widget.Select
	var posOptionsSelect *widget.Select
	yearOptionsSelect = widget.NewSelect([]string{}, func(value string) {
		year, _ := strconv.Atoi(value)
		selections.year = year
		log.Println("Select set to", value)
		gates := selectAllGates(&config.GateConfig, selections.year)
		var gateOptions []string
		for _, g := range gates {
			gateOptions = append(gateOptions, g.GateName)
		}
		selections.gate = ""
		gateOptionsSelect.SetSelected("")
		gateOptionsSelect.SetOptions(gateOptions)
	})
	gateOptionsSelect = widget.NewSelect([]string{}, func(value string) {
		selections.gate = value
		log.Println("Select set to", value)
	})
	posOptionsSelect = widget.NewSelect([]string{}, func(value string) {
		selections.pos = value
		log.Println("Select set to", value)
	})

	yearOptionsSelect.SetOptions(selectAllYears(&config.GateConfig))
	posOptionsSelect.SetOptions(getPositionOptions())
	button := widget.NewButton("Set Gates", func() {
		if selections.gate == "" || selections.pos == "" || selections.year == 0 {
			popupLabel := widget.NewLabel("All selections are required")
			popup := widget.NewModalPopUp(container.NewVBox(popupLabel), fyne.CurrentApp().Driver().CanvasForObject(posOptionsSelect))
			popup.Show()
			<-time.NewTimer(3 * time.Second).C
			popup.Hide()
			return
		}
		log.Println("tapped")
		setGatesRelativeTo(&config.GateConfig, selections.year, selections.gate, RelativePositionStr(selections.pos).Value())
	})
	myWindow.Resize(fyne.NewSize(500, 300))

	myWindow.SetContent(container.NewVBox(yearLabel, yearOptionsSelect, gateLabel, gateOptionsSelect, posLabel, posOptionsSelect, button))
	myWindow.ShowAndRun()
}
