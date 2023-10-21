import json

import graphviz


def main():
    graph = graphviz.Digraph()
    file_name = "test.json"
    with open(file_name, 'r') as file:
        model = json.loads(file.read())
        create_nodes(graph, model)


def create_nodes(graph: graphviz.Digraph, model: dict):
    for key in model:
        print(key, " type: ", type(model[key]))
        if isinstance(model[key], dict):
            create_nodes(graph, model[key])
        # graph.node('')


if __name__ == '__main__':
    main()
