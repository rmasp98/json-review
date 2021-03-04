package uiold_test

import (
	"testing"

	"github.com/rmasp98/kube-review/uiold"
)

func TestReturnsZeroDimensionIfNotResized(t *testing.T) {
	views := uiold.NewWindow(panelRelativeWidth, border, tbBaseBuffer)
	assertEqualDimensions(t, views, zeroDimensions)
}

func TestReturnsFullScaledSizeOfAllLayouts(t *testing.T) {
	views := uiold.NewWindow(panelRelativeWidth, border, tbBaseBuffer)
	views.Resize(100, 100, 3)
	assertEqualDimensions(t, views, expectedFull)
}

func TestTooNarrowForPanel(t *testing.T) {
	views := uiold.NewWindow(panelRelativeWidth, border, tbBaseBuffer)
	views.Resize(45, 100, 3)
	assertEqualDimensions(t, views, expectedPanelless)
}

func TestTooShortForTextboxes(t *testing.T) {
	views := uiold.NewWindow(panelRelativeWidth, border, tbBaseBuffer)
	views.Resize(100, 5, 3)
	assertEqualDimensions(t, views, expectedTBless)
}

func TestTooNarrowForEverything(t *testing.T) {
	views := uiold.NewWindow(panelRelativeWidth, border, tbBaseBuffer)
	views.Resize(2, 100, 3)
	assertEqualDimensions(t, views, zeroDimensions)
}

func TestTooShortForEverything(t *testing.T) {
	views := uiold.NewWindow(panelRelativeWidth, border, tbBaseBuffer)
	views.Resize(100, 2, 3)
	assertEqualDimensions(t, views, zeroDimensions)
}

func TestCanAlterSearchHeightWithoutEffectingOtherViews(t *testing.T) {
	views := uiold.NewWindow(panelRelativeWidth, border, tbBaseBuffer)
	views.Resize(100, 100, 10)
	assertEqualDimensions(t, views, expectedExtendedSearch)
}

func assertEqualDimensions(t *testing.T, window uiold.Window, expected [][]int) {
	for i := 0; i < 4; i++ {
		x0, y0, x1, y1 := window.GetDimensions(uiold.ViewEnum(i))
		if x0 != expected[i][0] || y0 != expected[i][1] || x1 != expected[i][2] || y1 != expected[i][3] {
			t.Errorf(`Got "[%d,%d,%d,%d]" but expected "%+v" for %+v`, x0, y0, x1, y1, expected[i], uiold.ViewEnum(i))
		}
	}
}

var (
	panelRelativeWidth = 0.2
	border             = 1
	tbBaseBuffer       = 3
	zeroDimensions     = [][]int{
		[]int{0, 0, 0, 0},
		[]int{0, 0, 0, 0},
		[]int{0, 0, 0, 0},
		[]int{0, 0, 0, 0},
	}
	expectedFull = [][]int{
		[]int{1, 1, 20, 96},
		[]int{21, 1, 99, 3},
		[]int{21, 4, 99, 96},
		[]int{1, 97, 99, 99},
	}

	expectedPanelless = [][]int{
		[]int{0, 0, 0, 0},
		[]int{1, 1, 44, 3},
		[]int{1, 4, 44, 96},
		[]int{1, 97, 44, 99},
	}

	expectedTBless = [][]int{
		[]int{1, 1, 20, 4},
		[]int{0, 0, 0, 0},
		[]int{21, 1, 99, 4},
		[]int{0, 0, 0, 0},
	}
	expectedExtendedSearch = [][]int{
		[]int{1, 1, 20, 96},
		[]int{21, 1, 99, 10},
		[]int{21, 4, 99, 96},
		[]int{1, 97, 99, 99},
	}
)
