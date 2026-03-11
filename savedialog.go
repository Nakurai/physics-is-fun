package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
)

func makeButtonImage() *widget.ButtonImage {
	idle := image.NewNineSliceColor(color.NRGBA{R: 80, G: 80, B: 80, A: 255})
	hover := image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255})
	pressed := image.NewNineSliceColor(color.NRGBA{R: 60, G: 60, B: 60, A: 255})
	return &widget.ButtonImage{Idle: idle, Hover: hover, Pressed: pressed}
}

func (l *Level4) buildSaveDialog() {
	// Semi-transparent dark overlay (full screen)
	overlay := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
	)

	// The dialog box itself, centered
	dialog := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
			widget.RowLayoutOpts.Padding(&widget.Insets{Top: 20, Bottom: 20, Left: 20, Right: 20}),
		)),
		widget.ContainerOpts.BackgroundImage(
			image.NewNineSliceColor(color.NRGBA{R: 30, G: 30, B: 30, A: 220}),
		),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
	)

	// Text input
	l.saveInput = widget.NewTextInput(
		widget.TextInputOpts.Image(&widget.TextInputImage{
			Idle: image.NewNineSliceColor(color.NRGBA{R: 60, G: 60, B: 60, A: 255}),
		}),
		widget.TextInputOpts.Face(&myFontFace),
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:  color.White,
			Caret: color.White,
		}),
		widget.TextInputOpts.Placeholder("Enter save name..."),
		widget.TextInputOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(200, 40),
		),
	)

	// Save button
	saveBtn := widget.NewButton(
		widget.ButtonOpts.Image(makeButtonImage()),
		widget.ButtonOpts.Text("Save", &myFontFace, &widget.ButtonTextColor{Idle: color.White}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			filename := l.saveInput.GetText() // <-- this is how you collect the input
			if filename != "" {
				l.saveGameToFile(filename + ".json")
				l.hideSaveDialog()
			}
		}),
	)

	dialog.AddChild(l.saveInput, saveBtn)
	overlay.AddChild(dialog)
	l.saveDialog = overlay
}

func (l *Level4) showSaveDialog() {
	l.showingSave = true
	l.rootContainer.AddChild(l.saveDialog) // adds it on top of everything
	l.saveInput.Focus(true)                // immediately focus the text field
}

func (l *Level4) hideSaveDialog() {
	l.showingSave = false
	l.saveInput.SetText("")                   // clear for next time
	l.rootContainer.RemoveChild(l.saveDialog) // remove from render tree
}

func (l *Level4) saveGameToFile(filename string) {
	jsonLevel, err := json.Marshal(struct {
		Obstacles []*Solid `json:"obstacles"`
		Balls     []Ball   `json:"balls"`
		Mode      string   `json:"mode"`
		IsPaused  bool     `json:"isPaused"`
	}{
		Obstacles: l.Obstacles,
		Balls:     l.Balls,
		Mode:      l.Mode,
		IsPaused:  l.IsPaused,
	})
	if err != nil {
		fmt.Printf("ERROR marshalling the game: %v\n", err)
		return
	}
	err = os.WriteFile("./save/"+filename, jsonLevel, 0644)
	if err != nil {
		fmt.Printf("ERROR writing the file: %v\n", err)
	}

}
