// Package main runs the app.
package main

import (
	"log"
	"os"
	"path"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

// AnimMode represents the animation state.
type AnimMode int

// AnimModes
const (
	Stopped AnimMode = iota
	Forward
	Reverse
	Bounce
)

var (
	mode     = Stopped
	watchDir = false
	index    = 0
	dir      = "."

	delay     time.Duration = 30
	watchTime time.Duration = 5

	entries []os.DirEntry
	img     *canvas.Image
	w       fyne.Window
)

func main() {
	// make gui
	a := app.New()
	w = a.NewWindow("Images")
	img = &canvas.Image{}
	img.FillMode = canvas.ImageFillOriginal
	w.SetContent(img)

	// events
	w.Canvas().SetOnTypedKey(handleKeys)

	// load image list and first image
	readDir()
	loadImage()

	w.ShowAndRun()
}

func loadImage() {
	// no entries, nothing to load.
	if len(entries) == 0 {
		index = 0
		return
	}

	// probably deleted some images since last load.
	if index >= len(entries) {
		index = 0
	}

	name := entries[index].Name()
	filepath := path.Join(dir, name)

	// make sure it exists before loading it (could have been deleted since last check)
	if _, err := os.Stat(filepath); err == nil {
		w.SetTitle(name)
		img.File = filepath
		img.Refresh()
	}
}

func animate() {
	for mode == Forward {
		loadImage()
		index++
		// loop back to start
		if index >= len(entries) {
			index = 0
		}
		time.Sleep(delay * time.Millisecond)
	}
}

func reverse() {
	for mode == Reverse {
		loadImage()
		index--
		// loop back to end
		if index < 0 {
			index = len(entries) - 1
		}
		time.Sleep(delay * time.Millisecond)
	}
}

func bounce() {
	direction := 1
	for mode == Bounce {
		loadImage()
		index += direction
		// go the other way
		if index >= len(entries) || index < 0 {
			direction *= -1
			index += direction
		}
		time.Sleep(delay * time.Millisecond)
	}
}

func watch() {
	for watchDir {
		readDir()
		time.Sleep(watchTime * time.Second)
	}
}

func readDir() {
	list, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	entries = list
}

func handleKeys(k *fyne.KeyEvent) {
	// quit - Q or ESC
	if k.Name == fyne.KeyEscape || k.Name == fyne.KeyQ {
		w.Close()
	}

	// next frame - right arrow
	if k.Name == fyne.KeyRight && mode == Stopped {
		index++
		index %= len(entries)
		loadImage()
	}

	// prev frame - left arrow
	if k.Name == fyne.KeyLeft && mode == Stopped {
		index--
		if index < 0 {
			index += len(entries)
		}
		loadImage()
	}

	// first frame - F
	if k.Name == fyne.KeyF && mode == Stopped {
		index = 0
		loadImage()
	}

	// last frame - L
	if k.Name == fyne.KeyL && mode == Stopped {
		index = len(entries) - 1
		loadImage()
	}

	// play forward - P
	if k.Name == fyne.KeyP {
		if mode == Forward {
			mode = Stopped
		} else {
			mode = Forward
			go animate()
		}
	}

	// play reverse - R
	if k.Name == fyne.KeyR {
		if mode == Reverse {
			mode = Stopped
		} else {
			mode = Reverse
			go reverse()
		}
	}

	// play bounce - B
	if k.Name == fyne.KeyB {
		if mode == Bounce {
			mode = Stopped
		} else {
			mode = Bounce
			go bounce()
		}
	}

	// stop playing - SPACE
	if k.Name == fyne.KeySpace {
		mode = Stopped
	}

	// increase animation speed - up arrow
	if k.Name == fyne.KeyUp {
		delay -= 10
		if delay < 0 {
			delay = 0
		}
	}

	// decrease animation speed - down arrow
	if k.Name == fyne.KeyDown {
		delay += 10
		if delay > 1000 {
			delay = 1000
		}
	}

	// 30 fps - 0
	if k.Name == fyne.Key0 {
		delay = 30
	}

	// 1 fps - 1
	if k.Name == fyne.Key1 {
		delay = 1000
	}

	// 2 fps - 2
	if k.Name == fyne.Key2 {
		delay = 1000 / 2
	}

	// 3 fps - 3
	if k.Name == fyne.Key3 {
		delay = 1000 / 3
	}

	// 4 fps - 4
	if k.Name == fyne.Key4 {
		delay = 1000 / 4
	}

	// 5 fps - 5
	if k.Name == fyne.Key5 {
		delay = 1000 / 5
	}

	// 6 fps - 6
	if k.Name == fyne.Key6 {
		delay = 1000 / 6
	}

	// 7 fps - 7
	if k.Name == fyne.Key7 {
		delay = 1000 / 7
	}

	// 8 fps - 8
	if k.Name == fyne.Key8 {
		delay = 1000 / 8
	}

	// 9 fps - 9
	if k.Name == fyne.Key9 {
		delay = 1000 / 9
	}

	// watch dir - W
	if k.Name == fyne.KeyW {
		watchDir = !watchDir
		if watchDir {
			go watch()
		}
	}

	// decrease watch time - comma (<)
	if k.Name == fyne.KeyComma {
		watchTime--
		if watchTime < 0 {
			watchTime = 0
		}
	}

	// increase watch time - comma (>)
	if k.Name == fyne.KeyPeriod {
		watchTime++
		if watchTime > 10 {
			watchTime = 10
		}
	}

	// manually refresh image list - F5
	if k.Name == fyne.KeyF5 {
		readDir()
	}
}
