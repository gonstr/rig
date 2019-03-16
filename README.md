## rig

Command line tool for managing rig templates.

```
A code-gen tool for k8s manifests.

Manage your manifests in versioned templates hosted in any git repository.
Complete documentation is available at https://github.com/gonstr/rig.

Usage:
  rig [command]

Available Commands:
  build       Builds a rig.yaml template to stdout
  help        Help about any command
  install     Installs a rig template in the current directory
  version     Print the version number of rig

Flags:
  -h, --help   help for rig

Use "rig [command] --help" for more information about a command.
```

#### Installing

Either download release assets from the latest github release or install using go:

```shell
go get github.com/gonstr/rig
```
