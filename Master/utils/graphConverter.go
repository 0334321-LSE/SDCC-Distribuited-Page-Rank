package utils

import (
	"Master/models"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
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
		nodeID, err := strconv.Atoi(node)
		if err != nil {
			log.Fatalf("Error occured during txt scanning %v", err)
		}
		outgoingLinks := parts[2:]
		outgoingLinks = FixStrings(outgoingLinks)
		outgoingLinksInt := ConvertOutLinks(outgoingLinks)
		actualNode := models.Node{ID: nodeID, OutLinks: outgoingLinksInt, PageRank: 1}
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

// ConvertOutLinks -> converts all the out-links string to integer
func ConvertOutLinks(list []string) []int {
	var correctArray []int
	for _, s := range list {
		nodeID, err := strconv.Atoi(s)
		if err != nil {
			log.Fatalf("Error occured during outlink convertion %v", err)
		}
		correctArray = append(correctArray, nodeID)
	}
	return correctArray
}
