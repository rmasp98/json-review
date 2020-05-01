package ui_test

import (
	"kube-review/ui"
	"testing"
)

func TestReturnsZeroDimensionIfNotResized(t *testing.T) {
	views := ui.GetWindow()
	assertEqualDimensions(t, views, zeroDimensions)
}

func TestReturnsFullScaledSizeOfAllLayouts(t *testing.T) {
	views := ui.GetWindow()
	views.Resize(100, 100)
	assertEqualDimensions(t, views, expectedFull)
}

func TestTooNarrowForPanel(t *testing.T) {
	views := ui.GetWindow()
	views.Resize(45, 100)
	assertEqualDimensions(t, views, expectedPanelless)
}

func TestTooShortForTextboxes(t *testing.T) {
	views := ui.GetWindow()
	views.Resize(100, 5)
	assertEqualDimensions(t, views, expectedTBless)
}

func TestTooNarrowForEverything(t *testing.T) {
	views := ui.GetWindow()
	views.Resize(2, 100)
	assertEqualDimensions(t, views, zeroDimensions)
}

func TestTooShortForEverything(t *testing.T) {
	views := ui.GetWindow()
	views.Resize(100, 2)
	assertEqualDimensions(t, views, zeroDimensions)
}

func assertEqualDimensions(t *testing.T, layouts *ui.Window, expected []ui.Dimensions) {
	for i := 0; i < 4; i++ {
		if layouts.GetDimensions(ui.ViewEnum(i)) != expected[i] {
			t.Errorf(`Got "%+v" but expected "%+v" for %+v`, layouts.GetDimensions(ui.ViewEnum(i)), expected[i], ui.ViewEnum(i))
		}
	}
}

var (
	zeroDimensions = []ui.Dimensions{
		ui.Dimensions{0, 0, 0, 0},
		ui.Dimensions{0, 0, 0, 0},
		ui.Dimensions{0, 0, 0, 0},
		ui.Dimensions{0, 0, 0, 0},
	}
	expectedFull = []ui.Dimensions{
		ui.Dimensions{1, 1, 20, 96},
		ui.Dimensions{21, 1, 99, 3},
		ui.Dimensions{21, 4, 99, 96},
		ui.Dimensions{1, 97, 99, 99},
	}

	expectedPanelless = []ui.Dimensions{
		ui.Dimensions{0, 0, 0, 0},
		ui.Dimensions{1, 1, 44, 3},
		ui.Dimensions{1, 4, 44, 96},
		ui.Dimensions{1, 97, 44, 99},
	}

	expectedTBless = []ui.Dimensions{
		ui.Dimensions{1, 1, 20, 4},
		ui.Dimensions{0, 0, 0, 0},
		ui.Dimensions{21, 1, 99, 4},
		ui.Dimensions{0, 0, 0, 0},
	}
)
