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
	gateOptions := selectAllGates(environment.config, 2026)

	gateLabel := widget.NewLabel("Select a gate")
	posLabel := widget.NewLabel("Position relative to gate")
	gateOptionsSelect := widget.NewSelect(gateOptions, func(value string) {
		log.Println("Select set to", value)
	})
	posOptions := []string{"Outside before", "Inside", "Outside after"}
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
