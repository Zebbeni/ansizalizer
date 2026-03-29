package adjust

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/style"
)

var (
	inputStyle = style.BgStyle().Width(20).AlignHorizontal(lipgloss.Left)
)

func (m Model) drawBrightnessForm() string {
	prompt, text := m.getInputColors(BrightnessForm)
	m.brightnessInput.SetWidth(4)

	styles := m.brightnessInput.Styles()
	styles.Focused.Prompt = styles.Focused.Prompt.Foreground(prompt)
	styles.Blurred.Prompt = styles.Blurred.Prompt.Foreground(prompt)
	styles.Focused.Text = styles.Focused.Text.Foreground(text)
	styles.Blurred.Text = styles.Blurred.Text.Foreground(text)
	styles.Cursor.Color = text
	styles.Cursor.Blink = m.brightnessInput.Focused()
	m.brightnessInput.SetStyles(styles)

	return inputStyle.Render(m.brightnessInput.View())
}

func (m Model) drawContrastForm() string {
	prompt, text := m.getInputColors(ContrastForm)
	m.contrastInput.SetWidth(4)

	styles := m.contrastInput.Styles()
	styles.Focused.Prompt = styles.Focused.Prompt.Foreground(prompt)
	styles.Blurred.Prompt = styles.Blurred.Prompt.Foreground(prompt)
	styles.Focused.Text = styles.Focused.Text.Foreground(text)
	styles.Blurred.Text = styles.Blurred.Text.Foreground(text)
	styles.Cursor.Color = text
	styles.Cursor.Blink = m.contrastInput.Focused()
	m.contrastInput.SetStyles(styles)

	return inputStyle.Render(m.contrastInput.View())
}

func (m Model) drawSlider(state State, value int) string {
	sliderWidth := m.width - 2
	if sliderWidth < 3 {
		sliderWidth = 3
	}

	pos := int(float64(value+100) / 200.0 * float64(sliderWidth-1))
	if pos < 0 {
		pos = 0
	}
	if pos >= sliderWidth {
		pos = sliderWidth - 1
	}

	left := strings.Repeat("-", pos)
	right := strings.Repeat("-", sliderWidth-1-pos)

	sliderColor := style.DimmedColor1
	if m.IsActive && m.focus == state {
		sliderColor = style.SelectedColor1
	}

	sliderStyle := style.BgStyle().Foreground(sliderColor).PaddingLeft(1)
	return sliderStyle.Render(left + "|" + right)
}

func (m Model) getInputColors(state State) (color.Color, color.Color) {
	if m.IsActive && m.focus == state {
		if m.active == state {
			return style.NormalColor1, style.SelectedColor1
		}
		return style.SelectedColor1, style.NormalColor1
	}
	return style.DimmedColor1, style.DimmedColor1
}
