package commands

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/opper-ai/oppercli/opperai"
)

type ListUsageCommand struct {
	StartDate    string
	EndDate      string
	ProjectName  string
	FunctionPath string
	Model        string
	Skip         int
	Limit        int
	Out          string
}

func (c *ListUsageCommand) Execute(ctx context.Context, client *opperai.Client) error {
	params := &opperai.UsageParams{
		StartDate:    c.StartDate,
		EndDate:      c.EndDate,
		ProjectName:  c.ProjectName,
		FunctionPath: c.FunctionPath,
		Model:        c.Model,
		Skip:         c.Skip,
		Limit:        c.Limit,
	}

	usage, err := client.Usage.List(ctx, params)
	if err != nil {
		return err
	}

	switch strings.ToLower(c.Out) {
	case "csv":
		w := csv.NewWriter(os.Stdout)
		defer w.Flush()

		// Write header
		header := []string{
			"Project",
			"Function",
			"Model",
			"Tokens Input",
			"Tokens Output",
			"Total Tokens",
			"Cost",
			"Created At",
		}
		if err := w.Write(header); err != nil {
			return fmt.Errorf("error writing CSV header: %v", err)
		}

		// Write data rows
		for _, item := range usage.Items {
			row := []string{
				item.ProjectName,
				item.FunctionPath,
				item.Model,
				fmt.Sprintf("%d", item.TokensInput),
				fmt.Sprintf("%d", item.TokensOutput),
				fmt.Sprintf("%d", item.TotalTokens),
				fmt.Sprintf("%.4f", item.Cost),
				item.CreatedAt.Format("2006-01-02 15:04:05"),
			}
			if err := w.Write(row); err != nil {
				return fmt.Errorf("error writing CSV row: %v", err)
			}
		}
		return nil
	case "":
		// Print stats
		fmt.Printf("Stats:\n")
		fmt.Printf("  Total Tokens Input:  %d\n", usage.Stats.TotalTokensInput)
		fmt.Printf("  Total Tokens Output: %d\n", usage.Stats.TotalTokensOutput)
		fmt.Printf("  Total Tokens:        %d\n", usage.Stats.TotalTokens)
		fmt.Printf("  Total Cost:          %.4f\n", usage.Stats.TotalCost)
		fmt.Printf("  Count:               %d\n", usage.Stats.Count)
		fmt.Printf("\n")

		// Print items
		fmt.Printf("Usage Items (Total: %d):\n", usage.Total)
		for _, item := range usage.Items {
			fmt.Printf("  Project: %s\n", item.ProjectName)
			fmt.Printf("  Function: %s\n", item.FunctionPath)
			fmt.Printf("  Model: %s\n", item.Model)
			fmt.Printf("  Tokens Input: %d\n", item.TokensInput)
			fmt.Printf("  Tokens Output: %d\n", item.TokensOutput)
			fmt.Printf("  Total Tokens: %d\n", item.TotalTokens)
			fmt.Printf("  Cost: %.4f\n", item.Cost)
			fmt.Printf("  Created At: %s\n", item.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("\n")
		}
		return nil
	default:
		return fmt.Errorf("unsupported output format: %s", c.Out)
	}
}
