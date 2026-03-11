package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

func (l *Level4) buildLoadDialog() {
	_, screenH := ebiten.WindowSize()
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

	// Scrollable list of save files
	saveList := widget.NewList(
		widget.ListOpts.Entries(l.listSaveFiles()),
		widget.ListOpts.ScrollContainerImage(&widget.ScrollContainerImage{
			Idle: image.NewNineSliceColor(color.NRGBA{R: 20, G: 20, B: 20, A: 255}),
			Mask: image.NewNineSliceColor(color.NRGBA{R: 20, G: 20, B: 20, A: 255}),
		}),
		widget.ListOpts.SliderParams(&widget.SliderParams{
			TrackImage: &widget.SliderTrackImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{R: 50, G: 50, B: 50, A: 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{R: 70, G: 70, B: 70, A: 255}),
			},
			HandleImage: &widget.ButtonImage{
				Idle:    image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
				Hover:   image.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 130, A: 255}),
				Pressed: image.NewNineSliceColor(color.NRGBA{R: 80, G: 80, B: 80, A: 255}),
			},
		}),
		widget.ListOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.MinSize(300, screenH/2),
			),
		),
		widget.ListOpts.EntryFontFace(&myFontFace),
		widget.ListOpts.EntryColor(&widget.ListEntryColor{
			Selected:                   color.NRGBA{R: 255, G: 255, B: 255, A: 255},
			Unselected:                 color.NRGBA{R: 180, G: 180, B: 180, A: 255},
			SelectedBackground:         color.NRGBA{R: 60, G: 100, B: 160, A: 255},
			SelectingBackground:        color.NRGBA{R: 50, G: 80, B: 130, A: 255},
			DisabledUnselected:         color.NRGBA{R: 100, G: 100, B: 100, A: 255},
			DisabledSelected:           color.NRGBA{R: 100, G: 100, B: 100, A: 255},
			DisabledSelectedBackground: color.NRGBA{R: 40, G: 40, B: 40, A: 255},
		}),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			return e.(string)
		}),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			filename := args.Entry.(string)
			l.loadGameFromFile(filename)
			l.hideLoadDialog()
		}),
	)

	// Cancel button
	cancelBtn := widget.NewButton(
		widget.ButtonOpts.Image(makeButtonImage()),
		widget.ButtonOpts.Text("Cancel", &myFontFace, &widget.ButtonTextColor{Idle: color.White}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			l.hideLoadDialog()
		}),
	)

	dialog.AddChild(saveList, cancelBtn)
	overlay.AddChild(dialog)
	l.loadDialog = overlay
}

func (l *Level4) showLoadDialog() {
	l.showingLoad = true
	l.buildLoadDialog()                    // rebuild each time to refresh the file list
	l.rootContainer.AddChild(l.loadDialog) // adds it on top of everything
}

func (l *Level4) hideLoadDialog() {
	l.showingLoad = false
	l.rootContainer.RemoveChild(l.loadDialog)
}

// listSaveFiles reads ./save/ and returns a slice of filenames as []interface{}
func (l *Level4) listSaveFiles() []interface{} {
	entries, err := os.ReadDir("./save")
	if err != nil {
		fmt.Printf("ERROR reading save folder: %v\n", err)
		return []interface{}{}
	}

	var files []interface{}
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files
}

func (l *Level4) loadGameFromFile(filename string) {
	data, err := os.ReadFile("./save/" + filename)
	if err != nil {
		fmt.Printf("ERROR reading file %s: %v\n", filename, err)
		return
	}

	var state struct {
		Obstacles []*Solid `json:"obstacles"`
		Balls     []Ball   `json:"balls"`
		Mode      string   `json:"mode"`
		IsPaused  bool     `json:"isPaused"`
	}

	if err := json.Unmarshal(data, &state); err != nil {
		fmt.Printf("ERROR unmarshalling save file %s: %v\n", filename, err)
		return
	}

	l.Obstacles = state.Obstacles
	l.Balls = state.Balls
	l.Mode = state.Mode
	l.IsPaused = state.IsPaused

	fmt.Printf("Game loaded from %s\n", filename)
}
