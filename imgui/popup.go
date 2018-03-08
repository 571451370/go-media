package imgui

type PopupPositionPolicy int

const (
	ImGuiPopupPositionPolicy_Default PopupPositionPolicy = iota
	ImGuiPopupPositionPolicy_ComboBox
)

func (c *Context) IsPopupOpen(str_id string) bool {
	return len(c.OpenPopupStack) > len(c.CurrentPopupStack) && c.OpenPopupStack[len(c.CurrentPopupStack)].PopupId == c.CurrentWindow.GetID(str_id)
}