package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/opper-ai/oppercli/opperai"
)

type ListTracesCommand struct{}

func (c *ListTracesCommand) Execute(ctx context.Context, client *opperai.Client) error {
	traces, err := client.Traces.List(ctx)
	if err != nil {
		return fmt.Errorf("error listing traces: %w", err)
	}

	// Print header
	fmt.Printf("\n%-36s  %-20s  %-10s  %-8s  %-15s  %-20s  %s\n",
		"UUID", "NAME", "STATUS", "SCORE", "DURATION", "PROJECT", "START TIME")
	fmt.Printf("%s\n", strings.Repeat("─", 125))

	for _, trace := range traces {
		// Calculate average score if there are scores
		scoreStr := "N/A"
		if len(trace.Scores) > 0 {
			var totalScore float64
			for _, score := range trace.Scores {
				totalScore += score.Score
			}
			avgScore := totalScore / float64(len(trace.Scores))
			scoreStr = fmt.Sprintf("%.0f%%", avgScore)
		}

		fmt.Printf("%-36s  %-20s  %-10s  %-8s  %-15.2fms  %-20s  %s\n",
			trace.UUID,
			truncateString(trace.Name, 20),
			trace.Status,
			scoreStr,
			trace.DurationMs,
			truncateString(trace.Project.Name, 20),
			trace.StartTime.Format(time.RFC3339),
		)
	}
	fmt.Println()
	return nil
}

type GetTraceCommand struct {
	TraceID string
}

func (c *GetTraceCommand) Execute(ctx context.Context, client *opperai.Client) error {
	trace, err := client.Traces.Get(ctx, c.TraceID)
	if err != nil {
		return fmt.Errorf("error getting trace: %w", err)
	}

	// Print trace details
	fmt.Printf("\nTrace: %s\n", trace.UUID)
	fmt.Printf("Name: %s\n", trace.Name)
	fmt.Printf("Status: %s\n", trace.Status)
	fmt.Printf("Project: %s\n", trace.Project.Name)
	fmt.Printf("Duration: %.2fms\n", trace.DurationMs)
	fmt.Printf("Start Time: %s\n", trace.StartTime.Format(time.RFC3339))
	fmt.Printf("End Time: %s\n", trace.EndTime.Format(time.RFC3339))
	if trace.Input != "" {
		fmt.Printf("Input: %s\n", trace.Input)
	}
	if trace.Output != nil {
		fmt.Printf("Output: %s\n", *trace.Output)
	}
	if len(trace.Scores) > 0 {
		var totalScore float64
		for _, score := range trace.Scores {
			totalScore += score.Score
		}
		avgScore := totalScore / float64(len(trace.Scores))
		fmt.Printf("Score: %.0f%%\n", avgScore)
	}

	// Print spans
	if len(trace.Spans) > 0 {
		fmt.Printf("\nSpans:\n")
		const (
			prefixWidth   = 12 // Width reserved for the tree prefix
			uuidWidth     = 48 // Width for UUID
			nameWidth     = 40 // Width for name
			scoreWidth    = 8  // Width for score
			durationWidth = 12 // Width for duration
			timeWidth     = 24 // Width for timestamp
		)

		// Print header with proper alignment
		fmt.Printf("%s%-*s  %-*s  %*s  %*s  %s\n",
			strings.Repeat(" ", prefixWidth),
			uuidWidth, "UUID",
			nameWidth, "NAME",
			scoreWidth, "SCORE",
			durationWidth, "DURATION",
			"START TIME")
		fmt.Printf("%s%s\n",
			strings.Repeat(" ", prefixWidth),
			strings.Repeat("─", uuidWidth+nameWidth+scoreWidth+durationWidth+timeWidth+8)) // 8 for spaces between columns

		// Create a map of parent UUID to child spans
		spansByParent := make(map[string][]*opperai.Span)
		var rootSpans []*opperai.Span

		// First pass: organize spans by parent
		for i := range trace.Spans {
			span := &trace.Spans[i]
			if span.ParentUUID == nil {
				rootSpans = append(rootSpans, span)
			} else {
				parentID := *span.ParentUUID
				spansByParent[parentID] = append(spansByParent[parentID], span)
			}
		}

		// Helper function to print span and its children recursively
		var printSpan func(span *opperai.Span, level int)
		printSpan = func(span *opperai.Span, level int) {
			// Create indentation based on level
			indent := strings.Repeat("    ", level)
			indentLen := len(indent)

			// Format score
			scoreStr := ""
			if span.Score != nil {
				scoreStr = fmt.Sprintf("%.0f%%", *span.Score)
			}

			// Convert duration from ms to s
			duration := span.DurationMs / 1000.0

			// Calculate remaining space for UUID to maintain column alignment
			remainingUUIDWidth := uuidWidth - indentLen
			if remainingUUIDWidth < 8 {
				remainingUUIDWidth = 8 // Minimum width for truncated UUID
			}

			// Print the main span line with aligned columns
			fmt.Printf("%s%-*s  %-*s  %*s  %*.3fs  %s\n",
				indent,
				remainingUUIDWidth, truncateString(span.UUID, remainingUUIDWidth),
				nameWidth, truncateString(span.Name, nameWidth),
				scoreWidth, scoreStr,
				durationWidth-1, duration, // -1 for the 's' suffix
				span.StartTime.Format(time.RFC3339),
			)

			// Print input/output with consistent indentation and alignment
			if span.Input != nil && *span.Input != "" {
				inputStr := strings.ReplaceAll(*span.Input, "\n", " ")
				fmt.Printf("%s    Input: %s\n",
					indent,
					inputStr)
			}
			if span.Output != nil && *span.Output != "" {
				outputStr := strings.ReplaceAll(*span.Output, "\n", " ")
				fmt.Printf("%s    Output: %s\n",
					indent,
					outputStr)
			}

			// Print child spans
			children := spansByParent[span.UUID]
			for _, child := range children {
				fmt.Println() // Add spacing between spans
				printSpan(child, level+1)
			}
		}

		// Print all root spans and their children
		for i, span := range rootSpans {
			if i > 0 {
				fmt.Println()
			}
			printSpan(span, 0)
		}
	}
	return nil
}