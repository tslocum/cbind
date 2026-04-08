package main

import (
	"fmt"
	"os"
	"strings"

	"codeberg.org/tslocum/cbind"
	"github.com/gdamore/tcell/v3"
)

func printInfo(mod tcell.ModMask, key tcell.Key, str string) string {
	modLabel := "ModNone"
	if mod != 0 {
		var m []string
		if mod&tcell.ModShift != 0 {
			m = append(m, "ModShift")
		}
		if mod&tcell.ModAlt != 0 {
			m = append(m, "ModAlt")
		}
		if mod&tcell.ModMeta != 0 {
			m = append(m, "ModMeta")
		}
		if mod&tcell.ModCtrl != 0 {
			m = append(m, "ModCtrl")
		}
		if mod&tcell.ModHyper != 0 {
			m = append(m, "ModHyper")
		}
		modLabel = strings.Join(m, "+")
	}
	var keyLabel string
	keyName := tcell.KeyNames[key]
	if keyName != "" {
		keyLabel = "Key" + keyName
	} else {
		switch key {
		case tcell.KeyRune:
			keyLabel = "KeyRune"
		default:
			keyLabel = "Unknown"
		}
	}
	var strLabel string
	if str != "" {
		strLabel = ", str '" + str + "'"
	}
	return fmt.Sprintf("mod %d (%s), key %d (%s)%s", mod, modLabel, key, keyLabel, strLabel)
}

func main() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	quit := make(chan struct{})

	quitApp := func(ev *tcell.EventKey) *tcell.EventKey {
		quit <- struct{}{}
		return nil
	}

	configuration := cbind.NewConfiguration()
	configuration.SetKey(tcell.ModNone, tcell.KeyEscape, quitApp)
	configuration.SetRune(tcell.ModCtrl, 'c', quitApp)

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))
	s.Clear()

	putln(s, 0, "No key press events have been detected yet.")
	putln(s, 2, "Press a key.")
	putln(s, 4, "Key event info will then be displayed here.")

	go func() {
		for ev := range s.EventQ() {
			switch ev := ev.(type) {
			case *tcell.EventResize:
				s.Sync()
			case *tcell.EventKey:
				s.SetStyle(tcell.StyleDefault.
					Foreground(tcell.ColorWhite).
					Background(tcell.ColorBlack))
				s.Clear()

				putln(s, 0, fmt.Sprintf("Decoded as %s", printInfo(ev.Modifiers(), ev.Key(), ev.Str())))

				str, err := cbind.Encode(ev.Modifiers(), ev.Key(), ev.Str())
				if err != nil {
					str = fmt.Sprintf("error: %s", err)
				}
				putln(s, 2, "Labeled as '"+str+"'")

				mod, key, str, err := cbind.Decode(str)
				if err != nil {
					putln(s, 4, err.Error())
				} else {
					putln(s, 4, fmt.Sprintf("Encoded as %s", printInfo(mod, key, str)))
				}

				configuration.Capture(ev)

				s.Sync()
			}
		}
	}()
	s.Show()

	<-quit
	s.Fini()
}

// putln and puts functions are copied from the tcell unicode demo.
// Apache License, Version 2.0

func putln(s tcell.Screen, y int, str string) {
	puts(s, tcell.StyleDefault, 0, y, str)
}

func puts(s tcell.Screen, style tcell.Style, x, y int, str string) {
	i := 0
	var deferred []rune
	dwidth := 0
	zwj := false
	for _, r := range str {
		if r == '\u200d' {
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
			deferred = append(deferred, r)
			zwj = true
			continue
		}
		if zwj {
			deferred = append(deferred, r)
			zwj = false
			continue
		}
		if len(deferred) != 0 {
			s.SetContent(x+i, y, deferred[0], deferred[1:], style)
			i += dwidth
		}
		deferred = nil
		dwidth = 1
		deferred = append(deferred, r)
	}
	if len(deferred) != 0 {
		s.SetContent(x+i, y, deferred[0], deferred[1:], style)
		i += dwidth
	}
}
