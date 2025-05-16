package main

import (
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("song-dl")

	myWindow.Resize(fyne.NewSize(600, 100))
	myWindow.CenterOnScreen()

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Enter URL here...")

	dirEntry := widget.NewEntry()
	dirEntry.SetPlaceHolder("Enter custom directory name here...")

	extensions := []string{"wav", "mp3"}
	extensionDropdown := widget.NewSelect(extensions, func(selected string) {})
	extensionDropdown.Selected = extensions[0]

	dirItem := widget.NewFormItem("Directory Name:", dirEntry)
	dirForm := widget.NewForm(dirItem)
	dirForm.Hide()

	urlItem := widget.NewFormItem("URL:", urlEntry)
	extensionItem := widget.NewFormItem("Extension:", extensionDropdown)

	customDirCheck := widget.NewCheck("Custom Directory", func(b bool) {
		dirForm.Hidden = !b
	})

	randomizeCheck := widget.NewCheck("Randomize", func(b bool) {})
	reverseCheck := widget.NewCheck("Reverse", func(b bool) {})

	checksBox := container.NewHBox(customDirCheck, randomizeCheck, reverseCheck)

	downloadButton := widget.NewButton("Download", func() {
		url := urlEntry.Text
		selectedExtension := extensionDropdown.Selected
		useCustomDir := customDirCheck.Checked
		shouldRandomize := randomizeCheck.Checked
		shouldReverse := reverseCheck.Checked
		dir := "%(playlist)s"

		if useCustomDir {
			dir = dirEntry.Text
		}

		cmd := exec.Command("yt-dlp", "-x", "--audio-format", selectedExtension, "-o", dir+"/%(autonumber)s - %(title)s.%(ext)s", url)

		if shouldRandomize {
			cmd.Args = append(cmd.Args, "--playlist-random")
		}

		if shouldReverse {
			cmd.Args = append(cmd.Args, "--playlist-reverse")
		}

		var b bool

		go func() {
			err := cmd.Run()
			myApp.Driver().DoFromGoroutine(func() {
				if err != nil {
					dialog.ShowError(err, myWindow)
				} else {
					dialog.ShowInformation("Download Complete", "Download Complete!", myWindow)
				}
			}, b)
		}()
	})

	myWindow.SetContent(container.NewVBox(dirForm, widget.NewForm(urlItem, extensionItem), checksBox, downloadButton))
	myWindow.ShowAndRun()
}
