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

	// Ensure cleanup functions are called when the application exits.
	defer func() {
		if gs != nil {
			gs.Shutdown()
		}
		// Closing DB here, ensures DB is closed even if GS is nil
		DB.Close()
	}()
	defer DB.Close()
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

	// This VBox will hold all the selection controls and the button at the top of the window.
	// We'll add the individual widgets to it after they are initialized.
	controlsVBox := container.NewVBox()

	// This container will act as the placeholder for the table.
	// It's placed in the center of the main Border layout.
	tablePlaceholderContainer := container.NewMax() // Use NewMax to ensure the table fills this space

	// This is the **main container** for the window. It uses a Border layout.
	// The `controlsVBox` will be at the top, and the `tablePlaceholderContainer` will be in the center,
	// allowing the table to expand and fill the rest of the window.
	mainLayoutContainer := container.NewBorder(
		controlsVBox,              // Top content: our selection controls
		nil,                       // Bottom content (none for now)
		nil,                       // Left content (none for now)
		nil,                       // Right content (none for now)
		tablePlaceholderContainer, // Center content: this will hold the table dynamically
	)
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

		newTable := widget.NewTable(
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
		newTable.SetColumnWidth(0, 300)

		// Remove old table and add the new one to the placeholder container.
		tablePlaceholderContainer.RemoveAll()   // Clear any existing content
		tablePlaceholderContainer.Add(newTable) // Add the newly created table
		tablePlaceholderContainer.Refresh()     // Refresh the placeholder container to show the new table

		// --- Real-time updates with GateSpectator ---
		// Shutdown any existing GateSpectator to avoid goroutine leaks.
		if gs != nil {
			gs.Shutdown()
		}
		gs = NewGateSpectator(&config.GateConfig, year)

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
				newTable.Refresh()
			}
		}()
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
	controlsVBox.Add(yearLabel)
	controlsVBox.Add(yearOptionsSelect)
	controlsVBox.Add(gateLabel)
	controlsVBox.Add(gateOptionsSelect)
	controlsVBox.Add(posLabel)
	controlsVBox.Add(posOptionsSelect)
	controlsVBox.Add(button)
	myWindow.Resize(fyne.NewSize(500, 600))

	// panelContainer.Add(container.NewVBox(yearLabel, yearOptionsSelect, gateLabel, gateOptionsSelect, posLabel, posOptionsSelect, button))
	myWindow.SetContent(mainLayoutContainer)
	myWindow.ShowAndRun()
}
