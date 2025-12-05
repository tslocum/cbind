package cbind

import (
	"fmt"
	"sync"

	"github.com/gdamore/tcell/v3"
)

type eventHandler func(ev *tcell.EventKey) *tcell.EventKey

// Configuration maps keys to event handlers and processes key events.
type Configuration struct {
	handlers map[string]eventHandler
	mutex    *sync.RWMutex
}

// NewConfiguration returns a new input configuration.
func NewConfiguration() *Configuration {
	c := Configuration{
		handlers: make(map[string]eventHandler),
		mutex:    new(sync.RWMutex),
	}

	return &c
}

// Set sets the handler for a key event string.
func (c *Configuration) Set(s string, handler func(ev *tcell.EventKey) *tcell.EventKey) error {
	mod, key, str, err := Decode(s)
	if err != nil {
		return err
	}

	if key == tcell.KeyRune {
		c.SetRune(mod, []rune(str)[0], handler)
	} else {
		c.SetKey(mod, key, handler)
	}
	return nil
}

// SetKey sets the handler for a key.
func (c *Configuration) SetKey(mod tcell.ModMask, key tcell.Key, handler func(ev *tcell.EventKey) *tcell.EventKey) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Convert KeyCtrlA-Z to rune format.
	if key >= tcell.KeyCtrlA && key <= tcell.KeyCtrlZ {
		mod |= tcell.ModCtrl
		r := 'a' + (key - tcell.KeyCtrlA)
		c.handlers[fmt.Sprintf("%d:%d", mod, r)] = handler
		return
	}

	// Convert Shift+Tab to Backtab.
	if mod&tcell.ModShift != 0 && key == tcell.KeyTab {
		mod ^= tcell.ModShift
		key = tcell.KeyBacktab
	}

	c.handlers[fmt.Sprintf("%d-%d", mod, key)] = handler
}

// SetRune sets the handler for a rune.
func (c *Configuration) SetRune(mod tcell.ModMask, ch rune, handler func(ev *tcell.EventKey) *tcell.EventKey) {
	// Some runes are identical to named keys. Set the bind on the matching
	// named key instead.
	switch ch {
	case '\t':
		c.SetKey(mod, tcell.KeyTab, handler)
		return
	case '\n':
		c.SetKey(mod, tcell.KeyEnter, handler)
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.handlers[fmt.Sprintf("%d:%d", mod, ch)] = handler
}

// Capture handles key events.
func (c *Configuration) Capture(ev *tcell.EventKey) *tcell.EventKey {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if ev == nil {
		return nil
	}

	mod := ev.Modifiers()
	key := ev.Key()
	str := ev.Str()

	// Convert KeyCtrlA-Z to rune format.
	if key >= tcell.KeyCtrlA && key <= tcell.KeyCtrlZ {
		mod |= tcell.ModCtrl
		str = string(rune('a' + (key - tcell.KeyCtrlA)))
		key = tcell.KeyRune
	}

	var keyName string
	if key != tcell.KeyRune {
		keyName = fmt.Sprintf("%d-%d", mod, key)
	} else {
		keyName = fmt.Sprintf("%d:%d", mod, []rune(str)[0])
	}

	handler := c.handlers[keyName]
	if handler != nil {
		return handler(ev)
	}
	return ev
}

// Clear removes all handlers.
func (c *Configuration) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.handlers = make(map[string]eventHandler)
}
