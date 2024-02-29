package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"google.golang.org/grpc"
	"log"
	"os"
)

// Define the node type
type NodeType struct {
	ID          string `json:"uid,omitempty"`
	A           string `json:"A,omitempty"`
	B           string `json:"B,omitempty"`
	ComboResult string `json:"ComboResult,omitempty"`
	// Add other fields as needed
}

func main() {
	// Connect to Dgraph server
	conn, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial Dgraph server: %v", err)
	}
	defer conn.Close()

	// Create a new Dgraph client
	client := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	// Define pagination parameters
	first := 1000 // Number of nodes to fetch per page
	offset := 0   // Initial offset

	// Create a new CSV file
	file, err := os.Create("export.csv")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	if err := writer.Write([]string{"A", "B", "ComboResult"}); err != nil {
		log.Fatalf("Failed to write CSV header: %v", err)
	}

	for {
		// Fetch nodes in batches
		nodes, err := fetchNodes(client, first, offset)
		if err != nil {
			log.Fatalf("Failed to fetch nodes: %v", err)
		}

		if len(nodes) == 0 {
			break
		}

		// Export nodes to CSV
		for _, node := range nodes {
			if err := writer.Write([]string{node.A, node.B, node.ComboResult}); err != nil {
				log.Fatalf("Failed to write CSV data: %v", err)
			}
		}

		offset += first
	}

	fmt.Println("Data exported to export.csv")
}

func fetchNodes(client *dgo.Dgraph, first, offset int) ([]NodeType, error) {
	// Define the query with pagination
	query := fmt.Sprintf(`
		{
			nodes(func: type(Combo), first: %d, offset: %d) {
				uid
				A
				B
				ComboResult
			}
		}
	`, first, offset)

	// Run the query
	resp, err := client.NewTxn().Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	// Unmarshal the response
	var data struct {
		Nodes []NodeType `json:"nodes"`
	}
	if err := json.Unmarshal(resp.GetJson(), &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return data.Nodes, nil
}
