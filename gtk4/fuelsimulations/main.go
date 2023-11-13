package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
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

	// Create labels for the input fields
	forwardLabel := gtk.NewLabel("Forward Tank (lbs):")
	aftLabel := gtk.NewLabel("Aft Tank (lbs):")
	brLabel := gtk.NewLabel("Burn Rate (lbs/min):")

	// Create a button to start the simulation
	startButton := gtk.NewButtonWithLabel("Start Simulation")

	// Create a label to display the results
	resultLabel := gtk.NewLabel("")
	resultLabel.SetCSSClasses([]string{"result-label"}) // Add CSS class for styling

	// Create a button to simulate wind effect
	windButton := gtk.NewButtonWithLabel("Toggle Wind Simulation")
	var windActive bool // Wind simulation state

	windButton.ConnectClicked(func() {
		windActive = !windActive // Toggle wind simulation
		if windActive {
			windButton.SetLabel("Wind Simulation ON")
		} else {
			windButton.SetLabel("Wind Simulation OFF")
		}
	})

	// Add widgets to the grid
	grid.Attach(forwardLabel, 0, 0, 1, 1)
	grid.Attach(forwardEntry, 1, 0, 1, 1)
	grid.Attach(aftLabel, 0, 1, 1, 1)
	grid.Attach(aftEntry, 1, 1, 1, 1)
	grid.Attach(brLabel, 0, 2, 1, 1)
	grid.Attach(brEntry, 1, 2, 1, 1)
	grid.Attach(startButton, 0, 3, 2, 1)
	grid.Attach(resultLabel, 0, 4, 2, 1)
	grid.Attach(windButton, 0, 5, 2, 1)

	// Set up the simulation logic
	startButton.ConnectClicked(func() {
		forwardText := forwardEntry.Text() // Corrected: Removed the error check
		aftText := aftEntry.Text()         // Corrected: Removed the error check
		brText := brEntry.Text()           // Corrected: Removed the error check

		forward, err := strconv.ParseFloat(forwardText, 64)
		if err != nil {
			resultLabel.SetText("Invalid forward tank value")
			return
		}
		aft, err := strconv.ParseFloat(aftText, 64)
		if err != nil {
			resultLabel.SetText("Invalid aft tank value")
			return
		}
		br, err := strconv.ParseFloat(brText, 64)
		if err != nil {
			resultLabel.SetText("Invalid burn rate value")
			return
		}

		// Disable the button to prevent multiple simulations
		startButton.SetSensitive(false)

		go simulateBurn(forward, aft, br, resultLabel, startButton, &windActive)
	})

	// Set the grid as the child of the window
	window.SetChild(grid)

	// Apply CSS for dark theme
	cssProvider := gtk.NewCSSProvider()
	cssProvider.LoadFromData(`
		.result-label {
			color: white;
		}
		window {
			background-color: #2e2e2e;
		}
	`)
	display := gdk.DisplayGetDefault()
	gtk.StyleContextAddProviderForDisplay(display, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	window.Show()
}

func simulateBurn(forward, aft, br float64, resultLabel *gtk.Label, startButton *gtk.Button, windActive *bool) {
	ticker := time.NewTicker(time.Second) // Changed to seconds for faster simulation
	elapsed := 0
	swapInterval := 5 * 60 // Swap tanks every 5 minutes (300 seconds)
	useForwardTank := true

	for range ticker.C {
		// Check for wind simulation and adjust burn rate
		actualBurnRate := br
		if *windActive {
			actualBurnRate *= 1.2 // Increase burn rate by 20% due to wind
		}

		// Burn fuel from the selected tank
		if useForwardTank {
			forward -= actualBurnRate / 60 // Convert burn rate to lbs/sec
			if forward < 0 {
				forward = 0
			}
		} else {
			aft -= actualBurnRate / 60 // Convert burn rate to lbs/sec
			if aft < 0 {
				aft = 0
			}
		}

		// Increment the elapsed time
		elapsed++

		// Swap tanks based on interval
		if elapsed%swapInterval == 0 {
			useForwardTank = !useForwardTank
		}

		// Calculate total time till out of fuel
		totalFuel := forward + aft
		timeTillEmpty := totalFuel / (actualBurnRate / 60) // in minutes
		if totalFuel == 0 {
			timeTillEmpty = 0
		}

		// Update the result label; must be done in the main thread
		glib.IdleAdd(func() bool {
			resultLabel.SetText(fmt.Sprintf("Forward Tank: %.2f lbs\nAft Tank: %.2f lbs\nTime till empty: %.2f minutes", forward, aft, timeTillEmpty))
			if forward == 0 && aft == 0 {
				ticker.Stop()
				startButton.SetSensitive(true)
				resultLabel.SetText(fmt.Sprintf("Simulation complete\nForward Tank: %.2f lbs\nAft Tank: %.2f lbs", forward, aft))
				return false // No further calls are required
			}
			return true // Continue calling until tanks are empty
		})

		// Stop the simulation if both tanks are empty
		if forward == 0 && aft == 0 {
			break
		}
	}

	// After the simulation is done, re-enable the start button
	glib.IdleAdd(func() bool {
		startButton.SetSensitive(true)
		return false // No further calls are required
	})
}
