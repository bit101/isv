// Package main runs the app.
package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/bit101/isv/res"
	flag "github.com/spf13/pflag"
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
	version  = "v0.1.0"
	stopping = false
	fpsIndex = 3
	fpsList  = []time.Duration{1, 5, 10, 30, 60, 1000}
	delay    = 1000 / fpsList[fpsIndex]

	watchTime time.Duration = 4

	entries     []string
	img         *canvas.Image
	w           fyne.Window
	modeLabel   *canvas.Text
	watchLabel  *canvas.Text
	initPlay    bool
	initReverse bool
	initBounce  bool
	initWatch   int
	initHelp    bool
	initVersion bool
)

func init() {
	flag.BoolVarP(&initPlay, "play", "p", false, "plays the image sequence on start")
	flag.BoolVarP(&initReverse, "reverse", "r", false, "plays the image sequence in reverse on start")
	flag.BoolVarP(&initBounce, "bounce", "b", false, "plays the image sequence back and forth on start")
	flag.IntVarP(&initWatch, "watch", "w", 0, "rescans dir every n (1-10) seconds")
	flag.BoolVarP(&initHelp, "help", "h", false, "shows this help")
	flag.BoolVarP(&initVersion, "version", "v", false, "shows the version number")

	flag.Usage = func() {
		fmt.Print("Usage:\n  isv [options] directory_path\nOptions:\n")
		flag.PrintDefaults()
		fmt.Println("Keys:")
		fmt.Println("  left/right: prev/next frame")
		fmt.Println("  f/l: first/last frame")
		fmt.Println("  p: play")
		fmt.Println("  r: reverse")
		fmt.Println("  b: bounce")
		fmt.Println("  space: stop")
		fmt.Println("  up/down: speed")
		fmt.Println("  w: watch")
		fmt.Println("  < / >: watch interval")
		fmt.Println("  F5: manual refresh image list")
		fmt.Println("  Q/Esc: quit")
	}
	flag.Parse()
	dir = flag.Arg(0)
	if dir == "" {
		dir = "."
	}
}

func main() {
	// make gui
	a := app.New()
	w = a.NewWindow("isv")
	w.SetFixedSize(true)

	img = &canvas.Image{}
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(400, 400))

	modeLabel = canvas.NewText("stopped", color.Black)
	modeLabel.TextSize = 10

	watchLabel = canvas.NewText("not watching", color.Black)
	watchLabel.TextSize = 10

	labels := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), modeLabel, layout.NewSpacer(), watchLabel, layout.NewSpacer())

	w.SetContent(container.New(layout.NewVBoxLayout(), img, labels))

	// events
	w.Canvas().SetOnTypedKey(handleKeys)

	if initHelp {
		flag.Usage()
		os.Exit(0)
	}

	if initVersion {
		fmt.Printf("isv version: %s\n", version)
		os.Exit(0)
	}

	// load image list and first image
	readDir()
	if initPlay {
		go animate()
	} else if initReverse {
		go reverse()
	} else if initBounce {
		go bounce()
	} else {
		loadImage()
	}

	if initWatch > 0 {
		watchTime = time.Duration(initWatch)
		if watchTime > 10 {
			watchTime = 10
		}
		watchDir = true
		updateWatchLabel()
		go watch()
	}

	w.ShowAndRun()
}

func loadImage() {
	// no entries, nothing to load.
	// read list again and jump.
	if len(entries) == 0 {
		index = 0
		readDir()
		img.Resource = res.Placeholder()
		return
	}

	// probably deleted some images since last load.
	// read list again and start again at 0.
	if index >= len(entries) {
		readDir()
		index = 0
	}

	name := entries[index]
	filepath := path.Join(dir, name)

	// make sure it exists before loading it (could have been deleted since last check)
	if _, err := os.Stat(filepath); err == nil {
		w.SetTitle(name)
		img.File = filepath
		img.Refresh()
	} else {
		// make sure we're up to date.
		readDir()
	}
}

func updateFPS() {
	fps := fpsList[fpsIndex]
	text := "stopped"

	if mode == Forward {
		text = "forward"
	} else if mode == Reverse {
		text = "reverse"
	} else if mode == Bounce {
		text = "bounce"
	}

	if mode == Stopped {
		modeLabel.Text = text
	} else if fps < 1000 {
		modeLabel.Text = fmt.Sprintf("%s: %d fps (attempted)", text, fps)
	} else {
		modeLabel.Text = fmt.Sprintf("%s: max fps", text)
	}
	modeLabel.Refresh()
}

func updateWatchLabel() {
	if watchDir {
		watchLabel.Text = fmt.Sprintf("watching: %ds", watchTime)
	} else {
		watchLabel.Text = "not watching"
	}
	watchLabel.Refresh()
}

func animate() {
	mode = Forward
	updateFPS()
	for mode == Forward {
		loadImage()
		index++
		// loop back to start
		if index >= len(entries) {
			index = 0
		}
		time.Sleep(delay * time.Millisecond)
	}
	if mode == Stopped {
		stopping = false
		updateFPS()
	}
}

func reverse() {
	mode = Reverse
	updateFPS()
	for mode == Reverse {
		loadImage()
		index--
		// loop back to end
		if index < 0 {
			index = len(entries) - 1
		}
		time.Sleep(delay * time.Millisecond)
	}
	if mode == Stopped {
		stopping = false
		updateFPS()
	}
}

func bounce() {
	direction := 1
	mode = Bounce
	updateFPS()
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
	if mode == Stopped {
		stopping = false
		updateFPS()
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
	entries = filterImages(list)
}

func filterImages(list []os.DirEntry) []string {
	names := []string{}
	for _, f := range list {
		name := f.Name()
		if strings.HasSuffix(name, ".png") ||
			strings.HasSuffix(name, ".jpg") ||
			strings.HasSuffix(name, ".jpeg") {
			names = append(names, name)
		}
	}
	return names
}

func handleKeys(k *fyne.KeyEvent) {
	// quit - Q or ESC
	if k.Name == fyne.KeyEscape || k.Name == fyne.KeyQ {
		w.Close()
	}

	// next frame - right arrow
	if k.Name == fyne.KeyRight && mode == Stopped {
		index++
		if index >= len(entries) {
			index = 0
		}
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
			stopping = true
		} else if !stopping {
			go animate()
		}
	}

	// play reverse - R
	if k.Name == fyne.KeyR {
		if mode == Reverse {
			mode = Stopped
			stopping = true
		} else if !stopping {
			go reverse()
		}
	}

	// play bounce - B
	if k.Name == fyne.KeyB {
		if mode == Bounce {
			mode = Stopped
			stopping = true
		} else if !stopping {
			go bounce()
		}
	}

	// stop playing - SPACE
	if k.Name == fyne.KeySpace {
		mode = Stopped
	}

	// increase animation speed - up arrow
	if k.Name == fyne.KeyUp {
		fpsIndex++
		if fpsIndex >= len(fpsList) {
			fpsIndex = len(fpsList) - 1
		}
		delay = 1000 / fpsList[fpsIndex]
		updateFPS()
	}

	// decrease animation speed - down arrow
	if k.Name == fyne.KeyDown {
		fpsIndex--
		if fpsIndex < 0 {
			fpsIndex = 0
		}
		delay = 1000 / fpsList[fpsIndex]
		updateFPS()
	}

	// watch dir - W
	if k.Name == fyne.KeyW {
		watchDir = !watchDir
		if watchDir {
			go watch()
		}
		updateWatchLabel()
	}

	// decrease watch time - comma (<)
	if k.Name == fyne.KeyComma {
		watchTime--
		if watchTime < 1 {
			watchTime = 1
		}
		updateWatchLabel()
	}

	// increase watch time - comma (>)
	if k.Name == fyne.KeyPeriod {
		watchTime++
		if watchTime > 10 {
			watchTime = 10
		}
		updateWatchLabel()
	}

	// manually refresh image list - F5
	if k.Name == fyne.KeyF5 {
		readDir()
	}
}
