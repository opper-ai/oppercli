package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/opper-ai/oppercli/opperai"
)

type ListIndexesCommand struct {
	Filter string
}

type CreateIndexCommand struct {
	Name string
}

type DeleteIndexCommand struct {
	Name string
}

type GetIndexCommand struct {
	Name string
}

type QueryIndexCommand struct {
	Name   string
	Query  string
	Filter string
}

type AddToIndexCommand struct {
	Name     string
	Key      string
	Content  string
	Metadata string
}

type UploadToIndexCommand struct {
	Name     string
	FilePath string
}

func (c *ListIndexesCommand) Execute(ctx context.Context, client *opperai.Client) error {
	indexes, err := client.Indexes.List("")
	if err != nil {
		return err
	}

	for _, index := range indexes {
		fmt.Println(index.Name)
	}
	return nil
}

func (c *CreateIndexCommand) Execute(ctx context.Context, client *opperai.Client) error {
	index, err := client.Indexes.Create(c.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Created index: %s\n", index.Name)
	return nil
}

func (c *DeleteIndexCommand) Execute(ctx context.Context, client *opperai.Client) error {
	err := client.Indexes.Delete(c.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted index: %s\n", c.Name)
	return nil
}

func (c *GetIndexCommand) Execute(ctx context.Context, client *opperai.Client) error {
	index, err := client.Indexes.Get(c.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Index: %s\n", index.Name)
	fmt.Printf("Created: %s\n", index.CreatedAt.Format(time.RFC3339))

	if len(index.Files) > 0 {
		fmt.Println("\nIndexed Files:")
		fmt.Printf("%-50s %-10s %-15s\n", "Name", "Size", "Status")
		fmt.Println(strings.Repeat("-", 75))

		for _, file := range index.Files {
			fmt.Printf("%-50s %-10d %-15s\n",
				truncateString(file.OriginalFilename, 47),
				file.Size,
				file.IndexStatus,
			)
		}
	} else {
		fmt.Println("\nNo files indexed yet")
	}

	return nil
}

func (c *QueryIndexCommand) Execute(ctx context.Context, client *opperai.Client) error {
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Filter), &filter); err != nil {
		return fmt.Errorf("invalid filter JSON: %v", err)
	}

	results, err := client.Indexes.Query(c.Name, c.Query, nil) // TODO: Convert filter to []Filter
	if err != nil {
		return err
	}

	// Print results
	for _, result := range results {
		fmt.Printf("Score: %f\nContent: %s\n\n", result.Score, result.Content)
	}
	return nil
}

func (c *AddToIndexCommand) Execute(ctx context.Context, client *opperai.Client) error {
	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(c.Metadata), &metadata); err != nil {
		return fmt.Errorf("invalid metadata JSON: %v", err)
	}

	doc := opperai.Document{
		Key:      c.Key,
		Content:  c.Content,
		Metadata: metadata,
	}

	err := client.Indexes.Add(c.Name, doc)
	if err != nil {
		return err
	}

	fmt.Printf("Added document with key '%s' to index '%s'\n", c.Key, c.Name)
	return nil
}

func (c *UploadToIndexCommand) Execute(ctx context.Context, client *opperai.Client) error {
	if _, err := os.Stat(c.FilePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", c.FilePath)
	}

	err := client.Indexes.UploadFile(c.Name, c.FilePath)
	if err != nil {
		return err
	}

	fmt.Printf("Uploaded file '%s' to index '%s'\n", c.FilePath, c.Name)
	return nil
}

// Export the function for testing
func ExecuteListIndexes(ctx context.Context, client *opperai.Client, w io.Writer, args []string) error {
	indexes, err := client.Indexes.List("")
	if err != nil {
		return err
	}

	for _, index := range indexes {
		fmt.Fprintln(w, index.Name)
	}
	return nil
}
