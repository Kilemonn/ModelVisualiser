# ModelVisualiser
A visualiser tool that is able to visualise complex JSON, YAML and XML model objects.

## Installation 

> go install github.com/Kilemonn/ModelVisualiser@latest

## Usage 

The application requires 1 argument to run this is `-i` which specifies the input file path.
The generated graph will be output in the same directory as the input file, with its file extension updated to appropriate file format it is being output to.

E.g. (Generating with output format set to jpg)
```bash
ModelVisualiser -i myfile.yaml -f jpg
# The graph will be generated and output as "myfile.jpg"
```

You can get application usage by using `-h` or `--help`.
