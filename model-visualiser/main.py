import json

import graphviz


def main():
    output_format = "png"
    file_name = "test.json"

    graph = graphviz.Digraph(node_attr={'shape': 'record'})
    with open(file_name, 'r') as file:
        model = json.loads(file.read())
        create_nodes(graph, model)
    graph.render(outfile=output_file_name(file_name, output_format), format=output_format)


def create_nodes(graph: graphviz.Digraph, model: dict, name_path: str = "root") -> str:
    placeholder_name = name_path + "_empty"
    with graph.subgraph(name="cluster_" + name_path) as g:
        g.attr(label=name_path)
        for key in model:
            key_type = type(model[key])
            property_name = name_path + "/" + key + " - " + key_type.__name__
            if isinstance(model[key], list):
                property_name = handle_list(g, graph, model[key], name_path + "/" + key, property_name)

            if isinstance(model[key], dict):
                handle_dict(g, graph, model[key], name_path + "/" + key)
            else:
                g.node(name=property_name, label=property_name, shape="record")
        g.node(name=placeholder_name, label="")
    return placeholder_name


def handle_list(subgraph: graphviz.Digraph, graph: graphviz.Digraph, model: list, key: str, property_name: str) -> str:
    if len(model) > 0:
        list_type = type(model[0])
        if isinstance(model[0], dict):
            r = create_nodes(graph, model[0], key)
            if r is None:
                # empty?
                pass
            else:
                subgraph.edge(key, r)
        else:
            property_name = property_name + "[" + list_type.__name__ + "]"
    return property_name


def handle_dict(subgraph: graphviz.Digraph, graph: graphviz.Digraph, model: dict, key: str):
    r = create_nodes(graph, model, key)
    if r is None:
        # empty object?
        pass
    else:
        subgraph.edge(key, r)


def output_file_name(input_file_name: str, output_format: str) -> str:
    PERIOD = "."
    split = input_file_name.split(PERIOD)
    if len(split) > 0:
        return split[0] + PERIOD + output_format
    return input_file_name + PERIOD + output_format


if __name__ == '__main__':
    main()
