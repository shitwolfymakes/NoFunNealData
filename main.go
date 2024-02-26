package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"google.golang.org/grpc"
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

	// Define the query
	query := `
        {
            nodes(func: type(Combo)) {
                uid
                A
				B
				ComboResult
            }
        }
    `

	// Run the query
	resp, err := client.NewReadOnlyTxn().Query(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}

	// Unmarshal the response
	var data struct {
		Nodes []NodeType `json:"nodes"`
	}
	if err := json.Unmarshal(resp.GetJson(), &data); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Export the data to JSON
	jsonData, err := json.MarshalIndent(data.Nodes, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Save JSON data to a file
	file, err := os.Create("export.json")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	fmt.Println("Data exported to export.json")
}
