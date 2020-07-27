package internal

import (
	"fmt"
	"io"
	"strings"
)

type TableWriter struct {
	titles []interface{}
	rows   [][]interface{}
	max    []int
}

func NewTableWriter() *TableWriter {
	return &TableWriter{
		titles: []interface{}{},
		rows:   [][]interface{}{},
		max:    []int{},
	}
}

func (t *TableWriter) ClearData() {
	t.rows = [][]interface{}{}
	for i := range t.max {
		t.max[i] = len(t.titles[i].(string))
	}
}

func (t *TableWriter) AddRow(data []string) error {
	if len(data) != len(t.titles) {
		return fmt.Errorf("wrong number of columns")
	}

	row := make([]interface{}, len(data), len(data))
	for i, d := range data {
		row[i] = data[i]
		if len(d) > t.max[i] {
			t.max[i] = len(d)
		}
	}
	t.rows = append(t.rows, row)
	return nil
}

func (t *TableWriter) SetColumnTitles(titles []string) {
	t.titles = []interface{}{}
	for i, title := range titles {
		t.titles = append(t.titles, title)
		if len(t.max) <= i {
			t.max = append(t.max, len(title))
		}
	}
}

func (t *TableWriter) Write(w io.Writer) {
	format := "|"
	separator := "+"
	for _, m := range t.max {
		format = format + fmt.Sprintf(" %%-%ds |", m)
		separator = separator + fmt.Sprintf("-%s-+", strings.Repeat("-", m))
	}
	format = format + "\n"
	separator = separator + "\n"
	fmt.Fprintf(w, separator)
	fmt.Fprintf(w, format, t.titles...)
	fmt.Fprintf(w, separator)
	for _, row := range t.rows {
		fmt.Fprintf(w, format, row...)
	}
	fmt.Fprintf(w, separator)
}
