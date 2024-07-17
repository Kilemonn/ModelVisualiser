from modelvisualiser.test.VisualiserTest import VisualiserTest
from modelvisualiser.visualiser import ModelVisualiser, Constants


class YamlVisualiserTest(VisualiserTest):
    def test_yaml_file(self):
        filepath = "yaml_files/test.yml"
        graph = ModelVisualiser.create_graph_from_filename(filepath)

        node_names = ["services", "services/queue", "services/queue/environment", "services/queue/healthcheck", "services/hello-world", "services/hello-world/depends_on", "services/hello-world/depends_on/queue"]
        expected_node_count = [3, 5, 3, 6, 3, 2, 2]
        root_node_count = 1
        non_existent_node = "some-other-node"

        self.verify_graph(graph, node_names, expected_node_count, root_node_count, non_existent_node)
