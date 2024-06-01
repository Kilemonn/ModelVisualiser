from modelvisualiser.test.VisualiserTest import VisualiserTest
from modelvisualiser.visualiser import Constants, ModelVisualiser


class JasonVisualiserTest(VisualiserTest):
    def test_json_file_1(self):
        filepath = "json_files/test1.json"
        graph = ModelVisualiser.create_graph_from_filename(filepath)
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
        self.assertEqual(8, len(nodes))

        nodes = self.get_nodes_in_subgraph(graph, f"\"{Constants.CLUSTER_PREFIX}root/object-key\"")
        self.assertEqual(2, len(nodes))

        nodes = self.get_nodes_in_subgraph(graph, f"\"{Constants.CLUSTER_PREFIX}root/object-array\"")
        self.assertEqual(3, len(nodes))

        nodes = self.get_nodes_in_subgraph(graph, f"\"{Constants.CLUSTER_PREFIX}root/empty-object\"")
        self.assertEqual(1, len(nodes))
