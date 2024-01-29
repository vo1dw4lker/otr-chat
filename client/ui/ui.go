package ui

import (
	"fmt"
	"github.com/rivo/tview"
	"io"
	"log"
	"strings"
)

func RunUI(rcvd, send chan string) error {
	app := tview.NewApplication()

	output := tview.NewTextView().
		SetChangedFunc(func() {
			app.Draw()
		}).SetScrollable(false)
	output.SetBorder(true).SetTitle("Chat")

	input := tview.NewTextArea()
	input.SetBorder(true)
	input.SetChangedFunc(inputChanged(output, input, send))

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(output, 0, 4, false).
		AddItem(input, 0, 1, true)
	app = app.SetRoot(layout, true)

	go listenForMessages(rcvd, output)
	if err := app.Run(); err != nil {
		return err
	}
	return nil
}

func inputChanged(writer io.Writer, input *tview.TextArea, send chan string) func() {
	return func() {
		text := input.GetText()
		if !strings.Contains(text, "\n") {
			return
		}

		send <- text

		// TODO: workaround, remove later (probably send to rcvd channel)
		_, _ = fmt.Fprintf(writer, "%s: %s", "user", text)
		input.SetText("", true)
	}
}

func listenForMessages(src chan string, writer io.Writer) {
	for {
		buf := <-src
		_, err := writer.Write([]byte(buf))
		if err != nil {
			log.Fatalln(err)
		}
	}
}
