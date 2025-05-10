package testutil

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"testing"

	"github.com/Kilemonn/ModelVisualiser/consts"
	"github.com/goccy/go-graphviz"
	"github.com/stretchr/testify/require"
)

func GetGraphCount(graph *graphviz.Graph) int {
	return len(GetSubGraphEntries(graph))
}

func GetSubGraphEntries(graph *graphviz.Graph) []string {
	subGraphNames := []string{}

	g, err := graph.GraphRoot().FirstSubGraph()
	for g != nil && err == nil {
		subGraphNames = append(subGraphNames, g.Label())
		g, err = g.NextSubGraph()
	}
	return subGraphNames
}

func SubGraphWithNameExists(graph *graphviz.Graph, subGraphName string) bool {
	subGraphEntries := GetSubGraphEntries(graph)
	for _, entry := range subGraphEntries {
		if strings.Contains(entry, subGraphName) {
			return true
		}
	}
	return false
}

func GetEmptyCount(graph *graphviz.Graph) int {
	return len(GetEmptyNodes(graph))
}

func GetEmptyNodes(graph *graphviz.Graph) []string {
	emptyNodes := []string{}

	g, _ := graph.GraphRoot().FirstSubGraph()
	for g != nil {
		n, _ := g.FirstNode()
		for n != nil {
			name, err := n.Name()
			if err != nil {
				fmt.Printf("Failed to get node name %s", err.Error())
				return emptyNodes
			}
			if strings.Contains(name, consts.EMPTY_SUFFIX) {
				emptyNodes = append(emptyNodes, name)
			}
			n, _ = g.NextNode(n)
		}
		g, _ = g.NextSubGraph()
	}

	return emptyNodes
}

func EmptyNodeWithNameExists(graph *graphviz.Graph, nodeName string) bool {
	emptyNodes := GetEmptyNodes(graph)
	for _, entry := range emptyNodes {
		if strings.Contains(entry, nodeName+consts.EMPTY_SUFFIX) {
			return true
		}
	}
	return false
}

func GetLinkageCount(graph *graphviz.Graph) int {
	return len(GetLinkages(graph))
}

func GetLinkages(graph *graphviz.Graph) []string {
	linkages := make(map[string]bool)

	g, _ := graph.GraphRoot().FirstSubGraph()
	for g != nil {
		n, _ := g.FirstNode()
		for n != nil {
			e, _ := n.Root().FirstEdge(n)
			for e != nil {
				name, err := e.Name()
				if err != nil {
					fmt.Printf("Failed to get edge name %s", err.Error())

					var slice []string
					for key := range maps.Keys(linkages) {
						slice = append(slice, key)
					}
					return slice
				}
				linkages[name] = true

				e, _ = n.Root().NextEdge(e, n)
			}
			n, _ = g.NextNode(n)
		}
		g, _ = g.NextSubGraph()
	}

	var slice []string
	for key := range maps.Keys(linkages) {
		slice = append(slice, key)
	}
	return slice
}

func LinkageExists(graph *graphviz.Graph, parentPropertyName string) bool {
	linkages := GetLinkages(graph)
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

	require.Equal(t, len(nodeNames)+1, GetGraphCount(graph))
	require.True(t, SubGraphWithNameExists(graph, consts.ROOT_PATH))

	for _, name := range nodeNames {
		require.True(t, SubGraphWithNameExists(graph, name))
	}

	// Make sure some other random subgraph name does not exist
	require.False(t, SubGraphWithNameExists(graph, nonExistentNode))

	require.Equal(t, len(nodeNames), GetEmptyCount(graph))
	for _, name := range nodeNames {
		require.True(t, EmptyNodeWithNameExists(graph, name))
	}

	// Make sure there is no empty node for root
	require.False(t, EmptyNodeWithNameExists(graph, consts.ROOT_PATH))
	require.False(t, EmptyNodeWithNameExists(graph, nonExistentNode))

	require.Equal(t, len(nodeNames), GetLinkageCount(graph))
	for _, name := range nodeNames {
		require.True(t, LinkageExists(graph, name))
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
