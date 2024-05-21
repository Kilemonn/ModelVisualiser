import unittest

import graphviz

import modelvisualiser.visualiser.ModelVisualiser
from modelvisualiser.visualiser import Constants


class VisualiserTest(unittest.TestCase):
    SUBGRAPH: str = "subgraph"
    EMPTY_LABEL: str = "[label=\"\"]"

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
