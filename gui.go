package main

import (
	"bufio"
	"fmt"
	"io"
	"log"

	"github.com/jroimartin/gocui"
)

// OutputDescription describes a command output.
type OutputDescription struct {
	Title string
	In    io.ReadCloser
}

// RunGUIOutput take a list of output description and print the command
// outputs inside the terminal.
func RunGUIOutput(outputs []OutputDescription) error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer g.Close()

	g.SetManagerFunc(func(g *gocui.Gui) error {
		return onUpdate(outputs, g)
	})

	// the SetKeybinding need to be set after SetManagerFunc
	g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	g.SetKeybinding("", gocui.KeyCtrlD, gocui.ModNone, quit)
	g.SetKeybinding("", 'q', gocui.ModNone, quit)

	for _, output := range outputs {
		go func(g *gocui.Gui, output OutputDescription) {

			reader := bufio.NewReader(output.In)
			for {
				line, prefix, err := reader.ReadLine()
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Println(err)
				}

				g.Update(func(g *gocui.Gui) error {
					view, err := g.View(output.Title)
					if err != nil {
						return err
					}

					if prefix {
						fmt.Fprint(view, string(line))
					} else {
						fmt.Fprintln(view, string(line))
					}

					return err
				})
			}
		}(g, output)
	}

	err = g.MainLoop()
	if err == gocui.ErrQuit {
		return nil
	}

	return err
}

func onUpdate(outputs []OutputDescription, g *gocui.Gui) error {
	maxX, maxY := g.Size()

	for idx, output := range outputs {
		view, err := g.SetView(output.Title, maxX/len(outputs)*idx, 0, maxX/len(outputs)*(idx+1), maxY)
		if err != nil && err != gocui.ErrUnknownView {
			return err
		}

		view.Title = " " + output.Title + " "
		view.Autoscroll = true
		view.Wrap = true
		view.Frame = true
	}

	return nil
}

func quit(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}
