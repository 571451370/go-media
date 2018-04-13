package imgui

import (
	"bytes"
	"fmt"
	"math"
	"strings"

	"github.com/qeedquan/go-media/math/f64"
	"github.com/qeedquan/go-media/math/mathutil"
)

// Shared state of InputText(), passed to callback when a ImGuiInputTextFlags_Callback* flag is used and the corresponding callback is triggered.
type TextEditCallbackData struct {
	EventFlag InputTextFlags // One of ImGuiInputTextFlags_Callback* // Read-only
	Flags     InputTextFlags // What user passed to InputText()      // Read-only
	ReadOnly  bool           // Read-only mode                       // Read-only

	// CharFilter event:
	EventChar rune // Character input                      // Read-write (replace character or set to zero)

	// Completion,History,Always events:
	// If you modify the buffer contents make sure you update 'BufTextLen' and set 'BufDirty' to true.
	EventKey       Key           // Key pressed (Up/Down/TAB)            // Read-only
	Buf            *bytes.Buffer // Current text buffer                  // Read-write (pointed data only, can't replace the actual pointer)
	BufDirty       bool          // Set if you modify Buf/BufTextLen!!   // Write
	CursorPos      int           //                                      // Read-write
	SelectionStart int           //                                      // Read-write (== to SelectionEnd when no selection)
	SelectionEnd   int           //                                      // Read-write
}

type TextEditCallback func(*TextEditCallbackData)

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
	InputTextFlagsCharsScientific     InputTextFlags = 1 << 17 // Allow 0123456789.+-*/eE (Scientific notation input)
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
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}
	text := fmt.Sprintf(format, args...)
	c.TextUnformatted(text)
}

func (c *Context) TextUnformatted(text string) {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}

	text_pos := f64.Vec2{
		window.DC.CursorPos.X,
		window.DC.CursorPos.Y + window.DC.CurrentLineTextBaseOffset,
	}
	wrap_pos_x := window.DC.TextWrapPos
	wrap_enabled := wrap_pos_x >= 0.0

	if len(text) > 2000 && !wrap_enabled {
		// Long text!
		// Perform manual coarse clipping to optimize for long multi-line text
		// From this point we will only compute the width of lines that are visible. Optimization only available when word-wrapping is disabled.
		// We also don't vertically center the text within the line full height, which is unlikely to matter because we are likely the biggest and only item on the line.
		line := 0
		line_height := c.GetTextLineHeight()
		clip_rect := window.ClipRect
		text_size := f64.Vec2{0, 0}

		if text_pos.Y <= clip_rect.Max.Y {
			pos := text_pos

			// Lines to skip (can't skip when logging text)
			if !c.LogEnabled {
				lines_skippable := int((clip_rect.Min.Y - text_pos.Y) / line_height)
				if lines_skippable > 0 {
					lines_skipped := 0
					for line < len(text) && lines_skipped < lines_skippable {
						line_end := strings.IndexRune(text[line:], '\n')
						if line_end < 0 {
							line_end = len(text) - line - 1
						}
						line += line_end + 1
						lines_skipped++
					}
					pos.Y += float64(lines_skipped) * line_height
				}
			}

			// Lines to render
			if line < len(text) {
				line_rect := f64.Rectangle{pos, pos.Add(f64.Vec2{math.MaxFloat32, line_height})}
				for line < len(text) {
					line_end := strings.IndexRune(text[line:], '\n')
					if c.IsClippedEx(line_rect, 0, false) {
						break
					}

					line_size := c.CalcTextSizeEx(text[line:], false, -1)
					text_size.X = math.Max(text_size.X, line_size.X)
					c.RenderTextEx(pos, text[line:], false)
					if line_end < 0 {
						line_end = len(text) - line - 1
					}
					line += line_end + 1
					line_rect.Min.Y += line_height
					line_rect.Max.Y += line_height
					pos.Y += line_height
				}

				// Count remaining lines
				lines_skipped := 0
				for line < len(text) {
					line_end := strings.IndexRune(text[line:], '\n')
					if line_end < 0 {
						line_end = len(text) - line - 1
					}
					line = line_end + 1
					lines_skipped++
				}
				pos.Y += float64(lines_skipped) * line_height
			}

			text_size.Y += (pos.Sub(text_pos)).Y
		}

		bb := f64.Rectangle{text_pos, text_pos.Add(text_size)}
		c.ItemSizeBB(bb)
		c.ItemAdd(bb, 0)
	} else {
		wrap_width := 0.0
		if wrap_enabled {
			wrap_width = c.CalcWrapWidthForPos(window.DC.CursorPos, wrap_pos_x)
		}
		text_size := c.CalcTextSizeEx(text, false, wrap_width)

		// Account of baseline offset
		bb := f64.Rectangle{text_pos, text_pos.Add(text_size)}
		c.ItemSize(text_size)
		if !c.ItemAdd(bb, 0) {
			return
		}

		// Render (we don't hide text after ## in this end-user function)
		c.RenderTextWrapped(bb.Min, text, wrap_width)
	}
}

func (c *Context) RenderTextWrapped(pos f64.Vec2, text string, wrap_width float64) {
	window := c.CurrentWindow

	if len(text) > 0 {
		window.DrawList.AddTextEx(c.Font, c.FontSize, pos, c.GetColorFromStyle(ColText), text, wrap_width, nil)
		if c.LogEnabled {
			c.LogRenderedText(&pos, text)
		}
	}
}

// Horizontal separating line.
func (c *Context) Separator() {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}

	var flags SeparatorFlags
	if flags&(SeparatorFlagsHorizontal|SeparatorFlagsVertical) == 0 {
		if window.DC.LayoutType == LayoutTypeHorizontal {
			flags |= SeparatorFlagsVertical
		} else {
			flags |= SeparatorFlagsHorizontal
		}
	}
	// Check that only 1 option is selected
	assert(mathutil.IsPow2(int(flags & (SeparatorFlagsHorizontal | SeparatorFlagsVertical))))

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
		c.LogRenderedText(nil, "--------------------------------\n")
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

func (c *Context) InputTextExCallback(label, buf string, size_arg f64.Vec2, flags InputTextFlags, callback func()) bool {
	return false
}

func (c *Context) AlignTextToFramePadding() {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}

	window.DC.CurrentLineHeight = math.Max(window.DC.CurrentLineHeight, c.FontSize+c.Style.FramePadding.Y*2)
	window.DC.CurrentLineTextBaseOffset = math.Max(window.DC.CurrentLineTextBaseOffset, c.Style.FramePadding.Y)
}

func (c *Context) TextWrapped(format string, args ...interface{}) {
	// Keep existing wrap position is one ia already set
	need_wrap := c.CurrentWindow.DC.TextWrapPos < 0.0
	if need_wrap {
		c.PushTextWrapPos(0.0)
	}
	c.Text(format, args...)
	if need_wrap {
		c.PopTextWrapPos()
	}
}

func (c *Context) PushTextWrapPos(wrap_pos_x float64) {
	window := c.GetCurrentWindow()
	window.DC.TextWrapPos = wrap_pos_x
	window.DC.TextWrapPosStack = append(window.DC.TextWrapPosStack, wrap_pos_x)
}

func (c *Context) PopTextWrapPos() {
	window := c.GetCurrentWindow()
	window.DC.TextWrapPosStack = window.DC.TextWrapPosStack[:len(window.DC.TextWrapPosStack)-1]
	window.DC.TextWrapPos = -1.0
	if len(window.DC.TextWrapPosStack) > 0 {
		window.DC.TextWrapPos = window.DC.TextWrapPosStack[len(window.DC.TextWrapPosStack)-1]
	}
}

func (c *Context) InputText(label, buf string, flags InputTextFlags, callback TextEditCallback) bool {
	// call InputTextMultiline()
	assert(flags&InputTextFlagsMultiline == 0)
	return c.InputTextEx(label, buf, f64.Vec2{0, 0}, flags, callback)
}

func (c *Context) InputTextMultiline(label, buf string, size f64.Vec2, flags InputTextFlags, callback TextEditCallback) bool {
	return c.InputTextEx(label, buf, size, flags|InputTextFlagsMultiline, callback)
}

// Edit a string of text
// NB: when active, hold on a privately held copy of the text (and apply back to 'buf'). So changing 'buf' while active has no effect.
// FIXME: Rather messy function partly because we are doing UTF8 > u16 > UTF8 conversions on the go to more easily handle stb_textedit calls. Ideally we should stay in UTF-8 all the time. See https://github.com/nothings/stb/issues/188
func (c *Context) InputTextEx(label, buf string, size_arg f64.Vec2, flags InputTextFlags, callback TextEditCallback) bool {
	return false
}