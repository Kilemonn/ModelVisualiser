package model_visualiser

import (
	"context"
	"os"
	"reflect"

	"github.com/Kilemonn/ModelVisualiser/consts"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"gopkg.in/yaml.v3"
)

type ModelVisualiser struct {
	g   *graphviz.Graphviz
	ctx context.Context
}

func NewModelVisualiser(ctx context.Context) (ModelVisualiser, error) {
	graph, err := graphviz.New(ctx)
	if err != nil {
		return ModelVisualiser{
			g:   graph,
			ctx: ctx,
		}, nil
	}
	return ModelVisualiser{}, err
}

func (mv ModelVisualiser) FromFile(fileName string) (*graphviz.Graph, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	data := make(map[any]any)

	err = yaml.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}

	// err = json.Unmarshal(content, &data)
	// if err != nil {
	// 	return nil, err
	// }

	return mv.createGraph(data)
}

func (mv ModelVisualiser) createGraph(data map[any]any) (*graphviz.Graph, error) {
	g, err := mv.g.Graph()
	if err != nil {
		return nil, err
	}

	_, err = mv.createNodes(g, data, consts.ROOT_PATH)
	if err != nil {
		return nil, err
	}

	return g, err
}

func (mv ModelVisualiser) createNodes(graph *graphviz.Graph, data map[any]any, namePath string) (string, error) {
	placeHolderName := namePath + consts.EMPTY_SUFFIX
	subgraph, err := graph.SubGraphByName(consts.CLUSTER_PREFIX + namePath)
	if err != nil {
		return "", err
	}

	// subgraph.Attr()
	for _, k := range data {
		propertyPath := namePath + "/" + k.(string)
		propertyName := propertyPath + " - " + string(reflect.TypeOf(k).String())

		valueType := reflect.TypeOf(data[k])

		if valueType.Kind() == reflect.Array || valueType.Kind() == reflect.Slice {
			mv.handleList(subgraph, graph, data[k].([]any), k.(string), propertyName)
		} else if valueType.Kind() == reflect.Map {
			mv.handleMap(subgraph, graph, data[k].(map[any]any), k.(string))
		} else {
			_, err := subgraph.CreateNodeByName(propertyName)
			if err != nil {
				return "", err
			}
		}
	}

	if namePath != consts.ROOT_PATH {
		_, err := subgraph.CreateNodeByName(placeHolderName)
		if err != nil {
			return "", err
		}
	}

	return placeHolderName, nil
}

func (mv ModelVisualiser) handleList(subgraph *cgraph.Graph, graph *graphviz.Graph, model []any, key string, propertyName string) error {
	if len(model) == 0 {
		return nil
	}

	listType := reflect.TypeOf(model[0])
	if listType.Kind() == reflect.Map {
		return mv.handleMap(subgraph, graph, model[0].(map[any]any), key)
	} else {
		propertyName = propertyName + "[" + listType.String() + "]"
		_, err := subgraph.CreateNodeByName(propertyName)
		return err
	}
}

func (mv ModelVisualiser) handleMap(subgraph *cgraph.Graph, graph *graphviz.Graph, model map[any]any, key string) error {
	placeHolderName, err := mv.createNodes(graph, model[0].(map[any]any), key)
	if err != nil {
		return err
	}
	return mv.createSubgraphEdgeBetweenNodesByName(subgraph, "", key, placeHolderName)
}

func (mv ModelVisualiser) createSubgraphEdgeBetweenNodesByName(subgraph *cgraph.Graph, edgeName string, node1Name string, node2Name string) error {
	node1, err := subgraph.NodeByName(node1Name)
	if err != nil {
		return err
	}

	node2, err := subgraph.NodeByName(node2Name)
	if err != nil {
		return err
	}

	_, err = subgraph.CreateEdgeByName(edgeName, node1, node2)
	return err
}

func (mv ModelVisualiser) createGraphEdgeBetweenNodesByName(graph *graphviz.Graph, edgeName string, node1Name string, node2Name string) error {
	node1, err := graph.NodeByName(node1Name)
	if err != nil {
		return err
	}

	node2, err := graph.NodeByName(node2Name)
	if err != nil {
		return err
	}

	_, err = graph.CreateEdgeByName(edgeName, node1, node2)
	return err
}

func (mv ModelVisualiser) ToFile(graph *graphviz.Graph, outputFile string) error {
	return mv.g.RenderFilename(mv.ctx, graph, graphviz.PNG, outputFile)
}

func (mv ModelVisualiser) Close() error {
	if mv.g != nil {
		return mv.g.Close()
	}
	return nil
}
