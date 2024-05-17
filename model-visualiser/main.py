import json
import sys

import graphviz

ROOT_PATH: str = "root"
PERIOD = "."
EMPTY_SUFFIX: str = "_empty"
CLUSTER_PREFIX: str = "cluster_"


def main():
    if len(sys.argv) <= 0:
        print("Expected JSON file path as first program argument.")
        return

    output_format = "png"
    file_name = sys.argv[1]
    graph = graphviz.Digraph()
    with open(file_name, "r") as file:
        model = json.loads(file.read())
        create_nodes(graph, model)
    graph.render(outfile=output_file_name(file_name, output_format), format=output_format)


def create_nodes(graph: graphviz.Digraph, model: dict, name_path: str = ROOT_PATH) -> str:
    placeholder_name = name_path + EMPTY_SUFFIX
    with graph.subgraph(name=CLUSTER_PREFIX + name_path) as g:
        g.attr(label=name_path)
        for key in model:
            key_type = type(model[key])
            property_path = name_path + "/" + key
            property_name = property_path + " - " + key_type.__name__
            if isinstance(model[key], list):
                property_name = handle_list(g, graph, model[key], property_path, property_name)
            if isinstance(model[key], dict):
                handle_dict(g, graph, model[key], property_path)
            else:
                g.node(name=property_name, label=property_name)

        g.node(name=placeholder_name, label="")

    return placeholder_name


def handle_list(subgraph: graphviz.Digraph, graph: graphviz.Digraph, model: list, key: str, property_name: str) -> str:
    if len(model) > 0:
        list_type = type(model[0])
        if isinstance(model[0], dict):
            r = create_nodes(graph, model[0], key)
            subgraph.edge(key, r)
        else:
            property_name = property_name + "[" + list_type.__name__ + "]"
    return property_name


def handle_dict(subgraph: graphviz.Digraph, graph: graphviz.Digraph, model: dict, key: str):
    r = create_nodes(graph, model, key)
    subgraph.edge(key, r)


def output_file_name(input_file_name: str, output_format: str) -> str:
    split = input_file_name.split(PERIOD)
    if len(split) > 0:
        return split[0] + PERIOD + output_format
    return input_file_name + PERIOD + output_format


if __name__ == '__main__':
    main()
