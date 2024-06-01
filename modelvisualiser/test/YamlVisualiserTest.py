from modelvisualiser.test.VisualiserTest import VisualiserTest
from modelvisualiser.visualiser import ModelVisualiser


class YamlVisualiserTest(VisualiserTest):
    def test_yaml_file(self):
        filepath = "yaml_files/test.yml"
        graph = ModelVisualiser.create_graph_from_filename(filepath)
        self.assertEqual(4, self.get_subgraph_count(graph))
