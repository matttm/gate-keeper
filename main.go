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
		environment.config.dbname,
	)
	myApp := app.New()
	myWindow := myApp.NewWindow("Gate Keeper")

	gateLabel := widget.NewLabel("Select a gate")
	posLabel := widget.NewLabel("Position relativr to gate")
	c1 := widget.NewSelect([]string{"Option 1", "Option 2"}, func(value string) {
		log.Println("Select set to", value)
	})
	c2 := widget.NewSelect([]string{"Option 1", "Option 2"}, func(value string) {
		log.Println("Select set to", value)
	})
	button := widget.NewButton("click me", func() {
		log.Println("tapped")
	})
	myWindow.Resize(fyne.NewSize(500, 300))

	myWindow.SetContent(container.NewVBox(gateLabel, c1, posLabel, c2, button))
	myWindow.ShowAndRun()
}
