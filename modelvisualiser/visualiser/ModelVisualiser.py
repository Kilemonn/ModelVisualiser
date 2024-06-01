import json
from typing import Any

import graphviz
import puremagic
import yaml

from modelvisualiser.visualiser import Constants

JSON_MIME: str = "application/json"
YAML_MIME: str = "application/x-yaml"


def create_graph_from_filename(filename: str) -> graphviz.Digraph:
    graph = graphviz.Digraph(node_attr={'shape': 'record'})
    model = load_file(filename)
    create_nodes(graph, model)
    return graph


def load_file(filename: str) -> Any:
    file_type = puremagic.magic_file(filename)
    # print(file_type)
    with open(filename, "r") as file:
        if len(file_type) > 0:
            if file_type[0].mime_type == JSON_MIME:
                return json.loads(file.read())
            elif file_type[0].mime_type == YAML_MIME:
                return yaml.safe_load(file.read())
        else:
            print("Unable to determine input file's type, attempting to parse as JSON.")
            return json.loads(file.read())


def create_nodes(graph: graphviz.Digraph, model: dict, name_path: str = Constants.ROOT_PATH) -> str:
    placeholder_name = name_path + Constants.EMPTY_SUFFIX
    with graph.subgraph(name=Constants.CLUSTER_PREFIX + name_path) as subgraph:
        subgraph.attr(label=name_path)
        for key in model:
            key_type = type(model[key])
            property_path = name_path + "/" + key
            property_name = property_path + " - " + key_type.__name__
            if isinstance(model[key], list) or isinstance(model[key], dict):
                if isinstance(model[key], list):
                    handle_list(subgraph, graph, model[key], property_path, property_name)
                elif isinstance(model[key], dict):
                    handle_dict(subgraph, graph, model[key], property_path)
            else:
                subgraph.node(name=property_name, label=property_name)

        # We shouldn't blindly create an empty node in the root since nothing can link back to it
        if name_path != Constants.ROOT_PATH:
            subgraph.node(name=placeholder_name, label="")

    return placeholder_name


def handle_list(subgraph: graphviz.Digraph, graph: graphviz.Digraph, model: list, key: str, property_name: str):
    if len(model) > 0:
        list_type = type(model[0])
        if isinstance(model[0], dict):
            placeholder_property = create_nodes(graph, model[0], key)
            subgraph.edge(property_name, placeholder_property)
        else:
            property_name = property_name + "[" + list_type.__name__ + "]"
            subgraph.node(name=property_name, label=property_name)


def handle_dict(subgraph: graphviz.Digraph, graph: graphviz.Digraph, model: dict, key: str):
    placeholder_property = create_nodes(graph, model, key)
    subgraph.edge(key, placeholder_property)


def output_file_name(input_file_name: str, output_format: str) -> str:
    split = input_file_name.rsplit(Constants.PERIOD, maxsplit=1)
    if len(split) > 0:
        return split[0] + Constants.PERIOD + output_format
    return input_file_name + Constants.PERIOD + output_format
