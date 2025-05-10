package main

import (
	"testing"

	"github.com/Kilemonn/ModelVisualiser/consts"
	"github.com/goccy/go-graphviz"
	"github.com/stretchr/testify/require"
)

func TestOutputFile_WithoutPeriod(t *testing.T) {
	output := "outfile"
	ext := "png"
	require.Equal(t, output+consts.PERIOD+ext, outputFileName(output, ext))
}

func TestOutputFile_WithPeriod(t *testing.T) {
	output := "outfile"
	ext := "png"
	require.Equal(t, output+consts.PERIOD+ext, outputFileName(output+".txt", ext))
}

func TestExtensionToFormat_ValidFormat(t *testing.T) {
	format := string(graphviz.PNG)
	defaultFormat := graphviz.SVG
	require.Equal(t, graphviz.PNG, extensionToFormat(format, defaultFormat))
}

func TestExtensionToFormat_EmptyString(t *testing.T) {
	format := ""
	defaultFormat := graphviz.SVG
	require.Equal(t, graphviz.SVG, extensionToFormat(format, defaultFormat))
}

func TestExtensionToFormat_InvalidFormat(t *testing.T) {
	format := "notValid"
	defaultFormat := graphviz.SVG
	require.Equal(t, graphviz.SVG, extensionToFormat(format, defaultFormat))
}
