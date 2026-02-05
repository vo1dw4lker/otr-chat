package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
	"time"
)

func RunChatUI(rcvd, send chan string) error {
	app := tview.NewApplication()

	pages := tview.NewPages()
	msgs := newMessages()

	chatPage, output := makeChatPage(app, pages, send, rcvd)
	pages.AddAndSwitchToPage("chat", chatPage, true)

	settingsPage := makeSettingsPage(pages, msgs)
	pages.AddPage("settings", settingsPage, true, false)

	app = app.SetRoot(pages, true)

	go listenForMessages(app, rcvd, msgs)
	go redraw(output, msgs)
	rcvd <- "You can switch to settings page by pressing Tab\n\n"
	if err := app.Run(); err != nil {
		return err
	}
	return nil
}

func makeChatPage(app *tview.Application, pages *tview.Pages, send, rcvd chan string) (layout *tview.Flex, output *tview.TextView) {
	output = tview.NewTextView().
		SetChangedFunc(func() {
			app.Draw()
		}).SetScrollable(false)
	output.SetBorder(true).SetTitle("Chat")

	input := tview.NewTextArea()
	input.SetBorder(true)
	input.SetChangedFunc(inputChanged(input, send, rcvd))
	input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			pages.SwitchToPage("settings")
			return nil
		}
		return event
	})
	layout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(output, 0, 4, false).
		AddItem(input, 0, 1, true)
	return layout, output
}

func makeSettingsPage(pages *tview.Pages, msgs *messages) *tview.Flex {
	dropdown := tview.NewDropDown().
		SetLabel("Clear messages after:")
	for _, option := range timeOptions {
		opt := option
		dropdown.AddOption(opt.Label, func() {
			msgs.SetTTL(opt.Duration)
		})
	}
	dropdown.SetCurrentOption(0)
	dropdown.SetBorder(true)
	dropdown.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			pages.SwitchToPage("chat")
		}
		return event
	})

	horizontal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(dropdown, 0, 2, true).
		AddItem(nil, 0, 1, false)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(horizontal, 3, 0, true).
		AddItem(nil, 0, 1, false)

	layout.SetBorder(true).SetTitle("Settings")

	return layout
}

func inputChanged(input *tview.TextArea, send, rcvd chan string) func() {
	return func() {
		text := input.GetText()
		if !strings.Contains(text, "\n") {
			return
		}

		send <- text

		rcvd <- fmt.Sprintf("%s: %s", "you", text)
		input.SetText("", true)
	}
}

func listenForMessages(app *tview.Application, src chan string, msgs *messages) {
	for buf := range src {
		msgs.AppendMessage(buf)
	}
	app.Stop()
}

func redraw(tv *tview.TextView, msgs *messages) {
	for {
		time.Sleep(50 * time.Millisecond)
		tv.SetText(msgs.String())
	}
}
