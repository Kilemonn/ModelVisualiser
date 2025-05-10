package visualiser

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/Kilemonn/ModelVisualiser/consts"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"gopkg.in/yaml.v3"
)

type Visualiser struct {
	g   *graphviz.Graphviz
	ctx context.Context
}

func NewVisualiser(ctx context.Context) (Visualiser, error) {
	graph, err := graphviz.New(ctx)
	if err != nil {
		return Visualiser{}, err
	}

	return Visualiser{
		g:   graph,
		ctx: ctx,
	}, nil
}

func (mv Visualiser) FromFile(fileName string) (*graphviz.Graph, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	data := make(map[string]any)

	// Just attempt to parse both json and yaml
	err = json.Unmarshal(content, &data)
	if err != nil {
		err = yaml.Unmarshal(content, &data)
		if err != nil {
			return nil, err
		}
	}

	return mv.createGraph(data)
}

func (mv Visualiser) createGraph(data map[string]any) (*graphviz.Graph, error) {
	// TODO: node_attr={'shape': 'record'}
	g, err := mv.g.Graph()
	if err != nil {
		return nil, err
	}

	_, err = mv.createNodes(g, data, consts.ROOT_PATH)
	if err != nil {
		return nil, err
	}

	err = mv.createEdges(g, data, consts.ROOT_PATH)

	return g, err
}

func (mv Visualiser) createNodes(graph *graphviz.Graph, data map[string]any, namePath string) (string, error) {
	placeHolderName := namePath + consts.EMPTY_SUFFIX
	subgraph, err := graph.CreateSubGraphByName(consts.CLUSTER_PREFIX + namePath)
	if err != nil {
		return "", err
	}

	subgraph = subgraph.SetLabel(namePath)
	for k, v := range data {
		propertyPath := namePath + "/" + k
		propertyName := propertyPath + " - " + string(reflect.TypeOf(k).String())
		fmt.Printf("Key: %s\nProperty name: %s\n", k, propertyName)

		valueType := reflect.TypeOf(v)
		if valueType.Kind() == reflect.Array || valueType.Kind() == reflect.Slice {
			mv.handleListNode(subgraph, graph, v.([]any), propertyPath, propertyPath+" - list")
		} else if valueType.Kind() == reflect.Map {
			err := mv.createNodeInSubGraph(subgraph, propertyPath, propertyName)
			if err != nil {
				return "", err
			}
			mv.handleMapNode(graph, v.(map[string]any), propertyPath)
		} else {
			err := mv.createNodeInSubGraph(subgraph, propertyPath, propertyName)
			if err != nil {
				return "", err
			}
		}
	}

	if namePath != consts.ROOT_PATH {
		err := mv.createNodeInSubGraph(subgraph, placeHolderName, "")
		if err != nil {
			return "", err
		}
	}

	return placeHolderName, nil
}

func (mv Visualiser) createEdges(graph *graphviz.Graph, data map[string]any, namePath string) error {
	subgraphName := consts.CLUSTER_PREFIX + namePath
	subgraph, err := graph.SubGraphByName(subgraphName)
	if err != nil || subgraph == nil {
		return fmt.Errorf("subgraph with name [%s] does not exist", subgraphName)
	}

	for k, v := range data {
		valueType := reflect.TypeOf(v)
		if valueType.Kind() == reflect.Map {
			mv.handleMapEdge(graph, v.(map[string]any), namePath+"/"+k)
		} else if valueType.Kind() == reflect.Array || valueType.Kind() == reflect.Slice {
			mv.handleListEdge(graph, v.([]any), namePath+"/"+k)
		}
	}

	return nil
}

func (mv Visualiser) createNodeInSubGraph(subgraph *cgraph.Graph, nodeName string, nodeLabel string) error {
	n, err := subgraph.CreateNodeByName(nodeName)
	if err == nil && n != nil {
		n.SetShape(cgraph.BoxShape)
		n.SetLabel(nodeLabel)
	}
	return err
}

func (mv Visualiser) handleListNode(subgraph *cgraph.Graph, graph *graphviz.Graph, model []any, key string, propertyName string) error {
	if len(model) == 0 {
		return nil
	}

	listType := reflect.TypeOf(model[0])
	if listType.Kind() == reflect.Map {
		err := mv.createNodeInSubGraph(subgraph, key, propertyName)
		if err != nil {
			return err
		}
		return mv.handleMapNode(graph, model[0].(map[string]any), key)
	} else {
		propertyName = propertyName + "[" + listType.String() + "]"
		return mv.createNodeInSubGraph(subgraph, propertyName, propertyName)
	}
}

func (mv Visualiser) handleMapNode(graph *graphviz.Graph, model map[string]any, key string) error {
	_, err := mv.createNodes(graph, model, key)
	return err
}

func (mv Visualiser) handleMapEdge(graph *graphviz.Graph, model map[string]any, key string) error {
	placeHolderName := key + consts.EMPTY_SUFFIX
	err := mv.createEdges(graph, model, key)
	if err != nil {
		return err
	}
	return mv.createSubgraphEdgeBetweenNodesByName(graph, key+consts.ARROW+placeHolderName, key, placeHolderName)
}

func (mv Visualiser) handleListEdge(graph *graphviz.Graph, model []any, key string) error {
	if len(model) == 0 {
		return nil
	}

	placeHolderName := key + consts.EMPTY_SUFFIX
	listType := reflect.TypeOf(model[0])
	if listType.Kind() == reflect.Map {
		// If its a list of a map type, continue the recursion
		err := mv.createSubgraphEdgeBetweenNodesByName(graph, key+consts.ARROW+placeHolderName, key, placeHolderName)
		if err != nil {
			return err
		}
		return mv.createEdges(graph, model[0].(map[string]any), key)
	}

	// If its a list of a non-complex type, then no linkages need to be made

	return nil
}

func (mv Visualiser) createSubgraphEdgeBetweenNodesByName(graph *graphviz.Graph, edgeName string, node1Name string, node2Name string) error {
	node1, _ := graph.NodeByName(node1Name)
	if node1 == nil {
		return fmt.Errorf("failed to create edge [%s] because node [%s] does not exist", edgeName, node1Name)
	}

	node2, _ := graph.NodeByName(node2Name)
	if node2 == nil {
		return fmt.Errorf("failed to create edge [%s] because node [%s] does not exist", edgeName, node2Name)
	}

	_, err := graph.CreateEdgeByName(edgeName, node1, node2)
	return err
}

func (mv Visualiser) ToFile(graph *graphviz.Graph, outputFile string, outputFormat graphviz.Format) error {
	return mv.g.RenderFilename(mv.ctx, graph, outputFormat, outputFile)
}

func (mv Visualiser) Close() error {
	return mv.g.Close()
}
