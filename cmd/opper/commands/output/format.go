package output

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

// Table prints data in a formatted table
func Table(headers []string, rows [][]string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Print headers
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	// Print separator
	seps := make([]string, len(headers))
	for i, h := range headers {
		seps[i] = strings.Repeat("─", len(h))
	}
	fmt.Fprintln(w, strings.Join(seps, "\t"))

	// Print rows
	for _, row := range rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}

	w.Flush()
}

// Plain prints data in plain text format
func Plain(w io.Writer, items []string) {
	for _, item := range items {
		fmt.Fprintln(w, item)
	}
}
