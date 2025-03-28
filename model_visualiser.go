package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/Kilemonn/ModelVisualiser/consts"
	"github.com/Kilemonn/ModelVisualiser/model_visualiser"
)

func main() {
	inputFile := flag.String(consts.INPUT_FILE, "", "Input File")
	outputFormat := flag.String(consts.OUTPUT_FORMAT, consts.DEFAULT_OUTPUT_FORMAT, "Output image format")

	flag.Parse()

	if inputFile == nil || *inputFile == "" {
		fmt.Printf("No input file provided.\n")
		flag.Usage()
		return
	}

	visualiser, err := model_visualiser.NewModelVisualiser(context.Background())
	if err != nil {
		fmt.Printf("Failed to create visualiser. Error: %s", err.Error())
	}
	defer visualiser.Close()

	g, err := visualiser.FromFile(*inputFile)
	if err != nil {
		fmt.Printf("Failed to generate graph from input file [%s] with error: [%s]\n", *inputFile, err.Error())
		return
	}

	outputFile := outputFileName(*inputFile, *outputFormat)
	err = visualiser.ToFile(g, outputFile)
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
