package main

import (
	"fmt"
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
	year string
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
	gates := selectAllGates(environment.config, 2026)
	var gateOptions []string
	for _, g := range gates {
		gateOptions = append(gateOptions, g.GateName)
	}

	yearLabel := widget.NewLabel("Select a year")
	gateLabel := widget.NewLabel("Select a gate")
	posLabel := widget.NewLabel("Position relative to gate")
	yearOptionsSelect := widget.NewSelect([]string{"2026"}, func(value string) {
		selections.year = value
		log.Println("Select set to", value)
	})
	gateOptionsSelect := widget.NewSelect(gateOptions, func(value string) {
		selections.gate = value
		log.Println("Select set to", value)
	})
	posOptionsSelect := widget.NewSelect(getPositionOptions(), func(value string) {
		selections.pos = value
		log.Println("Select set to", value)
	})
	button := widget.NewButton("Set Gates", func() {
		if selections.gate == "" || selections.pos == "" || selections.year == "" {
			log.Fatal("All selections are required")
		}
		log.Println("tapped")
		year, _ := strconv.Atoi(selections.year)
		fmt.Println(year, int(year))
		setGatesRelativeTo(environment.config, int(year), selections.gate, RelativePositionStr(selections.pos).Value())
	})
	myWindow.Resize(fyne.NewSize(500, 300))

	myWindow.SetContent(container.NewVBox(yearLabel, yearOptionsSelect, gateLabel, gateOptionsSelect, posLabel, posOptionsSelect, button))
	myWindow.ShowAndRun()
}
