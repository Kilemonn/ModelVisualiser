package main

import (
	"testing"

	"github.com/Kilemonn/ModelVisualiser/consts"
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
