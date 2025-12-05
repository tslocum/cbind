package cbind

import (
	"fmt"
	"testing"

	"github.com/gdamore/tcell/v3"
)

type testCase struct {
	mod     tcell.ModMask
	key     tcell.Key
	str     string
	encoded string
}

func (c testCase) String() string {
	var str string
	if c.str != "" {
		str = "-" + str
	}
	return fmt.Sprintf("%d-%d%s-%s", c.mod, c.key, str, c.encoded)
}

var testCases = []testCase{
	{mod: tcell.ModNone, key: tcell.KeyRune, str: "a", encoded: "a"},
	{mod: tcell.ModNone, key: tcell.KeyRune, str: "+", encoded: "+"},
	{mod: tcell.ModNone, key: tcell.KeyRune, str: ";", encoded: ";"},
	{mod: tcell.ModNone, key: tcell.KeyTab, str: string(rune(tcell.KeyTab)), encoded: "Tab"},
	{mod: tcell.ModNone, key: tcell.KeyEnter, str: string(rune(tcell.KeyEnter)), encoded: "Enter"},
	{mod: tcell.ModNone, key: tcell.KeyPgDn, str: "", encoded: "PageDown"},
	{mod: tcell.ModAlt, key: tcell.KeyRune, str: "a", encoded: "Alt+a"},
	{mod: tcell.ModAlt, key: tcell.KeyRune, str: "+", encoded: "Alt++"},
	{mod: tcell.ModAlt, key: tcell.KeyRune, str: ";", encoded: "Alt+;"},
	{mod: tcell.ModAlt, key: tcell.KeyRune, str: " ", encoded: "Alt+Space"},
	{mod: tcell.ModAlt, key: tcell.KeyRune, str: "1", encoded: "Alt+1"},
	{mod: tcell.ModAlt, key: tcell.KeyTab, str: string(rune(tcell.KeyTab)), encoded: "Alt+Tab"},
	{mod: tcell.ModAlt, key: tcell.KeyEnter, str: string(rune(tcell.KeyEnter)), encoded: "Alt+Enter"},
	{mod: tcell.ModAlt, key: tcell.KeyDelete, str: "", encoded: "Alt+Delete"},
	{mod: tcell.ModCtrl, key: tcell.KeyRune, str: "c", encoded: "Ctrl+c"},
	{mod: tcell.ModCtrl, key: tcell.KeyRune, str: "d", encoded: "Ctrl+d"},
	{mod: tcell.ModCtrl | tcell.ModAlt, key: tcell.KeyRune, str: "c", encoded: "Ctrl+Alt+c"},
	{mod: tcell.ModCtrl, key: tcell.KeyRune, str: " ", encoded: "Ctrl+Space"},
	{mod: tcell.ModCtrl | tcell.ModAlt, key: tcell.KeyRune, str: "+", encoded: "Ctrl+Alt++"},
	{mod: tcell.ModCtrl | tcell.ModShift, key: tcell.KeyRune, str: "+", encoded: "Ctrl+Shift++"},
}

func TestEncode(t *testing.T) {
	t.Parallel()

	for _, c := range testCases {
		encoded, err := Encode(c.mod, c.key, c.str)
		if err != nil {
			t.Errorf("failed to encode key %d %d %s: %s", c.mod, c.key, c.str, err)
		}
		if encoded != c.encoded {
			t.Errorf("failed to encode key %d %d %s: got %s, want %s", c.mod, c.key, c.str, encoded, c.encoded)
		}
	}
}

func TestDecode(t *testing.T) {
	t.Parallel()

	for _, c := range testCases {
		mod, key, str, err := Decode(c.encoded)
		if err != nil {
			t.Errorf("failed to decode key %s: %s", c.encoded, err)
		}
		if mod != c.mod {
			t.Errorf("failed to decode key %s: invalid modifiers: got %d, want %d", c.encoded, mod, c.mod)
		}
		if key != c.key {
			t.Errorf("failed to decode key %s: invalid key: got %d, want %d", c.encoded, key, c.key)
		}
		if str != c.str {
			t.Errorf("failed to decode key %s: invalid rune: got %s, want %s", c.encoded, str, c.str)
		}
	}
}
