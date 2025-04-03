package visualiser

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNodeCreationInSubgraphs(t *testing.T) {
	mainNodeName := "mainNode"
	subNode1Name := "node1"
	subGraph1Name := "subgraph1"
	subGraph2Name := "subgraph2"
	subNode2Name := "node2"

	v, err := NewVisualiser(context.Background())
	require.NoError(t, err)
	defer v.Close()

	g, err := v.g.Graph()
	require.NoError(t, err)

	mainNode, err := g.CreateNodeByName(mainNodeName)
	require.NoError(t, err)
	require.NotNil(t, mainNode)

	// Make sure main node is retrievable from main graph
	n, err := g.NodeByName(mainNodeName)
	require.NoError(t, err)
	require.NotNil(t, n)

	subgraph1, err := g.CreateSubGraphByName(subGraph1Name)
	require.NoError(t, err)
	require.NotNil(t, subgraph1)

	subNode1, err := subgraph1.CreateNodeByName(subNode1Name)
	require.NoError(t, err)
	require.NotNil(t, subNode1)

	// Make sure sub node is retrievable from main graph
	n, err = g.NodeByName(subNode1Name)
	require.NoError(t, err)
	require.NotNil(t, n)

	// Make sure we cannot get the main node from the subgraph
	n, err = subgraph1.NodeByName(mainNodeName)
	require.NoError(t, err)
	require.Nil(t, n)

	subgraph2, err := g.CreateSubGraphByName(subGraph2Name)
	require.NoError(t, err)
	require.NotNil(t, subgraph2)

	subNode2, err := subgraph2.CreateNodeByName(subNode2Name)
	require.NoError(t, err)
	require.NotNil(t, subNode2)

	// Make sure subgraph2 node is available in main graph
	n, err = g.NodeByName(subNode2Name)
	require.NoError(t, err)
	require.NotNil(t, n)

	// Make sure we cannot get node 2 in subgraph2 from subgraph 1
	n, err = subgraph1.NodeByName(subNode2Name)
	require.NoError(t, err)
	require.Nil(t, n)
}
