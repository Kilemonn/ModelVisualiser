import unittest

import graphviz

import modelvisualiser.visualiser.ModelVisualiser
from modelvisualiser.visualiser import Constants


class VisualiserTest(unittest.TestCase):
    SUBGRAPH: str = "subgraph"
    EMPTY_LABEL: str = "[label=\"\"]"
    ARROW: str = " -> "

    def testTest1Json(self):
        filepath = "json_files/test1.json"
        graph = modelvisualiser.visualiser.ModelVisualiser.create_graph_from_filename(filepath)
        self.assertEqual(4, self.get_subgraph_count(graph))

        self.assertTrue(self.subgraph_with_name_exists(graph, "root"))
        self.assertTrue(self.subgraph_with_name_exists(graph, "object-key"))
        self.assertTrue(self.subgraph_with_name_exists(graph, "object-array"))
        self.assertTrue(self.subgraph_with_name_exists(graph, "empty-object"))
        # Make sure some other random subgraph name does not exist
        self.assertFalse(self.subgraph_with_name_exists(graph, "some-other-object"))

        self.assertEqual(3, self.get_empty_count(graph))
        self.assertTrue(self.empty_with_name_exists(graph, "object-key"))
        self.assertTrue(self.empty_with_name_exists(graph, "object-array"))
        self.assertTrue(self.empty_with_name_exists(graph, "empty-object"))
        # Make sure there is no empty node for root
        self.assertFalse(self.empty_with_name_exists(graph, "root"))

        self.assertEqual(3, self.get_linkage_count(graph))
        self.assertTrue(self.linkage_exists(graph, "object-key"))
        self.assertTrue(self.linkage_exists(graph, "object-array"))
        self.assertTrue(self.linkage_exists(graph, "empty-object"))

        nodes = self.get_nodes_in_subgraph(graph, Constants.CLUSTER_PREFIX + "root")
        self.assertEqual(7, len(nodes))

        nodes = self.get_nodes_in_subgraph(graph, f"\"{Constants.CLUSTER_PREFIX}root/object-key\"")
        self.assertEqual(2, len(nodes))

        nodes = self.get_nodes_in_subgraph(graph, f"\"{Constants.CLUSTER_PREFIX}root/object-array\"")
        self.assertEqual(3, len(nodes))

        nodes = self.get_nodes_in_subgraph(graph, f"\"{Constants.CLUSTER_PREFIX}root/empty-object\"")
        self.assertEqual(1, len(nodes))

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
