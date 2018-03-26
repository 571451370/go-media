package imgui

import (
	"math"
	"strings"

	"github.com/qeedquan/go-media/math/f64"
)

type SeparatorFlags int

const (
	SeparatorFlagsHorizontal SeparatorFlags = 1 << 0 // Axis default to current layout type, so generally Horizontal unless e.g. in a menu bar
	SeparatorFlagsVertical   SeparatorFlags = 1 << 1
)

type InputTextFlags int

const (
	InputTextFlagsCharsDecimal        InputTextFlags = 1 << 0  // Allow 0123456789.+-*/
	InputTextFlagsCharsHexadecimal    InputTextFlags = 1 << 1  // Allow 0123456789ABCDEFabcdef
	InputTextFlagsCharsUppercase      InputTextFlags = 1 << 2  // Turn a..z into A..Z
	InputTextFlagsCharsNoBlank        InputTextFlags = 1 << 3  // Filter out spaces tabs
	InputTextFlagsAutoSelectAll       InputTextFlags = 1 << 4  // Select entire text when first taking mouse focus
	InputTextFlagsEnterReturnsTrue    InputTextFlags = 1 << 5  // Return 'true' when Enter is pressed (as opposed to when the value was modified)
	InputTextFlagsCallbackCompletion  InputTextFlags = 1 << 6  // Call user function on pressing TAB (for completion handling)
	InputTextFlagsCallbackHistory     InputTextFlags = 1 << 7  // Call user function on pressing Up/Down arrows (for history handling)
	InputTextFlagsCallbackAlways      InputTextFlags = 1 << 8  // Call user function every time. User code may query cursor position modify text buffer.
	InputTextFlagsCallbackCharFilter  InputTextFlags = 1 << 9  // Call user function to filter character. Modify data->EventChar to replace/filter input or return 1 to discard character.
	InputTextFlagsAllowTabInput       InputTextFlags = 1 << 10 // Pressing TAB input a '\t' character into the text field
	InputTextFlagsCtrlEnterForNewLine InputTextFlags = 1 << 11 // In multi-line mode unfocus with Enter add new line with Ctrl+Enter (default is opposite: unfocus with Ctrl+Enter add line with Enter).
	InputTextFlagsNoHorizontalScroll  InputTextFlags = 1 << 12 // Disable following the cursor horizontally
	InputTextFlagsAlwaysInsertMode    InputTextFlags = 1 << 13 // Insert mode
	InputTextFlagsReadOnly            InputTextFlags = 1 << 14 // Read-only mode
	InputTextFlagsPassword            InputTextFlags = 1 << 15 // Password mode display all characters as '*'
	InputTextFlagsNoUndoRedo          InputTextFlags = 1 << 16 // Disable undo/redo. Note that input text owns the text data while active if you want to provide your own undo/redo stack you need e.g. to call ClearActiveID().
	// [Internal]
	InputTextFlagsMultiline InputTextFlags = 1 << 20 // For internal use by InputTextMultiline()
)

type TextEditState struct {
	Id                   ID     // widget id owning the text state
	Text                 []rune // edit buffer, we need to persist but can't guarantee the persistence of the user-provided buffer. so we copy into own buffer.
	InitialText          []byte // backup of end-user buffer at the time of focus (in UTF-8, unaltered)
	TempTextBuffer       []byte
	CurLenA, CurLenW     int // we need to maintain our buffer length in both UTF-8 and wchar format.
	BufSizeA             int // end-user buffer size
	ScrollX              float64
	StbState             int
	CursorAnim           float64
	CursorFollow         bool
	SelectedAllMouseLock bool
}

func (c *Context) GetTextLineHeight() float64 {
	return c.FontSize
}

func (c *Context) GetTextLineHeightWithSpacing() float64 {
	return c.FontSize + c.Style.ItemSpacing.Y
}

func (c *Context) CalcTextSize(text string) f64.Vec2 {
	return c.CalcTextSizeEx(text, false, -1)
}

// Calculate text size. Text can be multi-line. Optionally ignore text after a ## marker.
// CalcTextSize("") should return ImVec2(0.0f, GImGui->FontSize)
func (c *Context) CalcTextSizeEx(text string, hide_text_after_double_hash bool, wrap_width float64) f64.Vec2 {
	text_display_end := len(text)
	if hide_text_after_double_hash {
		// Hide anything after a '##' string
		text_display_end = c.FindRenderedTextEnd(text)
	}

	font := c.Font
	font_size := c.FontSize
	if text_display_end == 0 {
		return f64.Vec2{0, font_size}
	}
	text_size, _ := font.CalcTextSizeA(font_size, math.MaxFloat32, wrap_width, text[:text_display_end])

	// Cancel out character spacing for the last character of a line (it is baked into glyph->AdvanceX field)
	font_scale := font_size / font.FontSize
	character_spacing_x := 1.0 * font_scale
	if text_size.X > 0.0 {
		text_size.X -= character_spacing_x
	}
	text_size.X = float64(int(text_size.X + 0.95))

	return text_size
}

func (c *Context) FindRenderedTextEnd(text string) int {
	text_display_end := strings.Index(text, "##")
	if text_display_end == -1 {
		text_display_end = len(text)
	}
	return text_display_end
}

func (c *Context) Text(format string, args ...interface{}) {
}

func (c *Context) TextUnformatted(text string) {
}

// Horizontal separating line.
func (c *Context) Separator() {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}

	var flags SeparatorFlags
	if flags&(SeparatorFlagsHorizontal|SeparatorFlagsVertical) == 0 {
		flags |= SeparatorFlagsVertical
		if window.DC.LayoutType == LayoutTypeHorizontal {
			flags |= SeparatorFlagsVertical
		} else {
			flags |= SeparatorFlagsHorizontal
		}
	}

	if flags&SeparatorFlagsVertical != 0 {
		c.VerticalSeparator()
		return
	}

	// Horizontal Separator
	if window.DC.ColumnsSet != nil {
		c.PopClipRect()
	}

	x1 := window.Pos.X
	x2 := window.Pos.X + window.Size.X
	if len(window.DC.GroupStack) > 0 {
		x1 += window.DC.IndentX
	}

	bb := f64.Rectangle{
		f64.Vec2{x1, window.DC.CursorPos.Y},
		f64.Vec2{x2, window.DC.CursorPos.Y + 1.0},
	}

	// NB: we don't provide our width so that it doesn't get feed back into AutoFit, we don't provide height to not alter layout
	c.ItemSize(f64.Vec2{0, 0})
	if !c.ItemAdd(bb, 0) {
		if window.DC.ColumnsSet != nil {
			c.PushColumnClipRect()
		}
		return
	}

	window.DrawList.AddLine(bb.Min, f64.Vec2{bb.Max.X, bb.Min.Y}, c.GetColorFromStyle(ColSeparator))
	if c.LogEnabled {
		c.LogRenderedText("--------------------------------\n")
	}

	if window.DC.ColumnsSet != nil {
		c.PushColumnClipRect()
		window.DC.ColumnsSet.LineMinY = window.DC.CursorPos.Y
	}
}

func (c *Context) VerticalSeparator() {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}

	y1 := window.DC.CursorPos.Y
	y2 := window.DC.CursorPos.Y + window.DC.CurrentLineHeight
	bb := f64.Rectangle{
		f64.Vec2{window.DC.CursorPos.X, y1},
		f64.Vec2{window.DC.CursorPos.X + 1.0, y2},
	}
	c.ItemSize(f64.Vec2{bb.Dx(), 0})
	if !c.ItemAdd(bb, 0) {
		return
	}

	window.DrawList.AddLine(f64.Vec2{bb.Min.X, bb.Min.Y}, f64.Vec2{bb.Min.X, bb.Max.Y}, c.GetColorFromStyle(ColSeparator))
	if c.LogEnabled {
		c.LogText(" |")
	}
}

func (c *Context) GuiLayoutTypeHorizontal() {
}

func (c *Context) LogText(text string) {
}

func (c *Context) LogRenderedText(text string) {
}
