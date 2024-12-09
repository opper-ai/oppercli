package commands

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/guptarohit/asciigraph"
	"github.com/opper-ai/oppercli/opperai"
)

type ListUsageCommand struct {
	FromDate    string
	ToDate      string
	Granularity string
	Fields      []string
	GroupBy     []string
	Out         string
	Graph       string
}

func formatCost(cost string) string {
	if f, err := strconv.ParseFloat(cost, 64); err == nil {
		return fmt.Sprintf("%.6f", f)
	}
	return cost
}

func parseCost(cost string) float64 {
	if f, err := strconv.ParseFloat(cost, 64); err == nil {
		return f
	}
	return 0
}

func getGraphValue(event opperai.UsageEvent, graphType string) float64 {
	switch graphType {
	case "cost":
		return parseCost(event.Cost)
	default: // count
		return float64(event.Count)
	}
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

	// Sort events by time bucket
	sort.Slice(events, func(i, j int) bool {
		return events[i].TimeBucket < events[j].TimeBucket
	})

	switch {
	case c.Graph != "" && (c.Graph == "count" || c.Graph == "cost"):
		// If we have group by, we need to show multiple lines
		if len(c.GroupBy) > 0 {
			// Group events by the group by field
			groups := make(map[string][]opperai.UsageEvent)
			var groupNames []string
			for _, event := range events {
				for _, field := range c.GroupBy {
					if val, ok := event.Fields[field]; ok {
						groupName := fmt.Sprintf("%v", val)
						if _, exists := groups[groupName]; !exists {
							groupNames = append(groupNames, groupName)
						}
						groups[groupName] = append(groups[groupName], event)
					}
				}
			}

			// Sort group names for consistent output
			sort.Strings(groupNames)

			// Get all unique time buckets
			timeMap := make(map[string]bool)
			for _, events := range groups {
				for _, event := range events {
					timeMap[event.TimeBucket] = true
				}
			}
			var timeBuckets []string
			for t := range timeMap {
				timeBuckets = append(timeBuckets, t)
			}
			sort.Strings(timeBuckets)

			// Create data series for each group
			var series [][]float64
			var labels []string
			for _, name := range groupNames {
				// Create a map of time bucket to value for this group
				valueMap := make(map[string]float64)
				for _, event := range groups[name] {
					valueMap[event.TimeBucket] = getGraphValue(event, c.Graph)
				}

				// Create series with 0 for missing points
				var data []float64
				for _, t := range timeBuckets {
					if val, ok := valueMap[t]; ok {
						data = append(data, val)
					} else {
						data = append(data, 0)
					}
				}
				series = append(series, data)
				labels = append(labels, name)
			}

			// Get time labels
			var timeLabels []string
			for _, t := range timeBuckets {
				parsed, _ := time.Parse(time.RFC3339, t)
				timeLabels = append(timeLabels, parsed.Format("2006-01-02 15:04"))
			}

			// Plot multiple series
			metric := "count"
			if c.Graph == "cost" {
				metric = "cost"
			}
			fmt.Printf("\n%s over time by %s:\n\n", strings.Title(metric), strings.Join(c.GroupBy, ", "))
			graph := asciigraph.PlotMany(series,
				asciigraph.Height(15),
				asciigraph.Width(100),
				asciigraph.Caption("Time →"),
				asciigraph.SeriesColors(
					asciigraph.Red,
					asciigraph.Green,
					asciigraph.Blue,
					asciigraph.Yellow,
				),
				asciigraph.LabelColor(asciigraph.White),
			)

			// Print graph
			fmt.Println(graph)

			// Print legend
			fmt.Println("\nLegend:")
			colors := []asciigraph.AnsiColor{
				asciigraph.Red,
				asciigraph.Green,
				asciigraph.Blue,
				asciigraph.Yellow,
			}
			for i, name := range labels {
				color := colors[i%len(colors)]
				fmt.Printf("%s%s%s: %s\n",
					color.String(),
					"─────",
					asciigraph.White.String(),
					name,
				)
			}
			fmt.Println()

			// Print time labels
			if len(timeLabels) > 0 {
				fmt.Println("Time points:")
				for i, label := range timeLabels {
					fmt.Printf("%d: %s\n", i+1, label)
				}
				fmt.Println()
			}

		} else {
			// Single line graph
			var data []float64
			var timeLabels []string
			for _, event := range events {
				data = append(data, getGraphValue(event, c.Graph))
				t, _ := time.Parse(time.RFC3339, event.TimeBucket)
				timeLabels = append(timeLabels, t.Format("2006-01-02 15:04"))
			}

			metric := "count"
			if c.Graph == "cost" {
				metric = "cost"
			}
			fmt.Printf("\n%s over time:\n\n", strings.Title(metric))
			graph := asciigraph.Plot(data,
				asciigraph.Height(15),
				asciigraph.Width(100),
				asciigraph.Caption("Time →"),
			)
			fmt.Println(graph)

			// Print time labels
			if len(timeLabels) > 0 {
				fmt.Println("\nTime points:")
				for i, label := range timeLabels {
					fmt.Printf("%d: %s\n", i+1, label)
				}
				fmt.Println()
			}
		}
		return nil

	case strings.ToLower(c.Out) == "csv":
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
					if v == nil {
						row = append(row, "")
					} else {
						row = append(row, fmt.Sprintf("%v", v))
					}
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
