package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	environment := GetEnvironment()
	InitializeDatabase(
		environment.user,
		environment.pass,
		environment.host,
		environment.port,
		environment.config.Dbname,
	)
	myApp := app.New()
	myWindow := myApp.NewWindow("Gate Keeper")
	gates := selectAllGates(environment.config, 2026)
	var gateOptions []string
	for _, g := range gates {
		gateOptions = append(gateOptions, g.GateName)
	}

	gateLabel := widget.NewLabel("Select a gate")
	posLabel := widget.NewLabel("Position relative to gate")
	gateOptionsSelect := widget.NewSelect(gateOptions, func(value string) {
		log.Println("Select set to", value)
	})
	posOptionsSelect := widget.NewSelect(posOptions, func(value string) {
		log.Println("Select set to", value)
	})
	button := widget.NewButton("Set Gates", func() {
		log.Println("tapped")
	})
	myWindow.Resize(fyne.NewSize(500, 300))

	myWindow.SetContent(container.NewVBox(gateLabel, gateOptionsSelect, posLabel, posOptionsSelect, button))
	myWindow.ShowAndRun()
}
