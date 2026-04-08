package cbind

import (
	"errors"
	"strings"

	"github.com/gdamore/tcell/v3"
)

// Modifier labels
const (
	LabelCtrl  = "ctrl"
	LabelAlt   = "alt"
	LabelMeta  = "meta"
	LabelShift = "shift"
)

// ErrInvalidKeyEvent is the error returned when encoding or decoding a key event fails.
var ErrInvalidKeyEvent = errors.New("invalid key event")

// UnifyEnterKeys is a flag that determines whether or not KPEnter (keypad
// enter) key events are interpreted as Enter key events. When enabled, Ctrl+J
// key events are also interpreted as Enter key events.
var UnifyEnterKeys = true

var fullKeyNames = map[string]string{
	"pgup": "PageUp",
	"pgdn": "PageDown",
	"esc":  "Escape",
}

// Decode decodes a string as a key or combination of keys.
func Decode(s string) (mod tcell.ModMask, key tcell.Key, str string, err error) {
	if len(s) == 0 {
		return 0, 0, "", ErrInvalidKeyEvent
	}

	// Special case for plus rune decoding
	if s[len(s)-1:] == "+" {
		key = tcell.KeyRune
		str = "+"

		if len(s) == 1 {
			return mod, key, str, nil
		} else if len(s) == 2 {
			return 0, 0, "", ErrInvalidKeyEvent
		} else {
			s = s[:len(s)-2]
		}
	}

	split := strings.Split(s, "+")
DECODEPIECE:
	for _, piece := range split {
		// Decode modifiers
		pieceLower := strings.ToLower(piece)
		switch pieceLower {
		case LabelCtrl:
			mod |= tcell.ModCtrl
			continue
		case LabelAlt:
			mod |= tcell.ModAlt
			continue
		case LabelMeta:
			mod |= tcell.ModMeta
			continue
		case LabelShift:
			mod |= tcell.ModShift
			continue
		}

		// Decode key
		for shortKey, fullKey := range fullKeyNames {
			if pieceLower == strings.ToLower(fullKey) {
				pieceLower = shortKey
				break
			}
		}
		switch pieceLower {
		case "space", "spacebar":
			key = tcell.KeyRune
			str = " "
			continue
		}
		for k, keyName := range tcell.KeyNames {
			if pieceLower == strings.ToLower(strings.ReplaceAll(keyName, "-", "+")) {
				key = k
				if key < 0x80 {
					str = string(rune(k))
				}
				continue DECODEPIECE
			}
		}

		// Decode rune
		if len(piece) > 1 {
			return 0, 0, "", ErrInvalidKeyEvent
		}

		key = tcell.KeyRune
		str = string(rune(piece[0]))
	}

	// Normalize Ctrl+A-Z to lowercase
	if mod&tcell.ModCtrl != 0 && key == tcell.KeyRune {
		str = strings.ToLower(str)
	}

	return mod, key, str, nil
}

// Encode encodes a key or combination of keys a string.
func Encode(mod tcell.ModMask, key tcell.Key, str string) (string, error) {
	var b strings.Builder
	var wrote bool

	if mod&tcell.ModCtrl != 0 {
		if key == tcell.KeyBackspace || key == tcell.KeyTab || key == tcell.KeyEnter {
			mod ^= tcell.ModCtrl
		} else {
			// Convert KeyCtrlA-Z to rune format.
			if key >= tcell.KeyCtrlA && key <= tcell.KeyCtrlZ {
				mod |= tcell.ModCtrl
				str = string(rune('a' + (key - tcell.KeyCtrlA)))
				key = tcell.KeyRune
			}
		}
	}

	if key != tcell.KeyRune {
		if UnifyEnterKeys && key == tcell.KeyCtrlJ {
			key = tcell.KeyEnter
		} else if key < 0x80 {
			str = string(rune(key))
		}
	}

	// Encode modifiers
	if mod&tcell.ModCtrl != 0 {
		b.WriteString(upperFirst(LabelCtrl))
		wrote = true
	}
	if mod&tcell.ModAlt != 0 {
		if wrote {
			b.WriteRune('+')
		}
		b.WriteString(upperFirst(LabelAlt))
		wrote = true
	}
	if mod&tcell.ModMeta != 0 {
		if wrote {
			b.WriteRune('+')
		}
		b.WriteString(upperFirst(LabelMeta))
		wrote = true
	}
	if mod&tcell.ModShift != 0 {
		if wrote {
			b.WriteRune('+')
		}
		b.WriteString(upperFirst(LabelShift))
		wrote = true
	}

	if key == tcell.KeyRune && str == " " {
		if wrote {
			b.WriteRune('+')
		}
		b.WriteString("Space")
	} else if key != tcell.KeyRune {
		// Encode key
		keyName := tcell.KeyNames[key]
		if keyName == "" {
			return "", ErrInvalidKeyEvent
		}
		keyName = strings.ReplaceAll(keyName, "-", "+")
		fullKeyName := fullKeyNames[strings.ToLower(keyName)]
		if fullKeyName != "" {
			keyName = fullKeyName
		}

		if wrote {
			b.WriteRune('+')
		}
		b.WriteString(keyName)
	} else {
		// Encode rune
		if wrote {
			b.WriteRune('+')
		}
		b.WriteString(str)
	}

	return b.String(), nil
}

func upperFirst(s string) string {
	if len(s) <= 1 {
		return strings.ToUpper(s)
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
