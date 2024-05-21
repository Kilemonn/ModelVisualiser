import sys

from modelvisualiser.visualiser.ModelVisualiser import create_graph_from_filename, output_file_name


def main():
    if len(sys.argv) <= 1:
        print("Expected JSON file path as first program argument.")
        return

    filename = sys.argv[1]
    graph = create_graph_from_filename(filename)

    output_format = "png"
    graph.render(outfile=output_file_name(filename, output_format), format=output_format)


if __name__ == '__main__':
    main()
