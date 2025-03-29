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

	return g, err
}

func (mv Visualiser) createNodes(graph *graphviz.Graph, data map[string]any, namePath string) (string, error) {
	placeHolderName := namePath + consts.EMPTY_SUFFIX
	subgraph, err := graph.CreateSubGraphByName(consts.CLUSTER_PREFIX + namePath)
	if err != nil {
		return "", err
	}

	// TODO
	// subgraph.Attr("name", namePath)
	for k, v := range data {
		propertyPath := namePath + "/" + k
		propertyName := propertyPath + " - " + string(reflect.TypeOf(k).String())
		fmt.Printf("Key: %s\nProperty name: %s\n", k, propertyName)

		valueType := reflect.TypeOf(v)
		if valueType.Kind() == reflect.Array || valueType.Kind() == reflect.Slice {
			mv.handleList(subgraph, graph, v.([]any), k, propertyName)
		} else if valueType.Kind() == reflect.Map {
			mv.handleMap(subgraph, graph, v.(map[string]any), k)
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

func (mv Visualiser) handleList(subgraph *cgraph.Graph, graph *graphviz.Graph, model []any, key string, propertyName string) error {
	if len(model) == 0 {
		return nil
	}

	listType := reflect.TypeOf(model[0])
	if listType.Kind() == reflect.Map {
		return mv.handleMap(subgraph, graph, model[0].(map[string]any), key)
	} else {
		propertyName = propertyName + "[" + listType.String() + "]"
		_, err := subgraph.CreateNodeByName(propertyName)
		return err
	}
}

func (mv Visualiser) handleMap(subgraph *cgraph.Graph, graph *graphviz.Graph, model map[string]any, key string) error {
	placeHolderName, err := mv.createNodes(graph, model, key)
	if err != nil {
		return err
	}
	return mv.createSubgraphEdgeBetweenNodesByName(subgraph, "", key, placeHolderName)
}

func (mv Visualiser) createSubgraphEdgeBetweenNodesByName(subgraph *cgraph.Graph, edgeName string, node1Name string, node2Name string) error {
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

func (mv Visualiser) ToFile(graph *graphviz.Graph, outputFile string, outputFormat graphviz.Format) error {
	return mv.g.RenderFilename(mv.ctx, graph, outputFormat, outputFile)
}

func (mv Visualiser) Close() error {
	return mv.g.Close()
}
