package main

import (
	"image/color"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Selections struct {
	gate string
	pos  string
	year int
}

func main() {
	config := GetConfig()
	InitializeDatabase(
		config.Credentials.User,
		config.Credentials.Pass,
		config.Credentials.Host,
		config.Credentials.Port,
		config.GateConfig.Dbname,
	)
	var gs *GateSpectator = nil
	defer gs.Shutdown()
	defer DB.Close()
	selections := &Selections{}
	myApp := app.New()
	myWindow := myApp.NewWindow("Gate Keeper")
	panelContainer := container.NewVBox()

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
		gates := selectAllGates(&config.GateConfig, selections.year)
		var gateOptions []string
		for _, g := range gates {
			gateOptions = append(gateOptions, g.GateName)
		}
		selections.gate = ""
		gateOptionsSelect.ClearSelected()
		gateOptionsSelect.SetOptions(gateOptions)

		// the rest of this cb is fr the table
		gs = NewGateSpectator(&config.GateConfig, year)

		table := widget.NewTable(
			func() (int, int) {
				return len(gates), 1 // Rows, Columns
			},
			func() fyne.CanvasObject {
				// This creates the template object for each cell
				bg := canvas.NewRectangle(color.White) // Default background for the cell
				label := widget.NewLabel("Cell Data")  // Content of the cell
				return container.NewStack(bg, label)   // Stack to place background behind label
			},
			func(id widget.TableCellID, o fyne.CanvasObject) {
				// This function is called to update the content of each cell
				stack := o.(*fyne.Container)
				bg := stack.Objects[0].(*canvas.Rectangle)
				label := stack.Objects[1].(*widget.Label)

				label.SetText(gates[id.Row].GateName)

				// Apply color based on the index
				if isGateOpen(gates[id.Row]) {
					bg.FillColor = color.RGBA{R: 0, G: 255, B: 0, A: 128}
				} else {
					bg.FillColor = color.RGBA{R: 255, G: 0, B: 0, A: 128} // Red with some transparency
				}
				bg.Refresh() // Important to refresh the rectangle after changing its color
			},
		)
		go func() {
			for _gates := range gs.gatesUpdate {
				// TODO: do this differently
				for _, g := range _gates {
					// find current gate in gates and update dates
					for _, _g := range gates {
						if g.GateName == _g.GateName {
							_g.Start = g.Start
							_g.End = g.End
						}
					}
				}
				table.Refresh()
			}
		}()
		panelContainer.Add(table)
	})
	gateOptionsSelect = widget.NewSelect([]string{}, func(value string) {
		selections.gate = value
	})
	posOptionsSelect = widget.NewSelect([]string{}, func(value string) {
		selections.pos = value
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
		setGatesRelativeTo(&config.GateConfig, selections.year, selections.gate, RelativePositionStr(selections.pos).Value())
	})
	myWindow.Resize(fyne.NewSize(500, 300))

	panelContainer.Add(container.NewVBox(yearLabel, yearOptionsSelect, gateLabel, gateOptionsSelect, posLabel, posOptionsSelect, button))
	myWindow.SetContent(panelContainer)
	myWindow.ShowAndRun()
}
