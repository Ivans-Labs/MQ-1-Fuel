package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

const (
	taxiAndTakeoffFuel = 35.0 // Taxi and Takeoff Fuel in lbs
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

	// Create input fields for forward, aft tank values, burn rate, descent fuel, groundspeed, and distance from base
	forwardEntry := gtk.NewEntry()
	aftEntry := gtk.NewEntry()
	brEntry := gtk.NewEntry()
	descentFuelEntry := gtk.NewEntry()
	groundspeedEntry := gtk.NewEntry()
	distanceFromBaseEntry := gtk.NewEntry()

	// Create labels for the input fields
	forwardLabel := gtk.NewLabel("Forward Tank (lbs):")
	aftLabel := gtk.NewLabel("Aft Tank (lbs):")
	brLabel := gtk.NewLabel("Burn Rate (lbs/hr):")
	descentFuelLabel := gtk.NewLabel("Descent Fuel (lbs):")
	groundspeedLabel := gtk.NewLabel("Groundspeed (knots):")
	distanceFromBaseLabel := gtk.NewLabel("Distance from Base (NM):")

	// Create a label to display the results
	resultLabel := gtk.NewLabel("")

	// Create a button to calculate the fuel metrics
	calculateButton := gtk.NewButtonWithLabel("Calculate Fuel Metrics")
	calculateButton.ConnectClicked(func() {
		// Retrieve and validate the user's input
		forward, err := strconv.ParseFloat(forwardEntry.Text(), 64)
		if handleError(err, resultLabel) {
			return
		}
		aft, err := strconv.ParseFloat(aftEntry.Text(), 64)
		if handleError(err, resultLabel) {
			return
		}
		br, err := strconv.ParseFloat(brEntry.Text(), 64) // Burn rate in lbs/hr
		if handleError(err, resultLabel) {
			return
		}
		br /= 60 // Convert burn rate to lbs/min for calculations
		descentFuel, err := strconv.ParseFloat(descentFuelEntry.Text(), 64)
		if handleError(err, resultLabel) {
			return
		}
		groundspeed, err := strconv.ParseFloat(groundspeedEntry.Text(), 64)
		if handleError(err, resultLabel) {
			return
		}
		distanceFromBase, err := strconv.ParseFloat(distanceFromBaseEntry.Text(), 64)
		if handleError(err, resultLabel) {
			return
		}

		// Calculate the required fuel values
		vfrReserve := 45.0 // VFR reserve in minutes (use 30 for day VFR)
		transitTime := distanceFromBase / groundspeed
		rtbFuel := transitTime * 60 * br // RTB fuel in lbs
		bingoFuel := taxiAndTakeoffFuel + (br * vfrReserve) + descentFuel + rtbFuel
		onStationTime := (forward + aft - bingoFuel) / br // ONSTA time in minutes

		// Convert ONSTA time to hours and minutes
		onstaHours := int(onStationTime) / 60
		onstaMinutes := int(onStationTime) % 60

		// Display the results
		resultLabel.SetText(fmt.Sprintf("Bingo Fuel: %.2f lbs\nONSTA Time: %02d hr %02d min\nRTB Fuel: %.2f lbs\nTransit Time: %.2f hr", bingoFuel, onstaHours, onstaMinutes, rtbFuel, transitTime))
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

// handleError is a helper function to handle input errors and update the result label
func handleError(err error, resultLabel *gtk.Label) bool {
	if err != nil {
		resultLabel.SetText("Error: " + err.Error())
		return true
	}
	return false
}

// No additional code needed after this point for the fuel calculator application
