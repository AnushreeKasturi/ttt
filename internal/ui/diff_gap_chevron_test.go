package ui

import (
	"strings"
	"testing"

	"github.com/eugenioenko/ttt/internal/core/diff"
)

func gapRowWith(t *testing.T, rows []string, where string) string {
	t.Helper()
	y := rowWith(rows, "⋯ 11 lines ⋯")
	if y < 0 {
		t.Fatalf("%s: no gap row:\n%s", where, strings.Join(rows, "\n"))
	}
	return rows[y]
}

func renderWidgetRows(w Widget, width, height int) []string {
	w.SetRect(Rect{W: width, H: height})
	grid := makeGrid(width, height)
	w.Render(NewRenderSurface(grid, Rect{W: width, H: height}))
	rows := make([]string, height)
	for y := range grid {
		var sb strings.Builder
		for _, c := range grid[y] {
			sb.WriteRune(c.Ch)
		}
		rows[y] = sb.String()
	}
	return rows
}

func TestDiffGapRowsFollowConfiguredChevron(t *testing.T) {
	old, cur := gapFile()
	check := func(where, row string) {
		t.Helper()
		if !strings.ContainsRune(row, '»') || strings.ContainsRune(row, '▶') {
			t.Fatalf("%s: gap row %q, want the configured '»'", where, row)
		}
	}

	for _, mode := range []DiffMode{DiffModeUnified, DiffModeSplit} {
		f := newLiveDiffFixture(t, old, cur, mode, 100, 20)
		f.d.SetChevrons('»', 'v')
		check("editable", gapRowWith(t, f.render(), "editable"))
	}

	ro := NewDiffEditorWidget("f.txt", diff.Parse(diff.Generate(old, cur, "f.txt")), old, cur, false)
	ro.SetUnified(true)
	ro.SetChevrons('»', 'v')
	check("read-only", gapRowWith(t, renderWidgetRows(ro, 100, 20), "read-only"))

	detail := NewCommitDetailWidget("/repo", "full", "abc1234", false)
	detail.SetDetail("Subject", []CommitDetailFile{{Path: "f.txt", Diff: diff.Parse(diff.Generate(old, cur, "f.txt"))}}, "")
	detail.SetChevrons('»', 'v')
	check("commit details", gapRowWith(t, renderWidgetRows(detail, 100, 30), "commit details"))
}
