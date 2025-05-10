package testutil

import (
	"maps"
	"slices"
	"strings"
	"testing"

	"github.com/Kilemonn/ModelVisualiser/consts"
	"github.com/goccy/go-graphviz"
	"github.com/stretchr/testify/require"
)

func GetGraphCount(t *testing.T, graph *graphviz.Graph) int {
	return len(GetSubGraphEntries(t, graph))
}

func GetSubGraphEntries(t *testing.T, graph *graphviz.Graph) []string {
	subGraphNames := []string{}

	g, err := graph.GraphRoot().FirstSubGraph()
	require.NoError(t, err)
	for g != nil && err == nil {
		subGraphNames = append(subGraphNames, g.Label())
		g, err = g.NextSubGraph()
		require.NoError(t, err)
	}
	return subGraphNames
}

func SubGraphWithNameExists(t *testing.T, graph *graphviz.Graph, subGraphName string) bool {
	subGraphEntries := GetSubGraphEntries(t, graph)
	for _, entry := range subGraphEntries {
		if strings.Contains(entry, subGraphName) {
			return true
		}
	}
	return false
}

func GetEmptyCount(t *testing.T, graph *graphviz.Graph) int {
	return len(GetEmptyNodes(t, graph))
}

func GetEmptyNodes(t *testing.T, graph *graphviz.Graph) []string {
	emptyNodes := []string{}

	g, err := graph.GraphRoot().FirstSubGraph()
	require.NoError(t, err)
	for g != nil {
		n, err := g.FirstNode()
		require.NoError(t, err)
		for n != nil {
			name, err := n.Name()
			require.NoError(t, err)
			if strings.Contains(name, consts.EMPTY_SUFFIX) {
				emptyNodes = append(emptyNodes, name)
			}
			n, err = g.NextNode(n)
			require.NoError(t, err)
		}
		g, err = g.NextSubGraph()
		require.NoError(t, err)
	}

	return emptyNodes
}

func EmptyNodeWithNameExists(t *testing.T, graph *graphviz.Graph, nodeName string) bool {
	emptyNodes := GetEmptyNodes(t, graph)
	for _, entry := range emptyNodes {
		if strings.Contains(entry, nodeName+consts.EMPTY_SUFFIX) {
			return true
		}
	}
	return false
}

func GetLinkageCount(t *testing.T, graph *graphviz.Graph) int {
	return len(GetLinkages(t, graph))
}

func GetLinkages(t *testing.T, graph *graphviz.Graph) []string {
	linkages := make(map[string]bool)

	g, err := graph.GraphRoot().FirstSubGraph()
	require.NoError(t, err)
	for g != nil {
		n, err := g.FirstNode()
		require.NoError(t, err)
		for n != nil {
			e, err := n.Root().FirstEdge(n)
			require.NoError(t, err)
			for e != nil {
				name, err := e.Name()
				require.NoError(t, err)
				linkages[name] = true

				e, err = n.Root().NextEdge(e, n)
				require.NoError(t, err)
			}
			n, err = g.NextNode(n)
			require.NoError(t, err)
		}
		g, err = g.NextSubGraph()
		require.NoError(t, err)
	}

	var slice []string
	for key := range maps.Keys(linkages) {
		slice = append(slice, key)
	}
	return slice
}

func LinkageExists(t *testing.T, graph *graphviz.Graph, parentPropertyName string) bool {
	linkages := GetLinkages(t, graph)
	for _, l := range linkages {
		split := strings.Split(l, consts.ARROW)
		if len(split) >= 2 && strings.Contains(split[0], parentPropertyName) && strings.Contains(split[1], parentPropertyName) && strings.Contains(split[1], parentPropertyName+consts.EMPTY_SUFFIX) {
			return true
		}
	}
	return false
}

func GetNodesInSubGraph(t *testing.T, graph *graphviz.Graph, subgraphName string) []string {
	nodeNames := []string{}
	subgraph, err := graph.SubGraphByName(subgraphName)
	require.NoError(t, err)

	n, err := subgraph.FirstNode()
	require.NoError(t, err)
	for n != nil {
		name, err := n.Name()
		require.NoError(t, err)
		nodeNames = append(nodeNames, name)

		n, err = subgraph.NextNode(n)
		require.NoError(t, err)
	}

	return nodeNames
}

func VerifyGraph(t *testing.T, graph *graphviz.Graph, nodeNames []string, expectedNodeCounts []int, rootNodeCount int, nonExistentNode string) {
	require.False(t, slices.Contains(nodeNames, nonExistentNode))
	require.Equal(t, len(nodeNames), len(expectedNodeCounts))

	require.Equal(t, len(nodeNames)+1, GetGraphCount(t, graph))
	require.True(t, SubGraphWithNameExists(t, graph, consts.ROOT_PATH))

	for _, name := range nodeNames {
		require.True(t, SubGraphWithNameExists(t, graph, name))
	}

	// Make sure some other random subgraph name does not exist
	require.False(t, SubGraphWithNameExists(t, graph, nonExistentNode))

	require.Equal(t, len(nodeNames), GetEmptyCount(t, graph))
	for _, name := range nodeNames {
		require.True(t, EmptyNodeWithNameExists(t, graph, name))
	}

	// Make sure there is no empty node for root
	require.False(t, EmptyNodeWithNameExists(t, graph, consts.ROOT_PATH))
	require.False(t, EmptyNodeWithNameExists(t, graph, nonExistentNode))

	require.Equal(t, len(nodeNames), GetLinkageCount(t, graph))
	for _, name := range nodeNames {
		require.True(t, LinkageExists(t, graph, name))
	}

	rootNodes := GetNodesInSubGraph(t, graph, consts.CLUSTER_PREFIX+consts.ROOT_PATH)
	require.Equal(t, len(rootNodes), rootNodeCount)

	for i := 0; i < len(nodeNames); i++ {
		nodeName := nodeNames[i]
		nodeCount := expectedNodeCounts[i]
		nodes := GetNodesInSubGraph(t, graph, consts.CLUSTER_PREFIX+consts.ROOT_PATH+"/"+nodeName)
		require.Equal(t, nodeCount, len(nodes))
	}
}
