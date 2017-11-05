package main

import (
	"github.com/gotk3/gotk3/gtk"
	"log"
	"github.com/gotk3/gotk3/cairo"
	"io/ioutil"
)

type Color struct {
	r, g, b int
}

var Scr []byte = make([]byte, 6912)

var Colors = map[int]Color{
	0x00: {0, 0, 0},
	0x08: {0, 0, 0},
	0x01: {0, 0, 0xc0},
	0x09: {0, 0, 0xff},
	0x02: {0xc0, 0, 0},
	0x0a: {0xff, 0, 0},
	0x03: {0xc0, 0, 0xc0},
	0x0b: {0xff, 0, 0xff},
	0x04: {0, 0xc0, 0},
	0x0c: {0, 0xff, 0},
	0x05: {0, 0xc0, 0xc0},
	0x0d: {0, 0xff, 0xff},
	0x06: {0xc0, 0xc0, 0},
	0x0e: {0xff, 0xff, 0},
	0x07: {0xc0, 0xc0, 0xc0},
	0x0f: {0xff, 0xff, 0xff},
}

func main() {
	loadScr()
	gtkInit()
}

func loadScr() {
	buf, err := ioutil.ReadFile("ne.scr")
	if err != nil {
		panic("File read error")
	}

	copy(Scr, buf)
}

func gtkInit() {
	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)
	// Create a new toplevel window, set its title, and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("Simple Example")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	canvas, err := gtk.DrawingAreaNew()
	if err != nil {
		log.Fatal("Unable to create canvas:", err)
	}
	canvas.Connect("draw", draw())
	// Add the label to the window.
	win.Add(canvas)
	// Set the default window size.
	win.SetDefaultSize(512, 384)
	// Recursively show all widgets contained in this window.
	win.ShowAll()
	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}

func draw() func(da *gtk.DrawingArea, cr *cairo.Context) {
	return func(da *gtk.DrawingArea, cr *cairo.Context) {
		for i := range Scr[:6144] {
			b := Scr[i]
			xBase := calcX(i)
			pattern := byte(1)
			y := calcY(i)
			for xOffset := 0; xOffset <= 7; xOffset++ {
				if (b & pattern) > 0 {
					red, green, blue := getColor(i)
					cr.SetSourceRGB(float64(red), float64(green), float64(blue))
				} else {
					red, green, blue := getBgColor(i)
					cr.SetSourceRGB(float64(red), float64(green), float64(blue))
				}
				cr.Rectangle(float64(xBase-xOffset+8)*2, float64(y)*2, 2.5, 2.5)
				cr.Fill()
				pattern = pattern << 1
			}
		}
	}
}

func calcX(addr int) int {
	return addr % 32 * 8
}

func calcY(addr int) int {
	correctedAddr := addr & 0x181F                        // 1100000011111
	correctedAddr = correctedAddr | ((addr & 0x700) >> 3) // 0011100000000
	correctedAddr = correctedAddr | ((addr & 0xE0) << 3)  // 0000011100000

	return correctedAddr / 32
}

func getColor(addr int) (int, int, int) {
	attrByte := getAttrByte(addr)
	ink := (attrByte & 0x7) + ((attrByte & 0x40) >> 3)
	color := Colors[int(ink)]

	return color.r, color.g, color.b
}

func getBgColor(addr int) (int, int, int) {
	attrByte := getAttrByte(addr)
	paper := ((attrByte & 0x38) >> 3) + ((attrByte & 0x40) >> 3)
	color := Colors[int(paper)]

	return color.r, color.g, color.b
}

func getAttrByte(addr int) byte {
	x := calcX(addr)
	y := calcY(addr)

	attrAddress := 0x1800 + y/8*32 + x/8

	return Scr[attrAddress]
}
