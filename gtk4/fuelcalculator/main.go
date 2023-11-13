package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

const (
	taxiAndTakeoffFuel = 35.0 // TM (Taxi and Takeoff Fuel) in lbs
)

func main() {
	app := gtk.NewApplication("com.example.fuelcalculator", gio.ApplicationFlagsNone)
	app.ConnectActivate(func() { activate(app) })
	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *gtk.Application) {
	window := gtk.NewApplicationWindow(app)
	window.SetTitle("Fuel Calculator")
	window.SetDefaultSize(400, 300)

	// Apply dark theme and white text using CSS
	cssProvider := gtk.NewCSSProvider()
	cssProvider.LoadFromData(`
		window {
			background-color: #2e2e2e;
			color: white;
		}
		entry {
			color: black;
		}
		button {
			color: black;
		}
		label {
			color: white;
		}
	`)
	display := gdk.DisplayGetDefault()
	gtk.StyleContextAddProviderForDisplay(display, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	// Create a grid to organize widgets
	grid := gtk.NewGrid()
	grid.SetRowSpacing(10)
	grid.SetColumnSpacing(10)
	grid.SetMarginStart(10)
	grid.SetMarginEnd(10)
	grid.SetMarginTop(10)
	grid.SetMarginBottom(10)

	// Create input fields for forward, aft tank values, and burn rate
	forwardEntry := gtk.NewEntry()
	aftEntry := gtk.NewEntry()
	brEntry := gtk.NewEntry()
	descentFuelEntry := gtk.NewEntry()
	groundspeedEntry := gtk.NewEntry()
	distanceFromBaseEntry := gtk.NewEntry()

	// Create labels for the input fields
	forwardLabel := gtk.NewLabel("Forward Tank (lbs):")
	aftLabel := gtk.NewLabel("Aft Tank (lbs):")
	brLabel := gtk.NewLabel("Burn Rate (lbs/min):")
	descentFuelLabel := gtk.NewLabel("Descent Fuel (lbs):")
	groundspeedLabel := gtk.NewLabel("Groundspeed (knots):")
	distanceFromBaseLabel := gtk.NewLabel("Distance from Base (nm):")

	// Create a label to display the results
	resultLabel := gtk.NewLabel("")

	// Create a button to calculate the fuel metrics
	calculateButton := gtk.NewButtonWithLabel("Calculate Fuel Metrics")
	calculateButton.ConnectClicked(func() {
		// Retrieve and validate the user's input
		forward, err := strconv.ParseFloat(forwardEntry.Text(), 64)
		if err != nil {
			log.Fatal("Invalid forward tank value")
			return
		}
		aft, err := strconv.ParseFloat(aftEntry.Text(), 64)
		if err != nil {
			log.Fatal("Invalid aft tank value")
			return
		}
		br, err := strconv.ParseFloat(brEntry.Text(), 64)
		if err != nil {
			log.Fatal("Invalid burn rate value")
			return
		}
		descentFuel, err := strconv.ParseFloat(descentFuelEntry.Text(), 64)
		if err != nil {
			log.Fatal("Invalid descent fuel value")
			return
		}
		groundspeed, err := strconv.ParseFloat(groundspeedEntry.Text(), 64)
		if err != nil {
			log.Fatal("Invalid groundspeed value")
			return
		}
		distanceFromBase, err := strconv.ParseFloat(distanceFromBaseEntry.Text(), 64)
		if err != nil {
			log.Fatal("Invalid distance from base value")
			return
		}

		// Calculate the required fuel values
		vfrReserve := 45.0                               // VFR reserve in minutes (use 30 for day VFR)
		rtbFuel := (distanceFromBase / groundspeed) * 60 // RTB fuel in lbs
		bingoFuel := taxiAndTakeoffFuel + (br * vfrReserve / 60) + descentFuel + rtbFuel
		onStationTime := (forward + aft - bingoFuel) / br // ONSTA time in minutes

		// Display the results
		resultLabel.SetText(fmt.Sprintf("Bingo Fuel: %.2f lbs\nONSTA Time: %.2f mins\nRTB Fuel: %.2f lbs", bingoFuel, onStationTime, rtbFuel))
	})

	// Add widgets to the grid
	grid.Attach(forwardLabel, 0, 0, 1, 1)
	grid.Attach(forwardEntry, 1, 0, 1, 1)
	grid.Attach(aftLabel, 0, 1, 1, 1)
	grid.Attach(aftEntry, 1, 1, 1, 1)
	grid.Attach(brLabel, 0, 2, 1, 1)
	grid.Attach(brEntry, 1, 2, 1, 1)
	grid.Attach(descentFuelLabel, 0, 3, 1, 1)
	grid.Attach(descentFuelEntry, 1, 3, 1, 1)
	grid.Attach(groundspeedLabel, 0, 4, 1, 1)
	grid.Attach(groundspeedEntry, 1, 4, 1, 1)
	grid.Attach(distanceFromBaseLabel, 0, 5, 1, 1)
	grid.Attach(distanceFromBaseEntry, 1, 5, 1, 1)
	grid.Attach(calculateButton, 0, 6, 2, 1)
	grid.Attach(resultLabel, 0, 7, 2, 1)

	// Set the grid as the child of the window
	window.SetChild(grid)

	// Show the window
	window.Show()
}

// No additional code needed after this point for the fuel calculator application
