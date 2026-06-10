package ui

import (
	"bytes"
	"strings"
	"testing"
)

func TestTableEmpty(t *testing.T) {
	tbl := NewTable()
	if !tbl.Empty() {
		t.Error("new table should be empty")
	}

	tbl.Header("NAME")
	tbl.Row(Text("foo"))
	if tbl.Empty() {
		t.Error("table with rows should not be empty")
	}
}

func TestTableWriteHeaders(t *testing.T) {
	tbl := NewTable()
	tbl.Header("NAME", "URI")
	tbl.Row(Text("lofi"), Text("spotify:playlist:abc123").Fixed())

	var buf bytes.Buffer
	tbl.write(&buf, 1, 200)
	out := buf.String()

	if !strings.Contains(out, "NAME") {
		t.Error("output should contain header NAME")
	}
	if !strings.Contains(out, "URI") {
		t.Error("output should contain header URI")
	}
}

func TestTableWriteRows(t *testing.T) {
	tbl := NewTable()
	tbl.Header("NAME")
	tbl.Row(Text("abbey road"))
	tbl.Row(Text("dark side"))

	var buf bytes.Buffer
	tbl.write(&buf, 2, 200)
	out := buf.String()

	if !strings.Contains(out, "abbey road") {
		t.Error("output should contain first row")
	}
	if !strings.Contains(out, "dark side") {
		t.Error("output should contain second row")
	}
}

func TestTableTruncatesNonFixed(t *testing.T) {
	tbl := NewTable()
	tbl.Header("NAME", "ID")
	// wide non-fixed name, short fixed id
	tbl.Row(Text("a very long track name that should get truncated"), Text("abc").Fixed())

	var buf bytes.Buffer
	tbl.write(&buf, 1, 40)
	out := buf.String()

	if !strings.Contains(out, "...") {
		t.Error("long non-fixed column should be truncated with ...")
	}
}

func TestTableFixedNotTruncated(t *testing.T) {
	tbl := NewTable()
	tbl.Header("NAME", "URI")
	uri := "spotify:playlist:37i9dQZF1DX3Ogo9pFvBkY"
	tbl.Row(Text("lofi"), Text(uri).Fixed())

	var buf bytes.Buffer
	tbl.write(&buf, 1, 40)
	out := buf.String()

	if !strings.Contains(out, uri) {
		t.Error("fixed column should never be truncated")
	}
}

func TestTableTitle(t *testing.T) {
	tbl := NewTable()
	tbl.Title(Text("tracks"))
	tbl.Header("NAME")
	tbl.Row(Text("song"))

	var buf bytes.Buffer
	tbl.write(&buf, 1, 200)
	out := buf.String()

	if !strings.Contains(out, "Showing 1 of 1 tracks") {
		t.Errorf("output should contain title line, got:\n%s", out)
	}
}

func TestTableShownCount(t *testing.T) {
	tbl := NewTable()
	tbl.Title(Text("tracks"))
	tbl.Header("NAME")
	tbl.Row(Text("song one"))
	tbl.Row(Text("song two"))
	tbl.Row(Text("song three"))

	var buf bytes.Buffer
	tbl.write(&buf, 2, 200)
	out := buf.String()

	if !strings.Contains(out, "Showing 2 of 3 tracks") {
		t.Errorf("output should show correct shown/total count, got:\n%s", out)
	}
}
