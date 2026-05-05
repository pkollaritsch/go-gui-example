package main

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"cogentcore.org/core/colors"
	"cogentcore.org/core/core"
	"cogentcore.org/core/cursors"
	"cogentcore.org/core/events"
	"cogentcore.org/core/events/key"
	"cogentcore.org/core/math32"
	"cogentcore.org/core/paint"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/abilities"
	"cogentcore.org/core/styles/units"
	"cogentcore.org/core/text/htmltext"
	"cogentcore.org/core/text/rich"
	"cogentcore.org/core/text/shaped"
	"cogentcore.org/core/text/text"
)

func main() {
	bodyDim := int(1000)           // width = height in pixels
	num := int(-1)                 // hold number which dialog window requests
	isDialogWindowBlocking := true // boolean holding value of switch position
	strStatus := "STATUS: waiting for input..."
	txtSh := shaped.NewShaper() // for shaping text drawn to canvas

	//######### Create Body and its body parts: #######################################################
	b := core.NewBody("Cogent Core Example")

	canvas := core.NewCanvas(b)

	core.NewText(b).SetText("Click inside canvas, and then press 'd' key to bring up dialog window\n which is (blocking or non-blocking):")

	blockingSwitch := core.NewSwitch(b).SetText("Blocking")
	core.Bind(&isDialogWindowBlocking, blockingSwitch) //bind switch value to affect isDialogWindowBlocking

	textStatusWidget := core.NewText(b).SetText(strStatus)
	core.Bind(&strStatus, textStatusWidget) //bind our strStatus var to the value this text widget

	//######### Add details to Body and its body parts: #######################################################
	b.Styler(func(s *styles.Style) {
		s.Background = colors.Uniform(colors.Lightblue)
		//Sizes body window but not the system window
		//Use Dot for actual pixels... do not use Px which equals 1/96 inch
		s.Min.X.Dot(float32(bodyDim) + 20) // Set minimum width  (+20 to account for scroll bar width)
		s.Min.Y.Dot(float32(bodyDim))      // Set minimum height
		s.Max.X.Dot(float32(bodyDim) + 20)
		s.Max.Y.Dot(float32(bodyDim))
	})

	canvas.Styler(func(s *styles.Style) {
		// All dimensions must be equal or circles look like ovals
		canvasDim := bodyDim - 400
		s.Min.X.Dot(float32(canvasDim)) // Set minimum width
		s.Min.Y.Dot(float32(canvasDim)) // Set minimum height
		s.Max.X.Dot(float32(canvasDim))
		s.Max.Y.Dot(float32(canvasDim))
		s.SetAbilities(true, abilities.Focusable) // MUST enable focus for the canvas to receive key events
	})

	canvas.SetDraw(func(pc *paint.Painter) {
		sz := canvas.Geom.Size.Actual.Content
		originalMatrix := pc.Transform
		//by default coordinates are interpreted by draw funct to be between 0 and 1.0, 1.0 meaning the max dimension of the canvas
		pc.Transform = pc.Transform.Scale(1/sz.X, 1/sz.Y) // Scale transform so all coordinates are interpreted as pixel dimensions

		pc.FillBox(math32.Vec2(0, 0), math32.Vec2(sz.X, sz.Y), colors.Uniform(colors.White))
		pc.FillBox(math32.Vec2(40, 80), math32.Vec2(150, 75), colors.Uniform(colors.Darkcyan))
		pc.Fill.Color = colors.Uniform(colors.Yellow)
		pc.Circle(100, 100, 50) // DrawCircle uses (centerX, centerY, radius)
		pc.Stroke.Color = colors.Uniform(colors.Blue)
		pc.Stroke.Width.Dot(3)
		pc.Stroke.Dashes = []float32{}                         //Solid line
		CanvasDrawLine(pc, 0, 0, float64(sz.X), float64(sz.Y)) //Draw solid backslash line
		pc.Stroke.Color = colors.Uniform(colors.Orange)
		pc.Stroke.Dashes = []float32{18, 8}                    //small dotted line
		CanvasDrawLine(pc, float64(sz.X), 0, 0, float64(sz.Y)) //Draw dotted slash line
		strPixelWid := float64(500)
		strPixelHgt := float64(500)
		fontHgtInPixels := float32(18)
		myStr := "This is <a>HTML</a> <b>formatted</b> <i>text</i>"
		CanvasDrawString(pc, txtSh, myStr, 100, 200, strPixelWid, strPixelHgt, fontHgtInPixels, colors.Red)

		pc.Draw()
		pc.Transform = originalMatrix // Restore the matrix
	})

	//##################### Click Events ###############################################################
	b.On(events.MouseDown, func(e events.Event) {
		e.SetHandled() // Prevents the click from being ignored or passed on
		strStatus = "STATUS: Body focused via click"
		textStatusWidget.Update()
		textStatusWidget.SetFocus() //Needed so that canvas does not have focus
	})

	canvas.On(events.MouseDown, func(e events.Event) {
		e.SetHandled() // Prevents the click from being ignored or passed on
		strStatus = "STATUS: Canvas focused via click"
		textStatusWidget.Update()
		canvas.SetFocus()
	})

	blockingSwitch.OnChange(func(e events.Event) {
		strStatus = fmt.Sprintf("STATUS:  Blocking switch set to: %v\n", isDialogWindowBlocking)
		textStatusWidget.Update()
	})

	//##################### Key Press and Other Events ###############################################################
	canvas.On(events.KeyChord, func(e events.Event) {
		if e.KeyChord() == key.Chord("d") { // 'd' for dialog window
			num = -1
			//######## Create body and body parts ###################
			d := core.NewBody("Input Positive Integer") // Create new dialog window ('d' for dialog)
			input := core.NewTextField(d).SetPlaceholder("Enter number...")
			buttonRow := core.NewFrame(d)
			buttonRow.Styler(func(s *styles.Style) {
				s.Direction = styles.Row
				s.Gap.Set(units.Dp(8))
			})
			btnSubmit := core.NewButton(buttonRow).SetText("Submit")
			btnClose := core.NewButton(buttonRow).SetText("Close")
			txtClose := core.NewText(d).SetText("Click the button above to exit or click here.")
			txtClose.Styler(func(s *styles.Style) {
				s.Font.SetDecoration(rich.Underline)
				s.Color = colors.Uniform(colors.Blue)
				s.SetAbilities(true, abilities.Clickable)
				s.Cursor = cursors.Arrow //So cursor icon does not change to an I-bar icon
			})
			//######## Events #######################################
			d.OnShow(func(e events.Event) {
				input.SetFocus()
			})
			btnClose.OnClick(func(e events.Event) {
				num = -1
				strStatus = fmt.Sprintf("STATUS:  Dialog window's Close button clicked. num = %d\n", num)
				textStatusWidget.Update()
				d.Close()
			})
			txtClose.OnClick(func(e events.Event) {
				num = -1
				strStatus = fmt.Sprintf("STATUS:  Dialog window's Exit text clicked. num = %d\n", num)
				textStatusWidget.Update()
				d.Close()
			})
			btnSubmit.OnClick(func(e events.Event) {
				val := input.Text()
				val_i, err := strconv.Atoi(val)
				if err != nil || val_i < 0 {
					// Show an error message if the input is bad
					core.MessageSnackbar(d, "Invalid input! Please enter 0 or higher.")
					input.SetText("") // Clear it
					return
				}
				num = val_i
				strStatus = fmt.Sprintf("STATUS:  Dialog window's Submit button clicked. num = %d\n", num)
				textStatusWidget.Update()
				d.Close() // Close window on success
			})
			//######## Run Dialog Window ###############################
			if isDialogWindowBlocking == true {
				d.RunDialog(b) //blocks the parent, b, from receiving focus until this child window, d, is dismissed
			} else {
				d.RunWindow() //creates a peer window that can still be moved behind the first.
			}
		}
	})

	//##################  STARTUP MAIN WINDOW #######################
	b.RunMainWindow()
}

func CanvasDrawLine(pc *paint.Painter, x1 float64, y1 float64, x2 float64, y2 float64) {
	pc.MoveTo(float32(x1), float32(y1))
	pc.LineTo(float32(x2), float32(y2))
	pc.Draw() // Must have to complete the line
}

func CanvasDrawString(pc *paint.Painter, txtSh shaped.Shaper, myStr string, x float64, y float64, strPixelWid float64, strPixelHgt float64,
	fontHgtInPixels float32, color color.Color) {
	if strings.Contains(myStr, " ") {
		myStr = strings.ReplaceAll(myStr, " ", "&nbsp;&nbsp;")
	}
	tsty := text.NewStyle()
	fsty := rich.NewStyle()
	tsty.FontSize.Dp(fontHgtInPixels) //Use this since it doesn't honor:  pc.Context().Style.Font.Size = 18
	tsty.Color = color
	tsty.ToDots(&pc.UnitContext)
	tx, err := htmltext.HTMLToRich([]byte(myStr), fsty, nil)
	if err != nil {
		panic(err)
	}
	lns := txtSh.WrapLines(tx, fsty, tsty, math32.Vec2(float32(strPixelWid), float32(strPixelHgt)))
	pc.DrawText(lns, math32.Vec2(float32(x), float32(y)))
}
