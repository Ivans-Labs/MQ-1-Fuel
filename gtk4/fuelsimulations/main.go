package main

import (
	"fmt"
	"math"
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
	window.SetDefaultSize(800, 300) // Increased width to accommodate warning system

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

	// Create a button to simulate ice on wings
	iceButton := gtk.NewButtonWithLabel("Toggle Ice Simulation")
	var iceActive bool // Ice on wings simulation state

	// Create a button to simulate degraded engine performance
	engineButton := gtk.NewButtonWithLabel("Toggle Engine Degradation")
	var engineDegraded bool // Degraded engine simulation state

	// Create a label to display warnings and cautions
	warningLabel := gtk.NewLabel("")
	warningLabel.SetCSSClasses([]string{"warning-label"}) // Add CSS class for styling

	// Add widgets to the grid
	grid.Attach(forwardLabel, 0, 0, 1, 1)
	grid.Attach(forwardEntry, 1, 0, 1, 1)
	grid.Attach(aftLabel, 0, 1, 1, 1)
	grid.Attach(aftEntry, 1, 1, 1, 1)
	grid.Attach(brLabel, 0, 2, 1, 1)
	grid.Attach(brEntry, 1, 2, 1, 1)
	grid.Attach(startButton, 0, 3, 2, 1)
	grid.Attach(resultLabel, 0, 4, 2, 1)
	grid.Attach(windButton, 0, 5, 1, 1)
	grid.Attach(iceButton, 1, 5, 1, 1)
	grid.Attach(engineButton, 0, 6, 2, 1)
	grid.Attach(warningLabel, 2, 0, 1, 7) // Span 7 rows for warnings

	// Set up the simulation logic
	startButton.ConnectClicked(func() {
		forwardText := forwardEntry.Text()
		aftText := aftEntry.Text()
		brText := brEntry.Text()

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

		go simulateBurn(forward, aft, br, resultLabel, startButton, &windActive, &iceActive, &engineDegraded, warningLabel)
	})

	// Toggle wind simulation
	windButton.ConnectClicked(func() {
		windActive = !windActive // Toggle wind simulation
		if windActive {
			windButton.SetLabel("Wind Simulation ON")
		} else {
			windButton.SetLabel("Wind Simulation OFF")
		}
	})

	// Toggle ice simulation
	iceButton.ConnectClicked(func() {
		iceActive = !iceActive // Toggle ice simulation
		if iceActive {
			iceButton.SetLabel("Ice Simulation ON")
		} else {
			iceButton.SetLabel("Ice Simulation OFF")
		}
	})

	// Toggle engine degradation
	engineButton.ConnectClicked(func() {
		engineDegraded = !engineDegraded // Toggle engine degradation
		if engineDegraded {
			engineButton.SetLabel("Engine Degradation ON")
		} else {
			engineButton.SetLabel("Engine Degradation OFF")
		}
	})

	// Set the grid as the child of the window
	window.SetChild(grid)

	// Apply CSS for dark theme and white text
	cssProvider := gtk.NewCSSProvider()
	cssProvider.LoadFromData(`
	window {
		background-color: #2e2e2e;
		color: white; /* This will set the default text color to white for the window */
	}
	label {
		color: white; /* This ensures all labels have white text */
	}
	entry {
		color: black; /* Input text color for entries */
		background-color: white; /* Background color for entries */
	}
	button {
		color: white; /* Text color for buttons */
		background-color: black; /* Background color for buttons */
	}
	.warning-label {
		color: red;
		font-weight: bold;
		border-radius: 5px;
		border: 1px solid red;
		padding: 5px;
		margin-top: 10px;
	}
`)
	display := gdk.DisplayGetDefault()
	gtk.StyleContextAddProviderForDisplay(display, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	window.Show()
}

func simulateBurn(forward, aft, br float64, resultLabel *gtk.Label, startButton *gtk.Button, windActive *bool, iceActive *bool, engineDegraded *bool, warningLabel *gtk.Label) {
	ticker := time.NewTicker(time.Second) // Ticks every second for real-time simulation
	defer ticker.Stop()

	for range ticker.C {
		actualBurnRate := br / 60 // Convert burn rate to lbs per second
		if *windActive {
			actualBurnRate *= 1.2 // Increase burn rate by 20% due to wind
		}
		if *iceActive {
			actualBurnRate *= 1.1 // Increase burn rate by 10% due to ice
		}
		if *engineDegraded {
			actualBurnRate *= 1.3 // Increase burn rate by 30% due to degraded engine
		}

		burnAmount := actualBurnRate / 2 // Split the burn rate between the two tanks per second
		forward -= burnAmount
		aft -= burnAmount

		if forward < 0 {
			forward = 0
		}
		if aft < 0 {
			aft = 0
		}

		totalFuel := forward + aft
		timeTillEmpty := math.Max(forward, aft) / (actualBurnRate * 60) // in hours

		hoursLeft := int(timeTillEmpty)
		minutesLeft := int((timeTillEmpty - float64(hoursLeft)) * 60)

		glib.IdleAdd(func() bool {
			resultLabel.SetText(fmt.Sprintf("Forward Tank: %.2f lbs\nAft Tank: %.2f lbs\nTotal Fuel: %.2f lbs\nTime till empty: %02d hours and %02d minutes", forward, aft, totalFuel, hoursLeft, minutesLeft))
			if forward == 0 || aft == 0 {
				ticker.Stop()
				startButton.SetSensitive(true)
				resultLabel.SetText(fmt.Sprintf("Simulation complete\nForward Tank: %.2f lbs\nAft Tank: %.2f lbs\nTotal Fuel: %.2f lbs", forward, aft, totalFuel))
				return false // No further calls are required
			}
			return true // Continue calling until one of the tanks is empty
		})

		glib.IdleAdd(func() bool {
			warningText := ""
			if *windActive {
				warningText += "Caution: Wind Active\n"
			}
			if *iceActive {
				warningText += "Warning: Ice on Wings\n"
			}
			if *engineDegraded {
				warningText += "Warning: Engine Degradation\n"
			}
			warningLabel.SetText(warningText)
			return true // Continue updating warnings
		})
	}

	// After the simulation is done, re-enable the start button
	glib.IdleAdd(func() bool {
		startButton.SetSensitive(true)
		return false // No further calls are required
	})
}
