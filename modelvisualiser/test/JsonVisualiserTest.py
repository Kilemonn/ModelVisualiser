from modelvisualiser.test.VisualiserTest import VisualiserTest
from modelvisualiser.visualiser import Constants, ModelVisualiser


class JasonVisualiserTest(VisualiserTest):
    def test_json_file_1(self):
        filepath = "json_files/test1.json"
        graph = ModelVisualiser.create_graph_from_filename(filepath)
        node_names = ["object-key", "object-array", "empty-object"]
        expected_node_count = [2, 3, 1]
        root_node_count = 8
        non_existent_node = "some-other-node"
        self.verify_graph(graph, node_names, expected_node_count, root_node_count, non_existent_node)
