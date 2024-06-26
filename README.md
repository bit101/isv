# isv

An Image Sequence Viewer

## Use case

You are generating up to several hundred image frames which will later be compiled into an animated gif or video. The generation of frames may take many minutes. You'd like to see how the animation is shaping up during the rendering process.

## What isv does

### Viewing images

Open `isv` in the folder where the images are located, or run `isv path/to/image/folder` from another location. The first frame will be displayed.

Press the right arrow key to advance through subsequent frames. Press the left arrow key to go backwards through the frames. It will loop when it reaches the end of the list in either direction.

OK, most image viewers will give you the above functionaliy. What else?

### Animating

Press `p` and the frames will play back in sequence. Press `r` and they'll play back in reverse. Press `b` and they will bounce - play forward to the last frame then reverse to the first frame, etc. Bounce is particularly useful during the early stages of a render where not too many frames are complete and the sudden jump back to the first frame can be jarring. The bounce effect looks a lot smoother.

Press the same key again to turn off animation, or just press the space bar.

`isv` will attempt to run this animation at 30 fps. How well it does that depends on your hardware and the size of the images. I've found it does well with images up to 200-300 kb, and starts to gett a little stuttery with anything much larger. Note that `isv` is NOT meant as a general purpose image viewer and is not at all optimized for larger images.

You can use the up arrow key to increase the playback speed, which will eventually remove any delay between frames and play them back as fast as your computer can handle them. Press the down arrow key to slow down the speed. Minimum speed is 1 fps. 

### Watching the image list

Again, the main use case is previewing an image sequence that is currently being rendered. By default, `isv` reads the directory on startup and caches the image list. But as new frames come in, you want to update that list from time to time.

Press `w` to start "watching" the directory. Technically, it's just reloading every 4 seconds by default. Press `<` (`,`) to increase the frequency of reloading the list by one second increments, down to once per second. Press `>` (`.`) to increase it by one second increments up to once every ten seconds. Press `w` again to stop watching.

### Cheat sheet

- `p` play forward
- `r` play reverse
- `b` play bounce
- `right arrow` next frame
- `left arrow` prev frame
- `f` first frame
- `l` last frame
- `up arrow` increase playback rate
- `down arrow` decrease playback rate
- `w` to start watching the directory
- `<` (or `,`) to increase refresh rate
- `>` (or `.`) to decrease refresh rate
- `space` stop
- `q` or `esc` quit

### Command line usage

Although `isv` is a graphical ui app, it is meant to be run as a command line app. At this point there is no visible ui at all, only the keyboard shortcuts above. The usage for the tool is:

```
Usage:
  isv [options] directory_path

Options:
  -b, --bounce      plays the image sequence back and forth on start
  -h, --help        shows this help
  -p, --play        plays the image sequence on start
  -r, --reverse     plays the image sequence in reverse on start
  -v, --version     shows the version number
  -w, --watch int   rescans dir every n (1-10) seconds
```
The default directory path if not specified is the current directory. 

So you can set it to play and watch the directory every 2 seconds by running `isv -p -w 2`. Or even `isv -pw 2`.

## Limitations

As already mentioned, `isv` is NOT meant to be a general purpose image viewer. It is not optimized for large images, only smaller frames of the size you might make a gif or movie for posting on social media.

Any stated frame rates are only an approximation and not guaranteed.

`isv` only previews `jpg`, `jpeg` and `png` files currently.

## Installation

### Locally

0. Install Go.
1. Check out the repo.
2. Run `go mod tidy` to update dependencies.
3. Run `go build` to create the binary.
4. Or run `go install` to install the binary on your system.

You may still need additional tool chain items, like gcc or clang installed.

### Go install

0. Install Go.
1. Run `go install github.com/bit101/isv`

You may still need additional tool chain items, like gcc or clang installed.

### Binaries

The precompiled binaries should work for all of the targeted systems - Linux amd64, Windows amd64, MacOS arm64. No other dependencies required.

Download binary for your platform from [https://github.com/bit101/isv/releases](https://github.com/bit101/isv/releases) and move to a location in your path.

#### Known issue on MacOS

Running the precompiled binary on MacOS may cause an error.

From [https://docs.fyne.io/faq/troubleshoot](https://docs.fyne.io/faq/troubleshoot):

> Q: Apple macOS says my app is damaged when it is downloaded
>
> A: When files are downloaded on a macOS computer they are marked with a “quarantine” flag so they are checked by the OS for problems. If your application is signed with a certificate purchased from Apple this is not a problem. However if you want to share your software without that cost this error may appear - and on M1/2 computers it is not possible to use the System Settings to allow the app to run.
>
> The fix is to remove the quarantine flag, which you can do by opening the Terminal and executing the following command:
>
> `sudo xattr -r -d com.apple.quarantine MyApp.app`
