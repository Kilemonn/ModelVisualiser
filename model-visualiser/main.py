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
                l: list = model[key]
                if len(l) > 0:
                    list_type = type(l[0])
                    if isinstance(l[0], dict):
                        r = create_nodes(graph, l[0], key)
                        if r is None:
                            # empty?
                            pass
                        else:
                            g.edge(key, r)
                    property_name = key + " - " + key_type.__name__ + "[" + list_type.__name__ + "]"

            if ret is None:
                ret = property_name
            g.node(name=property_name, label=property_name, shape="record")
            if isinstance(model[key], dict):
                r = create_nodes(graph, model[key], key)
                if r is None:
                    # empty object?
                    pass
                else:
                    g.edge(key, r)

        if ret is None:
            ret = name + "_empty"
            g.node(name=ret, label="")
    return ret


if __name__ == '__main__':
    main()
