package utils

import (
	"PageRank/models"
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Convert -> Converts txt file that contains nodes list of adjacency
//			  into array of nodes pointer

func Convert(graphPath string) []*models.Node {

	var graph []*models.Node

	inputFile, err := os.Open(graphPath)
	if err != nil {
		fmt.Println("Error during file opening: ", err)
		return nil
	}
	defer func(inputFile *os.File) {
		err := inputFile.Close()
		if err != nil {
			fmt.Println("Something went wrong during file closing")
		}
	}(inputFile)

	scanner := bufio.NewScanner(inputFile)
	// Divides each row into node and outgoing links and append them to graph array
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		node := parts[0]
		outgoingLinks := parts[2:]
		outgoingLinks = FixStrings(outgoingLinks)
		actualNode := models.Node{ID: node, OutLinks: outgoingLinks, PageRank: 1}
		graph = append(graph, &actualNode)
	}

	// Initialize PageRank to 1/N for each node
	for i := 0; i < len(graph); i++ {
		graph[i].PageRank = 1.0 / float64(len(graph))
	}

	return graph
}

// FixStrings -> removes all commas into array of string
func FixStrings(list []string) []string {
	var correctArray []string
	for _, s := range list {
		correctArray = append(correctArray, strings.ReplaceAll(s, ",", ""))
	}
	return correctArray
}
