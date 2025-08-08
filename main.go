package main

import (
	"LiveBuilder/frontend"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	fmt.Println("Live Builder")
	plc := frontend.NewPackageListContainer()
	myApp := app.New()
	myWindow := myApp.NewWindow("test window")

	cnt := plc.GetContainer()
	fmt.Println(cnt.Size())
	myWindow.SetContent(cnt)

	windowPadding := fyne.NewSize(30, 60)
	windowSize := fyne.NewSize(
		cnt.Size().Width+windowPadding.Width,
		cnt.Size().Height+windowPadding.Height,
	)

	myWindow.Resize(windowSize)

	myWindow.ShowAndRun()

}
