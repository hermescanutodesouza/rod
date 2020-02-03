package keys

import (
	"runtime"
	"time"
)

// Key contains information for generating a key press based off the unicode
// value.
//
// Example data for the following runes:
// 									'\r'  '\n'  | ','  '<'    | 'a'   'A'  | '\u0a07'
// 									_____________________________________________________
type Key struct {
	// Code is the key code:
	// 								"Enter"     | "Comma"     | "KeyA"     | "MediaStop"
	Code string

	// Key is the key value:
	// 								"Enter"     | ","   "<"   | "a"   "A"  | "MediaStop"
	Key string

	// Text is the text for printable keys:
	// 								"\r"  "\r"  | ","   "<"   | "a"   "A"  | ""
	Text string

	// Unmodified is the unmodified text for printable keys:
	// 								"\r"  "\r"  | ","   ","   | "a"   "a"  | ""
	Unmodified string

	// Native is the native scan code.
	// 								0x13  0x13  | 0xbc  0xbc  | 0x61  0x41 | 0x00ae
	Native int64

	// Windows is the windows scan code.
	// 								0x13  0x13  | 0xbc  0xbc  | 0x61  0x41 | 0xe024
	Windows int64

	// Shift indicates whether or not the Shift modifier should be sent.
	// 								false false | false true  | false true | false
	Shift bool

	// Print indicates whether or not the character is a printable character
	// (ie, should a "char" event be generated).
	// 								true  true  | true  true  | true  true | false
	Print bool
}

// KeyParams dispatches a key event to the page.
type KeyParams struct {
	Type                  string     `json:"type"`                            // Type of the key event.
	Modifiers             int64      `json:"modifiers"`                       // Bit field representing pressed modifier keys. Alt=1, Ctrl=2, Meta/Command=4, Shift=8 (default: 0).
	Timestamp             *time.Time `json:"timestamp,omitempty"`             // Time at which the event occurred.
	Text                  string     `json:"text,omitempty"`                  // Text as generated by processing a virtual key code with a keyboard layout. Not needed for for keyUp and rawKeyDown events (default: "")
	UnmodifiedText        string     `json:"unmodifiedText,omitempty"`        // Text that would have been generated by the keyboard if no modifiers were pressed (except for shift). Useful for shortcut (accelerator) key handling (default: "").
	KeyIdentifier         string     `json:"keyIdentifier,omitempty"`         // Unique key identifier (e.g., 'U+0041') (default: "").
	Code                  string     `json:"code,omitempty"`                  // Unique DOM defined string value for each physical key (e.g., 'KeyA') (default: "").
	Key                   string     `json:"key,omitempty"`                   // Unique DOM defined string value describing the meaning of the key in the context of active modifiers, keyboard layout, etc (e.g., 'AltGr') (default: "").
	WindowsVirtualKeyCode int64      `json:"windowsVirtualKeyCode,omitempty"` // Windows virtual key code (default: 0).
	NativeVirtualKeyCode  int64      `json:"nativeVirtualKeyCode,omitempty"`  // Native virtual key code (default: 0).
	AutoRepeat            bool       `json:"autoRepeat"`                      // Whether the event was generated from auto repeat (default: false).
	IsKeypad              bool       `json:"isKeypad"`                        // Whether the event was generated from the keypad (default: false).
	IsSystemKey           bool       `json:"isSystemKey"`                     // Whether the event was a system key event (default: false).
	Location              int64      `json:"location,omitempty"`              // Whether the event was from the left or right side of the keyboard. 1=Left, 2=Right (default: 0).
}

// Encode encodes a keyDown, char, and keyUp sequence for the specified rune.
func Encode(r rune) []*KeyParams {
	// force \n -> \r
	if r == '\n' {
		r = '\r'
	}

	// if not known key, encode as unidentified
	v := Keys[r]

	// create
	keyDown := KeyParams{
		Type:                  "keyDown",
		Key:                   v.Key,
		Code:                  v.Code,
		NativeVirtualKeyCode:  v.Native,
		WindowsVirtualKeyCode: v.Windows,
	}
	if runtime.GOOS == "darwin" {
		keyDown.NativeVirtualKeyCode = 0
	}
	if v.Shift {
		keyDown.Modifiers |= 8
	}

	keyUp := keyDown
	keyUp.Type = "keyUp"

	// printable, so create char event
	if v.Print {
		keyChar := keyDown
		keyChar.Type = "char"
		keyChar.Text = v.Text
		keyChar.UnmodifiedText = v.Unmodified

		// the virtual key code for char events for printable characters will
		// be different than the defined keycode when not shifted...
		//
		// specifically, it always sends the ascii value as the scan code,
		// which is available as the rune.
		keyChar.NativeVirtualKeyCode = int64(r)
		keyChar.WindowsVirtualKeyCode = int64(r)

		return []*KeyParams{&keyDown, &keyChar, &keyUp}
	}

	return []*KeyParams{&keyDown, &keyUp}
}