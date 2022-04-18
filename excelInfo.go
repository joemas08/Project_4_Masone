package main

import (
	"fmt"
	"github.com/blizzy78/ebitenui"
	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xuri/excelize/v2"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"log"
	"strconv"
)

var populationApp GuiApp
var textWidget *widget.Text

func main() {
	ebiten.SetWindowSize(900, 800)
	ebiten.SetWindowTitle("2020-2021 United States Population Change Data")

	populationApp = GuiApp{AppUI: MakeUIWindow()}

	err := ebiten.RunGame(&populationApp)
	if err != nil {
		log.Fatalln("Error running User Interface Demo", err)
	}
}

func (g GuiApp) Update() error {
	g.AppUI.Update()
	return nil
}

func (g GuiApp) Draw(screen *ebiten.Image) {
	g.AppUI.Draw(screen)
}

func (g GuiApp) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

type GuiApp struct {
	AppUI *ebitenui.UI
}

func MakeUIWindow() (GUIhandler *ebitenui.UI) {
	background := image.NewNineSliceColor(color.Gray16{})
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}),
			widget.GridLayoutOpts.Spacing(0, 20))),
		widget.ContainerOpts.BackgroundImage(background))
	textInfo := widget.TextOptions{}.Text("                                         "+
		"2020-2021 United States Population Change Data", basicfont.Face7x13, color.White)
	resources, err := newListResources()
	if err != nil {
		log.Println(err)
	}
	allStates := loadStates()
	dataAsGeneric := make([]interface{}, len(allStates))
	for position, state := range allStates {
		dataAsGeneric[position] = state
	}

	listWidget := widget.NewList(
		widget.ListOpts.Entries(dataAsGeneric),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			fullName := "%s  -------->  Pop. Change : %d  ---  %% Change of Pop. : %.2f %%"
			fullName = fmt.Sprintf(fullName, e.(State).StateName, e.(State).PopChange, e.(State).PercentChange)
			return fullName
		}),
		widget.ListOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(resources.image)),
		widget.ListOpts.SliderOpts(
			widget.SliderOpts.Images(resources.track, resources.handle),
			widget.SliderOpts.HandleSize(resources.handleSize),
			widget.SliderOpts.TrackPadding(resources.trackPadding)),
		widget.ListOpts.EntryColor(resources.entry),
		widget.ListOpts.EntryFontFace(resources.face),
		widget.ListOpts.EntryTextPadding(resources.entryPadding),
		widget.ListOpts.HideHorizontalSlider(),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			//do something when a list item changes
		}))
	textWidget = widget.NewText(textInfo)
	rootContainer.AddChild(textWidget)
	rootContainer.AddChild(listWidget)

	GUIhandler = &ebitenui.UI{Container: rootContainer}
	return GUIhandler
}

func loadStates() []State {
	sliceOfStates := make([]State, 51, 55)
	excelFile, err := excelize.OpenFile("countyPopChange2020-2021.xlsx")
	if err != nil {
		log.Fatalln(err)
	}

	allRows, err := excelFile.GetRows("co-est2021-alldata")
	if err != nil {
		log.Fatalln(err)
	}
	location := 0
	percentChng := 0.0
	for number, row := range allRows {
		if number < 1 || number == 330 {
			continue
		}
		if row[5] == row[6] {
			statePopChange, _ := strconv.Atoi(row[11])
			popEst2020, _ := strconv.Atoi(row[8])
			percentChng = float64(statePopChange) / float64(popEst2020) * 100
			currentState := State{row[5], statePopChange, percentChng}
			sliceOfStates[location] = currentState
			location++

		}
	}
	return sliceOfStates
}
