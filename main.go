package main

import (
	"log"
	"strconv"

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
	environment := GetEnvironment()
	InitializeDatabase(
		environment.user,
		environment.pass,
		environment.host,
		environment.port,
		environment.config.Dbname,
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
		gates := selectAllGates(environment.config, selections.year)
		var gateOptions []string
		for _, g := range gates {
			gateOptions = append(gateOptions, g.GateName)
		}
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

	yearOptionsSelect.SetOptions(selectAllYears(environment.config))
	posOptionsSelect.SetOptions(getPositionOptions())
	button := widget.NewButton("Set Gates", func() {
		if selections.gate == "" || selections.pos == "" || selections.year == 0 {
			log.Fatal("All selections are required")
		}
		log.Println("tapped")
		setGatesRelativeTo(environment.config, selections.year, selections.gate, RelativePositionStr(selections.pos).Value())
	})
	myWindow.Resize(fyne.NewSize(500, 300))

	myWindow.SetContent(container.NewVBox(yearLabel, yearOptionsSelect, gateLabel, gateOptionsSelect, posLabel, posOptionsSelect, button))
	myWindow.ShowAndRun()
}
