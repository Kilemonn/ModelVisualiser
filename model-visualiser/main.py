import json

import graphviz


def main():
    output_format = "png"
    graph = graphviz.Digraph(node_attr={'shape': 'record'})
    file_name = "test.json"
    with open(file_name, 'r') as file:
        model = json.loads(file.read())
        create_nodes(graph, model)

    graph.render(outfile=file_name.split(".")[0] + "." + output_format, format=output_format)


def create_nodes(graph: graphviz.Digraph, model: dict, name: str = "root") -> str:
    ret = None
    with graph.subgraph(name="cluster_" + name) as g:
        g.attr(label=name)
        for key in model:
            key_type = type(model[key])
            property_name = key + " - " + key_type.__name__
            if isinstance(model[key], list):
                property_name = handle_list(g, graph, model[key], key, property_name)

            if ret is None:
                ret = property_name
            g.node(name=property_name, label=property_name, shape="record")
            if isinstance(model[key], dict):
                handle_dict(g, graph, model[key], key)

        if ret is None:
            ret = name + "_empty"
            g.node(name=ret, label="")
    return ret


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


if __name__ == '__main__':
    main()
