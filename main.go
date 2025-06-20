package main

import (
	"image/color"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Selections struct {
	gate string // Selected gate name
	pos  string // Selected position (relative to gate)
	year int    // Selected year
}

func main() {
	config := GetConfig()
	// Initialize the database connection using config values
	InitializeDatabase(
		config.Credentials.User,
		config.Credentials.Pass,
		config.Credentials.Host,
		config.Credentials.Port,
		config.GateConfig.Dbname,
	)
	var gs *GateSpectator = nil // Used for real-time gate updates
	// REQUIRED FOR FLAMEGRAPH
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Ensure cleanup functions are called when the application exits.
	defer func() {
		if gs != nil {
			gs.Shutdown() // Clean up GateSpectator goroutines
		}
		DB.Close() // Always close DB connection
	}()
	defer DB.Close() // Redundant, but ensures DB is closed

	selections := &Selections{} // Holds current user selections
	myApp := app.New()
	myWindow := myApp.NewWindow("Gate Keeper")

	// Create labels for selection controls
	yearLabel := widget.NewLabel("Select a year")
	gateLabel := widget.NewLabel("Select a gate")
	posLabel := widget.NewLabel("Position relative to gate")

	// Declare select fields for year, gate, and position
	var yearOptionsSelect *widget.Select
	var gateOptionsSelect *widget.Select
	var posOptionsSelect *widget.Select

	// VBox to hold selection controls and button
	controlsVBox := container.NewVBox()

	// Placeholder for the dynamic table
	tablePlaceholderContainer := container.NewMax() // Ensures table fills available space

	// Main layout: controls at top, table in center
	mainLayoutContainer := container.NewBorder(
		controlsVBox,              // Top: selection controls
		nil,                       // Bottom: none
		nil,                       // Left: none
		nil,                       // Right: none
		tablePlaceholderContainer, // Center: table
	)

	// Year select: populates gates and table when changed
	yearOptionsSelect = widget.NewSelect([]string{}, func(value string) {
		year, _ := strconv.Atoi(value)
		selections.year = year
		gates := selectAllGates(&config.GateConfig, selections.year)
		_isTimelineLinear := isTimelineLinear(gates)
		var gateOptions []string
		for _, g := range gates {
			gateOptions = append(gateOptions, g.GateName)
		}
		selections.gate = ""
		gateOptionsSelect.ClearSelected()
		gateOptionsSelect.SetOptions(gateOptions)

		// Create a new table for the selected year/gates
		newTable := widget.NewTable(
			func() (int, int) {
				return len(gates), 2 // Rows, Columns
			},
			func() fyne.CanvasObject {
				// Template for each cell: background + label
				bg := canvas.NewRectangle(color.White)
				label := widget.NewLabel("Cell Data")
				return container.NewStack(bg, label)
			},
			func(id widget.TableCellID, o fyne.CanvasObject) {
				// Update cell content and color
				stack := o.(*fyne.Container)
				bg := stack.Objects[0].(*canvas.Rectangle)
				label := stack.Objects[1].(*widget.Label)

				if id.Col == 0 {
					label.SetText(gates[id.Row].GateName)
				}
				if id.Col == 1 {
					p := RelativePosition(getGatePosition(gates[id.Row]))
					positions := []string{"Past", "Open", "Future"}
					label.SetText(positions[p+1])
				}

				// Color code: yellow if not linear, green if open, red otherwise
				if !_isTimelineLinear {
					bg.FillColor = color.RGBA{R: 255, G: 255, B: 0, A: 128}
				} else if isGateOpen(gates[id.Row]) {
					bg.FillColor = color.RGBA{R: 0, G: 255, B: 0, A: 128}
				} else {
					bg.FillColor = color.RGBA{R: 255, G: 0, B: 0, A: 128}
				}
				bg.Refresh()
			},
		)
		newTable.SetColumnWidth(0, 300)
		newTable.SetColumnWidth(1, 125)

		// Replace old table with new one
		myWindow.Resize(fyne.NewSize(500, 600))
		tablePlaceholderContainer.RemoveAll()
		tablePlaceholderContainer.Add(newTable)
		tablePlaceholderContainer.Refresh()

		// --- Real-time updates with GateSpectator ---
		// Shutdown any existing GateSpectator to avoid goroutine leaks.
		if gs != nil {
			gs.Shutdown()
		}
		gs = NewGateSpectator(&config.GateConfig, year)

		// Listen for gate updates and refresh table
		go func() {
			for _gates := range gs.gatesUpdate {
				_isTimelineLinear = isTimelineLinear(_gates)
				// Update gate times in the current gates slice
				for _, g := range _gates {
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

	// Gate select: updates selection
	gateOptionsSelect = widget.NewSelect([]string{}, func(value string) {
		selections.gate = value
	})

	// Position select: updates selection
	posOptionsSelect = widget.NewSelect([]string{}, func(value string) {
		selections.pos = value
	})

	// Populate year and position options
	yearOptionsSelect.SetOptions(selectAllYears(&config.GateConfig))
	posOptionsSelect.SetOptions(getPositionOptions())

	// Button to set gates relative to selection
	button := widget.NewButton("Set Gates", func() {
		if selections.gate == "" || selections.pos == "" || selections.year == 0 {
			// Show popup if any selection is missing
			popupLabel := widget.NewLabel("All selections are required")
			popup := widget.NewModalPopUp(container.NewVBox(popupLabel), fyne.CurrentApp().Driver().CanvasForObject(posOptionsSelect))
			popup.Show()
			<-time.NewTimer(3 * time.Second).C
			popup.Hide()
			return
		}
		// Update gates in DB based on selection
		setGatesRelativeTo(&config.GateConfig, selections.year, selections.gate, RelativePositionStr(selections.pos).Value())
	})

	// Add controls to the VBox
	controlsVBox.Add(yearLabel)
	controlsVBox.Add(yearOptionsSelect)
	controlsVBox.Add(gateLabel)
	controlsVBox.Add(gateOptionsSelect)
	controlsVBox.Add(posLabel)
	controlsVBox.Add(posOptionsSelect)
	controlsVBox.Add(button)

	myWindow.Resize(fyne.NewSize(500, 320))
	myWindow.SetContent(mainLayoutContainer)
	myWindow.ShowAndRun()
}
