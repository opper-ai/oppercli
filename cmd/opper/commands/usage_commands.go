package commands

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/opper-ai/oppercli/opperai"
)

type ListUsageCommand struct {
	FromDate    string
	ToDate      string
	Granularity string
	Fields      []string
	GroupBy     []string
	Out         string
}

func formatCost(cost string) string {
	if f, err := strconv.ParseFloat(cost, 64); err == nil {
		return fmt.Sprintf("%.6f", f)
	}
	return cost
}

func (c *ListUsageCommand) Execute(ctx context.Context, client *opperai.Client) error {
	params := &opperai.UsageParams{
		FromDate:    c.FromDate,
		ToDate:      c.ToDate,
		Granularity: c.Granularity,
		Fields:      c.Fields,
		GroupBy:     c.GroupBy,
	}

	usage, err := client.Usage.List(ctx, params)
	if err != nil {
		return err
	}

	events := *usage

	switch strings.ToLower(c.Out) {
	case "csv":
		w := csv.NewWriter(os.Stdout)
		defer w.Flush()

		// Build headers based on available fields
		headers := []string{"Time Bucket", "Cost", "Count"}
		if len(events) > 0 {
			// Get all dynamic field names from the first event
			var dynamicFields []string
			for k := range events[0].Fields {
				dynamicFields = append(dynamicFields, k)
			}
			// Sort field names for consistent output
			sort.Strings(dynamicFields)
			headers = append(headers, dynamicFields...)
		}

		if err := w.Write(headers); err != nil {
			return fmt.Errorf("error writing CSV header: %v", err)
		}

		// Write data rows
		for _, event := range events {
			row := []string{
				event.TimeBucket,
				formatCost(event.Cost),
				fmt.Sprintf("%d", event.Count),
			}

			// Add dynamic fields in the same order as headers
			for _, h := range headers[3:] { // Skip the first 3 standard fields
				if v, ok := event.Fields[h]; ok {
					row = append(row, fmt.Sprintf("%v", v))
				} else {
					row = append(row, "")
				}
			}

			if err := w.Write(row); err != nil {
				return fmt.Errorf("error writing CSV row: %v", err)
			}
		}
		return nil

	default:
		fmt.Printf("Usage Events:\n\n")
		for _, event := range events {
			fmt.Printf("Time Bucket: %s\n", event.TimeBucket)
			fmt.Printf("Cost: %s\n", formatCost(event.Cost))
			fmt.Printf("Count: %d\n", event.Count)

			// Sort field names for consistent output
			var fields []string
			for k := range event.Fields {
				fields = append(fields, k)
			}
			sort.Strings(fields)

			// Print dynamic fields
			for _, k := range fields {
				fmt.Printf("%s: %v\n", k, event.Fields[k])
			}
			fmt.Println()
		}
		return nil
	}
}
