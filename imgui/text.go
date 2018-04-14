package imgui

import (
	"bytes"
	"fmt"
	"math"
	"strings"

	"github.com/qeedquan/go-media/math/f64"
	"github.com/qeedquan/go-media/math/mathutil"
	"github.com/qeedquan/go-media/stb/stbte"
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
	Ctx                  *Context
	Id                   ID     // widget id owning the text state
	Text                 []rune // edit buffer, we need to persist but can't guarantee the persistence of the user-provided buffer. so we copy into own buffer.
	InitialText          []byte // backup of end-user buffer at the time of focus (in UTF-8, unaltered)
	TempTextBuffer       []byte
	CurLenA, CurLenW     int // we need to maintain our buffer length in both UTF-8 and wchar format.
	BufSizeA             int // end-user buffer size
	ScrollX              float64
	StbState             stbte.State
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
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	style := &c.Style
	io := &c.IO

	// Can't use both together (they both use up/down keys)
	assert(!(flags&InputTextFlagsCallbackHistory != 0 && flags&InputTextFlagsMultiline != 0))
	// Can't use both together (they both use tab key)
	assert(!(flags&InputTextFlagsCallbackCompletion != 0 && flags&InputTextFlagsAllowTabInput != 0))

	is_multiline := flags&InputTextFlagsMultiline != 0
	is_editable := flags&InputTextFlagsReadOnly == 0
	is_password := flags&InputTextFlagsPassword != 0
	is_undoable := flags&InputTextFlagsNoUndoRedo == 0

	// Open group before calling GetID() because groups tracks id created during their spawn
	if is_multiline {
		c.BeginGroup()
	}
	id := window.GetID(label)
	label_size := c.CalcTextSizeEx(label, true, -1)

	// Arbitrary default of 8 lines high for multi-line
	item_size := label_size.Y
	if is_multiline {
		item_size = c.GetTextLineHeight() * 8.0
	}
	size := c.CalcItemSize(size_arg, c.CalcItemWidth(), item_size+style.FramePadding.Y*2.0)
	frame_bb := f64.Rectangle{window.DC.CursorPos, window.DC.CursorPos.Add(size)}
	total_bb_x := 0.0
	if label_size.X > 0 {
		total_bb_x = style.ItemInnerSpacing.X + label_size.X
	}
	total_bb := f64.Rectangle{frame_bb.Min, frame_bb.Max.Add(f64.Vec2{total_bb_x, 0.0})}

	draw_window := window
	if is_multiline {
		c.ItemAddEx(total_bb, id, &frame_bb)
		if !c.BeginChildFrame(id, frame_bb.Size(), 0) {
			c.EndChildFrame()
			c.EndGroup()
			return false
		}
		draw_window = c.GetCurrentWindow()
		size.X -= draw_window.ScrollbarSizes.X
	} else {
		c.ItemSizeBBEx(total_bb, style.FramePadding.Y)
		if !c.ItemAddEx(total_bb, id, &frame_bb) {
			return false
		}
	}

	hovered := c.ItemHoverable(frame_bb, id)
	if hovered {
		c.MouseCursor = MouseCursorTextInput
	}

	// Password pushes a temporary font with only a fallback glyph
	if is_password {
		glyph := c.Font.FindGlyph('*')
		password_font := &c.InputTextPasswordFont
		password_font.FontSize = c.Font.FontSize
		password_font.Scale = c.Font.Scale
		password_font.DisplayOffset = c.Font.DisplayOffset
		password_font.Ascent = c.Font.Ascent
		password_font.Descent = c.Font.Descent
		password_font.ContainerAtlas = c.Font.ContainerAtlas
		password_font.FallbackGlyph = glyph
		password_font.FallbackAdvanceX = glyph.AdvanceX
		assert(len(password_font.Glyphs) == 0 && len(password_font.IndexAdvanceX) == 0 && len(password_font.IndexLookup) == 0)
		c.PushFont(password_font)
	}

	// NB: we are only allowed to access 'edit_state' if we are the active widget.
	edit_state := &c.InputTextState

	// Using completion callback disable keyboard tabbing
	focus_requested := c.FocusableItemRegisterEx(window, id, flags&(InputTextFlagsCallbackCompletion|InputTextFlagsAllowTabInput) != 0)
	focus_requested_by_code := focus_requested && (window.FocusIdxAllCounter == window.FocusIdxAllRequestCurrent)
	focus_requested_by_tab := focus_requested && !focus_requested_by_code

	user_clicked := hovered && io.MouseClicked[0]
	user_scrolled := is_multiline && c.ActiveId == 0 && edit_state.Id == id && c.ActiveIdPreviousFrame == draw_window.GetIDNoKeepAlive("#SCROLLY")
	user_nav_input_start := (c.ActiveId != id) && ((c.NavInputId == id) || (c.NavActivateId == id && c.NavInputSource == InputSourceNavKeyboard))

	clear_active_id := false

	select_all := (c.ActiveId != id) && ((flags&InputTextFlagsAutoSelectAll) != 0 || user_nav_input_start) && (!is_multiline)
	if focus_requested || user_clicked || user_scrolled || user_nav_input_start {
		if c.ActiveId != id {
			// Start edition
			// Take a copy of the initial buffer value (both in original UTF-8 format and converted to wchar)
			// From the moment we focused we are ignoring the content of 'buf' (unless we are in read-only mode)
			prev_len_w := edit_state.CurLenW
			edit_state.Text = []rune(buf)
			edit_state.InitialText = []byte(buf)
			edit_state.CurLenW = len(edit_state.Text)
			edit_state.CurLenA = len(buf)
			edit_state.CursorAnimReset()

			// Preserve cursor position and undo/redo stack if we come back to same widget
			// FIXME: We should probably compare the whole buffer to be on the safety side. Comparing buf (utf8) and edit_state.Text (wchar).
			recycle_state := (edit_state.Id == id) && (prev_len_w == edit_state.CurLenW)
			if recycle_state {
				// Recycle existing cursor/selection/undo stack but clamp position
				// Note a single mouse click will override the cursor/position immediately by calling stb_textedit_click handler.
				edit_state.CursorClamp()
			} else {
				edit_state.Id = id
				edit_state.ScrollX = 0.0
				edit_state.StbState.Init(!is_multiline)
				if !is_multiline && focus_requested_by_code {
					select_all = true
				}
			}

			if flags&InputTextFlagsAlwaysInsertMode != 0 {
				edit_state.StbState.SetInsertMode(true)
			}
			if !is_multiline && (focus_requested_by_tab || (user_clicked && io.KeyCtrl)) {
				select_all = true
			}
		}
		c.SetActiveID(id, window)
		c.SetFocusID(id, window)
		c.FocusWindow(window)
		if !is_multiline && flags&InputTextFlagsCallbackHistory == 0 {
			c.ActiveIdAllowNavDirFlags |= ((1 << uint(DirUp)) | (1 << uint(DirDown)))
		}
	} else if io.MouseClicked[0] {
		// Release focus when we click outside
		clear_active_id = true
	}

	value_changed := false
	enter_pressed := false
	if c.ActiveId == id {
		if !is_editable && !c.ActiveIdIsJustActivated {
			// When read-only we always use the live data passed to the function
			if len(buf) > len(edit_state.Text) {
				edit_state.Text = append(edit_state.Text, make([]rune, len(buf)-len(edit_state.Text))...)
			}
			edit_state.CurLenW = len(edit_state.Text)
			edit_state.CurLenA = len(buf)
			edit_state.CursorClamp()
		}

		edit_state.BufSizeA = len(buf)

		// Although we are active we don't prevent mouse from hovering other elements unless we are interacting right now with the widget.
		// Down the line we should have a cleaner library-wide concept of Selected vs Active.
		c.ActiveIdAllowOverlap = !io.MouseDown[0]
		c.WantTextInputNextFrame = 1

		// Edit in progress
		mouse_x := (io.MousePos.X - frame_bb.Min.X - style.FramePadding.X) + edit_state.ScrollX
		mouse_y := c.FontSize * 0.5
		if is_multiline {
			mouse_y = io.MousePos.Y - draw_window.DC.CursorPos.Y - style.FramePadding.Y
		}
		// OS X style: Double click selects by word instead of selecting whole text
		osx_double_click_selects_words := io.OptMacOSXBehaviors
		if select_all || (hovered && !osx_double_click_selects_words && io.MouseDoubleClicked[0]) {
			edit_state.SelectAll()
			edit_state.SelectedAllMouseLock = true
		} else if hovered && osx_double_click_selects_words && io.MouseDoubleClicked[0] {
			// Select a word only, OS X style (by simulating keystrokes)
			edit_state.OnKeyPressed(stbte.K_WORDLEFT)
			edit_state.OnKeyPressed(stbte.K_WORDRIGHT | stbte.K_SHIFT)
		} else if io.MouseClicked[0] && !edit_state.SelectedAllMouseLock {
			if hovered {
				edit_state.StbState.Click(edit_state, mouse_x, mouse_y)
				edit_state.CursorAnimReset()
			}
		} else if io.MouseDown[0] && !edit_state.SelectedAllMouseLock && (io.MouseDelta.X != 0.0 || io.MouseDelta.Y != 0.0) {
			edit_state.StbState.Drag(edit_state, mouse_x, mouse_y)
			edit_state.CursorAnimReset()
			edit_state.CursorFollow = true
		}

		if io.InputCharacters[0] != 0 {
			// Process text input (before we check for Return because using some IME will effectively send a Return?)
			// We ignore CTRL inputs, but need to allow ALT+CTRL as some keyboards (e.g. German) use AltGR (which _is_ Alt+Ctrl) to input certain characters.
			if !(io.KeyCtrl && !io.KeyAlt) && is_editable && !user_nav_input_start {
			}
			// Consume characters
			for i := range c.IO.InputCharacters {
				c.IO.InputCharacters[i] = 0
			}
		}
	}

	cancel_edit := false
	if c.ActiveId == id && !c.ActiveIdIsJustActivated && !clear_active_id {
		// Handle key-presses
		k_mask := 0
		if io.KeyShift {
			k_mask = stbte.K_SHIFT
		}
		// OS X style: Shortcuts using Cmd/Super instead of Ctrl
		_, _, _, _, _ = is_undoable, value_changed, enter_pressed, cancel_edit, k_mask
	}

	return false
}

func (t *TextEditState) Init(ctx *Context) {
	*t = TextEditState{
		Ctx: ctx,
	}
}

func (t *TextEditState) CursorAnimReset() {
	// After a user-input the cursor stays on for a while without blinking
	t.CursorAnim = -0.30
}

func (t *TextEditState) CursorClamp() {
	t.StbState.SetCursor(mathutil.Min(t.StbState.Cursor(), t.CurLenW))
	t.StbState.SetSelectStart(mathutil.Min(t.StbState.SelectStart(), t.CurLenW))
	t.StbState.SetSelectEnd(mathutil.Min(t.StbState.SelectEnd(), t.CurLenW))
}

func (t *TextEditState) SelectAll() {
	t.StbState.SetSelectStart(0)
	t.StbState.SetCursor(t.CurLenW)
	t.StbState.SetSelectEnd(t.CurLenW)
	t.StbState.SetHasPreferredX(false)
}

func (t *TextEditState) OnKeyPressed(key int) {
	t.StbState.Key(t, key)
	t.CursorFollow = true
	t.CursorAnimReset()
}

func (t *TextEditState) GetChar(idx int) rune {
	return t.Text[idx]
}

func (t *TextEditState) GetWidth(line_start_idx, char_idx int) float64 {
	ctx := t.Ctx

	const STB_TEXTEDIT_GETWIDTH_NEWLINE = -1
	c := t.Text[line_start_idx+char_idx]
	if c == '\n' {
		return STB_TEXTEDIT_GETWIDTH_NEWLINE
	}
	return ctx.Font.GetCharAdvance(c) * (ctx.FontSize / ctx.Font.FontSize)
}

func (t *TextEditState) InsertChars(pos int, new_text []rune) bool {
	t.Text = append(t.Text[:pos], append(new_text, t.Text[pos:]...)...)
	t.CurLenW += len(new_text)
	t.CurLenA += TextCountUtf8BytesFromStr(new_text)
	return true
}

func (t *TextEditState) DeleteChars(pos, n int) {
	// We maintain our buffer length in both UTF-8 and wchar formats
	t.CurLenA -= TextCountUtf8BytesFromStr(t.Text[pos:])
	t.CurLenW -= n

	// Offset remaining text
	copy(t.Text[pos:pos+n], t.Text[pos+n:])
	t.Text = t.Text[:len(t.Text)-n]
}

func (t *TextEditState) LayoutRow(r *stbte.TextEditRow, line_start_idx int) {
	ctx := t.Ctx
	text := t.Text
	size, text_remaining, _ := ctx.InputTextCalcTextSizeW(text[line_start_idx:], true)
	r.SetX0(0.0)
	r.SetX1(size.X)
	r.SetBaselineYDelta(size.Y)
	r.SetYMin(0)
	r.SetYMax(size.Y)
	r.SetNumChars(text_remaining)
}

func (t *TextEditState) Len() int {
	return t.CurLenW
}

func (t *TextEditState) MoveWordLeft(n int) int {
	return 0
}

func (t *TextEditState) MoveWordRight(n int) int {
	return 0
}

func (c *Context) InputTextCalcTextSizeW(text []rune, stop_on_new_line bool) (text_size f64.Vec2, remaining int, out_offset f64.Vec2) {
	return
}