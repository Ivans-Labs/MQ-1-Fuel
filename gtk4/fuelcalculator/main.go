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
	gtk.StyleContextAddProviderForDisplay(gdk.DisplayGetDefault(), cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

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

	// Create radio buttons for Day VFR and Night VFR
	dayVFRButton := gtk.Button(dayVFRButton)
	dayVFRButtonLabel := gtk.NewLabel("Day VFR (30 mins)")
	dayVFRButton.SetChild(dayVFRButtonLabel) // Manually set the label as a child of the radio button

	nightVFRButton := gtk.Button(nighVFR)
	nightVFRButtonLabel := gtk.NewLabel("Night VFR (45 mins)")
	nightVFRButton.SetChild(nightVFRButtonLabel) // Manually set the label as a child of the radio button

	// Create a button to calculate the fuel metrics
	calculateButton := gtk.NewButtonWithLabel("Calculate Fuel Metrics")

	// Create a label to display the results
	resultLabel := gtk.NewLabel("")

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
	grid.Attach(dayVFRButton, 0, 6, 2, 1)
	grid.Attach(nightVFRButton, 0, 7, 2, 1)
	grid.Attach(calculateButton, 0, 8, 2, 1)
	grid.Attach(resultLabel, 0, 9, 2, 1)

	// Set up the calculation logic
	calculateButton.ConnectClicked(func() {
		forwardText := forwardEntry.Text()
		aftText := aftEntry.Text()
		brText := brEntry.Text()
		descentFuelText := descentFuelEntry.Text()
		groundspeedText := groundspeedEntry.Text()
		distanceFromBaseText := distanceFromBaseEntry.Text()

		// Parse input values
		forward, _ := strconv.ParseFloat(forwardText, 64)
		aft, _ := strconv.ParseFloat(aftText, 64)
		br, _ := strconv.ParseFloat(brText, 64)
		descentFuel, _ := strconv.ParseFloat(descentFuelText, 64)
		groundspeed, _ := strconv.ParseFloat(groundspeedText, 64)
		distanceFromBase, _ := strconv.ParseFloat(distanceFromBaseText, 64)
		dayVFR := dayVFRButton.GetActive()

		// Calculate Bingo fuel and related metrics
		bingoFuel, onstaTime, rtbTime := calculateFuelMetrics(forward, aft, br, descentFuel, groundspeed, distanceFromBase, dayVFR)

		// Display the results
		resultLabel.SetText(fmt.Sprintf("Bingo Fuel: %.2f lbs\nONSTA Time: %.2f minutes\nRTB Time: %.2f minutes", bingoFuel, onstaTime, rtbTime))
	})

	// Set the grid as the child of the window
	window.SetChild(grid)
	window.Show()
}

func calculateFuelMetrics(forward, aft, br, descentFuel, groundspeed, distanceFromBase float64, dayVFR bool) (bingoFuel, onstaTime, rtbTime float64) {
	totalFuel := forward + aft
	vfrFuel := br / 60 * 30 // Day VFR reserve (30 mins)
	if !dayVFR {
		vfrFuel = br / 60 * 45 // Night VFR reserve (45 mins)
	}

	// Calculate Bingo Fuel
	bingoFuel = taxiAndTakeoffFuel + vfrFuel + descentFuel + (distanceFromBase / groundspeed * 60 * br / 60)

	// Calculate ONSTA Time (Time on station)
	onstaTime = (totalFuel - bingoFuel) / br

	// Calculate RTB (Return to Base) Time
	rtbTime = distanceFromBase / groundspeed

	return bingoFuel, onstaTime, rtbTime
}
