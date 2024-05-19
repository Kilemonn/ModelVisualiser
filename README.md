# ModelVisualiser
A visualiser tool that is able to visualise complex JSON, YAML and XML model objects.

# Getting Started

## Install Pipenv

Install pipenv:
> pip install pipenv

Call pipenv to install required dependencies:
> pipenv install

Install new dependency:
> pipenv install \<package-name\>

## Locking and Updating pipe

`pipenv lock` — records the new requirements to the Pipfile.lock file.

`pipenv update` — records the new requirements to the Pipfile.lock file and installs the missing dependencies on the Python interpreter.

# Usage

The program can be run using the `python3` interpreter. The application is expecting a single commandline argument which is the path to the input file that you wish to visualise.

> python3 main.py path/to/file.json

The output is always a `.png` file for now.
