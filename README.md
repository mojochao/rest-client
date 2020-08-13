# rest-client

The `rest-client` application is a command-line interface application that
is capable of executing JetBrains REST client `*.http` or `*.rest` request
files against multiple environments defined in `rest-client.env.json` 
environment files.

## Installation

    go get github.com/mojochao/rest-client

## Usage

The `envs` command lists environments in the environment file.

    $ rest-client [-e ENV] [-p PATH...] envs  

The `reqs` command lists requests defined in request files.

    $ rest-client [-e ENV] [-p PATH...] reqs

The `exec` command executes requests defined in request files with an environment
defined in an environment file.

    $ rest-client [-e ENV] [-p PATH...] exec [REQUESTS...] 

If `-e` option is not provided, defaults to `rest-client.end.json` file in the
current working directory.

The `-p` option provides file or directory paths to process.  If not provided,
defaults to all `*.http` files in the current working directory.

The `REQUESTS` args provides one or more request names to execute.
