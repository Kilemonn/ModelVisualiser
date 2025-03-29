package main

import (
	"context"
	"flag"
	"fmt"
	"slices"
	"strings"

	"github.com/Kilemonn/ModelVisualiser/consts"
	"github.com/Kilemonn/ModelVisualiser/visualiser"
	"github.com/goccy/go-graphviz"
)

func main() {
	inputFile := flag.String(consts.INPUT_FILE, "", "Input File")
	outputFormat := flag.String(consts.OUTPUT_FORMAT, string(graphviz.PNG), "Output image format")

	flag.Parse()

	if inputFile == nil || *inputFile == "" {
		fmt.Printf("No input file provided.\n")
		flag.Usage()
		return
	}

	visualiser, err := visualiser.NewVisualiser(context.Background())
	if err != nil {
		fmt.Printf("Failed to create visualiser. Error: %s", err.Error())
	}
	defer visualiser.Close()

	graph, err := visualiser.FromFile(*inputFile)
	if err != nil {
		fmt.Printf("Failed to generate graph from input file [%s] with error: [%s]\n", *inputFile, err.Error())
		return
	}
	defer graph.Close()

	outputFile := outputFileName(*inputFile, *outputFormat)
	err = visualiser.ToFile(graph, outputFile, extensionToFormat(*outputFormat, graphviz.PNG))
	if err != nil {
		fmt.Printf("Failed to write generated graph to file [%s] with error: [%s]\n", outputFile, err.Error())
		return
	}
}

func outputFileName(inputFilename string, outputFormat string) string {
	index := strings.LastIndex(inputFilename, consts.PERIOD)
	if index == -1 {
		return inputFilename + consts.PERIOD + outputFormat
	}
	return inputFilename[:index] + consts.PERIOD + outputFormat
}

func extensionToFormat(outputFormat string, defaultFormat graphviz.Format) graphviz.Format {
	lower := strings.ToLower(outputFormat)
	formats := []graphviz.Format{graphviz.PNG, graphviz.JPG, graphviz.XDOT, graphviz.SVG}

	formatStrings := make([]string, 0)
	for _, format := range formats {
		formatStrings = append(formatStrings, string(format))
	}
	index := slices.Index(formatStrings, lower)

	if index == -1 {
		return defaultFormat
	}
	return formats[index]
}
