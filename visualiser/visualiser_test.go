package visualiser

import (
	"context"
	"testing"

	"github.com/Kilemonn/ModelVisualiser/visualiser/testutil"
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

func TestFromFile_YamlTestFile(t *testing.T) {
	ctx := context.Background()
	visualiser, err := NewVisualiser(ctx)
	require.NoError(t, err)

	graph, err := visualiser.FromFile("test_files/yaml_files/test.yml")
	require.NoError(t, err)
	defer graph.Close()

	nodeNames := []string{"services", "services/queue", "services/queue/environment", "services/queue/healthcheck", "services/hello-world", "services/hello-world/depends_on", "services/hello-world/depends_on/queue"}
	expectedNodeCount := []int{3, 5, 3, 6, 3, 2, 2}
	rootNodeCount := 1
	nonExistentNode := "some-other-node"

	testutil.VerifyGraph(t, graph, nodeNames, expectedNodeCount, rootNodeCount, nonExistentNode)
}

func TestFromFile_JsonFile(t *testing.T) {
	ctx := context.Background()
	visualiser, err := NewVisualiser(ctx)
	require.NoError(t, err)

	graph, err := visualiser.FromFile("test_files/json_files/test1.json")
	require.NoError(t, err)
	defer graph.Close()

	nodeNames := []string{"object-key", "object-array", "empty-object"}
	expectedNodeCount := []int{2, 3, 1}
	rootNodeCount := 9
	nonExistentNode := "some-other-node"

	testutil.VerifyGraph(t, graph, nodeNames, expectedNodeCount, rootNodeCount, nonExistentNode)
}
