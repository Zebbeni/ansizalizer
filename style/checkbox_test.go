package style

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestMain(m *testing.M) {
	SetTheme(LightOnDark)
	m.Run()
}

func TestRenderCheckboxChecked(t *testing.T) {
	result := RenderCheckbox("Test", true, false)
	if !strings.Contains(result, "☑") {
		t.Error("checked checkbox should contain ☑")
	}
}

func TestRenderCheckboxUnchecked(t *testing.T) {
	result := RenderCheckbox("Test", false, false)
	if !strings.Contains(result, "☐") {
		t.Error("unchecked checkbox should contain ☐")
	}
}

func TestRenderCheckboxContainsLabel(t *testing.T) {
	result := RenderCheckbox("MyLabel", true, false)
	if !strings.Contains(result, "MyLabel") {
		t.Errorf("checkbox should contain label, got %q", result)
	}
}

func TestRenderCheckboxFixedWidthAlignment(t *testing.T) {
	short := RenderCheckboxFixedWidth("A", true, false, 20)
	long := RenderCheckboxFixedWidth("LongLabelHere", true, false, 20)

	shortW := lipgloss.Width(short)
	longW := lipgloss.Width(long)

	if shortW != longW {
		t.Errorf("fixed width checkboxes should have same width: short=%d, long=%d", shortW, longW)
	}
}
