## infra-cli

Command line tool for generating and managing infra configuration files.

```shell
Usage
  $ infra-cli init
  $ infra-cli validate [FILE]
  $ infra-cli build-k8s [FILE] [CLUSTER] [IMAGETAG]
  $ infra-cli breakout [FILE] [OUTPUTPATH]

Commands
  init
    Creates build.yaml and k8s.yaml files in the current working dir.

  validate
    Validates a build.yaml or k8s.yaml file.

  build-k8s
    Converts a k8s.yaml to Kubernetes manifests and prints them to stdout.

  breakout
    Converts build.yaml or k8s.yaml to Jenkinsfile or Kubernetes manifest
    files respectivly and writes them to file system.

Examples
  $ infra-cli init
  $ infra-cli validate build.yaml
  $ infra-cli build-k8s k8s.yaml int v1.2.3
  $ infra-cli breakout k8s.yaml .
```

#### Install and use with npm

```shell
$ npm install -g git+ssh://git@github.com/fareoffice/infra-cli.git
$ infra-cli
```

#### Install and use with yarn

```shell
$ yarn global add git+ssh://git@github.com/fareoffice/infra-cli.git
$ infra-cli
```

#### Use with docker run

```shell
$ docker run -v ${PWD}:/cwd fareoffice/infra-cli init
$ docker run -v ${PWD}:/cwd fareoffice/infra-cli breakout /cwd/build.yaml
```
