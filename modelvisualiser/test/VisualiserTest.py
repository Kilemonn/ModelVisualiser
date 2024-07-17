import unittest

import graphviz

from modelvisualiser.visualiser import Constants


class VisualiserTest(unittest.TestCase):
    SUBGRAPH: str = "subgraph"
    EMPTY_LABEL: str = "[label=\"\"]"
    ARROW: str = " -> "

    def get_subgraph_count(self, graph: graphviz.Digraph) -> int:
        return len(self.get_subgraph_entries(graph))

    def get_subgraph_entries(self, graph: graphviz.Digraph) -> list[str]:
        return [entry for entry in graph.body if VisualiserTest.SUBGRAPH in entry]

    def subgraph_with_name_exists(self, graph: graphviz.Digraph, subgraph_name: str) -> bool:
        subgraph_entries = self.get_subgraph_entries(graph)
        occurrences = [subgraph_entry for subgraph_entry in subgraph_entries if Constants.CLUSTER_PREFIX in subgraph_entry and subgraph_name in subgraph_entry]
        return len(occurrences) > 0

    def get_empty_entries(self, graph: graphviz.Digraph) -> list[str]:
        return [entry for entry in graph.body if VisualiserTest.EMPTY_LABEL in entry and Constants.EMPTY_SUFFIX in entry]

    def get_empty_count(self, graph: graphviz.Digraph) -> int:
        return len(self.get_empty_entries(graph))

    def empty_with_name_exists(self, graph: graphviz.Digraph, empty_name: str) -> bool:
        empty_entries = self.get_empty_entries(graph)
        occurrences = [empty_entry for empty_entry in empty_entries if empty_name + Constants.EMPTY_SUFFIX in empty_entry]
        return len(occurrences) > 0

    def get_linkages(self, graph: graphviz.Digraph) -> list[str]:
        return [entry for entry in graph.body if VisualiserTest.ARROW in entry]

    def get_linkage_count(self, graph: graphviz.Digraph) -> int:
        return len(self.get_linkages(graph))

    def linkage_exists(self, graph: graphviz.Digraph, parent_property_name: str) -> bool:
        linkages = self.get_linkages(graph)
        for linkage in linkages:
            split = linkage.split(VisualiserTest.ARROW)
            if len(split) >= 2 and parent_property_name in split[0] and parent_property_name in split[1] and parent_property_name + Constants.EMPTY_SUFFIX in split[1]:
                return True
        return False

    def get_nodes_in_subgraph(self, graph: graphviz.Digraph, subgraph_name: str) -> list[str]:
        index = self.get_index_of_subgraph_start(graph, subgraph_name)
        if index == -1:
            return []

        result = []
        num_of_tabs = graph.body[index].count("\t")
        for i in range(index + 1, len(graph.body)):
            if graph.body[i].count("\t") > num_of_tabs:
                if not graph.body[i].strip().startswith("label="):
                    result.append(graph.body[i])
            else:
                break
        return result

    def get_index_of_subgraph_start(self, graph: graphviz.Digraph, subgraph_name: str) -> int:
        for i, entry in enumerate(graph.body):
            if VisualiserTest.SUBGRAPH in entry and subgraph_name + " {" in entry:
                return i
        return -1

    def verify_graph(self, graph: graphviz.Digraph, node_names: list[str], expected_node_counts: list[int], root_node_count: int, non_existent_node: str):
        self.assertFalse(non_existent_node in node_names)
        self.assertEqual(len(node_names), len(expected_node_counts))

        self.assertEqual(len(node_names) + 1, self.get_subgraph_count(graph))
        self.assertTrue(self.subgraph_with_name_exists(graph, "root"))
        for node_name in node_names:
            self.assertTrue(self.subgraph_with_name_exists(graph, node_name))

        # Make sure some other random subgraph name does not exist
        self.assertFalse(self.subgraph_with_name_exists(graph, non_existent_node))

        self.assertEqual(len(node_names), self.get_empty_count(graph))
        for node_name in node_names:
            self.assertTrue(self.empty_with_name_exists(graph, node_name))

        # Make sure there is no empty node for root
        self.assertFalse(self.empty_with_name_exists(graph, "root"))
        self.assertFalse(self.empty_with_name_exists(graph, non_existent_node))

        self.assertEqual(len(node_names), self.get_linkage_count(graph))
        for node_name in node_names:
            self.assertTrue(self.linkage_exists(graph, node_name))

        nodes = self.get_nodes_in_subgraph(graph, Constants.CLUSTER_PREFIX + "root")
        self.assertEqual(root_node_count, len(nodes))

        for node_name, node_count in zip(node_names, expected_node_counts):
            nodes = self.get_nodes_in_subgraph(graph, f"\"{Constants.CLUSTER_PREFIX}root/{node_name}\"")
            self.assertEqual(node_count, len(nodes))
