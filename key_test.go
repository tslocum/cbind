package cbind

import (
	"testing"

	"github.com/gdamore/tcell"
)

type testCase struct {
	mod     tcell.ModMask
	key     tcell.Key
	ch      rune
	encoded string
}

var testCases = []testCase{
	{mod: tcell.ModNone, key: tcell.KeyRune, ch: 'a', encoded: "a"},
	{mod: tcell.ModNone, key: tcell.KeyRune, ch: '+', encoded: "+"},
	{mod: tcell.ModNone, key: tcell.KeyRune, ch: ';', encoded: ";"},
	{mod: tcell.ModNone, key: tcell.KeyEnter, ch: 0, encoded: "Enter"},
	{mod: tcell.ModAlt, key: tcell.KeyRune, ch: 'a', encoded: "Alt+a"},
	{mod: tcell.ModAlt, key: tcell.KeyRune, ch: '+', encoded: "Alt++"},
	{mod: tcell.ModAlt, key: tcell.KeyRune, ch: ';', encoded: "Alt+;"},
	{mod: tcell.ModAlt, key: tcell.KeyEnter, ch: 0, encoded: "Alt+Enter"},
	{mod: tcell.ModAlt, key: tcell.KeyRune, ch: ' ', encoded: "Alt+Space"},
	{mod: tcell.ModAlt, key: tcell.KeyBackspace2, ch: 0, encoded: "Alt+Backspace"},
	{mod: tcell.ModAlt, key: tcell.KeyPgDn, ch: 0, encoded: "Alt+PageDown"},
	{mod: tcell.ModCtrl | tcell.ModAlt, key: tcell.KeyRune, ch: '+', encoded: "Ctrl+Alt++"},
	{mod: tcell.ModCtrl | tcell.ModShift, key: tcell.KeyRune, ch: '+', encoded: "Ctrl+Shift++"},
}

func TestEncode(t *testing.T) {
	t.Parallel()

	for _, c := range testCases {
		encoded, err := Encode(c.mod, c.key, c.ch)
		if err != nil {
			t.Errorf("failed to encode key %d %d %d: %s", c.mod, c.key, c.ch, err)
		}
		if encoded != c.encoded {
			t.Errorf("failed to encode key %d %d %d: got %s, want %s", c.mod, c.key, c.ch, encoded, c.encoded)
		}
	}
}

func TestDecode(t *testing.T) {
	t.Parallel()

	for _, c := range testCases {
		mod, key, ch, err := Decode(c.encoded)
		if err != nil {
			t.Errorf("failed to decode key %s: %s", c.encoded, err)
		}
		if mod != c.mod {
			t.Errorf("failed to decode key %s: invalid modifiers: got %d, want %d", c.encoded, mod, c.mod)
		}
		if key != c.key {
			t.Errorf("failed to decode key %s: invalid key: got %d, want %d", c.encoded, key, c.key)
		}
		if ch != c.ch {
			t.Errorf("failed to decode key %s: invalid rune: got %d, want %d", c.encoded, ch, c.ch)
		}
	}
}