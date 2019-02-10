## infra-cli

Command line tool for generating and managing infra configuration files.

```shell
Usage
    $ infra-cli init
    $ infra-cli validate [FILE]
    $ infra-cli breakout [-k keep-input] [FILE] [OUTPUTPATH]

Commands
    init
        Creates build.yaml and k8s.yaml files in the current working dir.

    validate
        Validates a build.ysml or k8s.yaml file

    breakout
        Converts build.yaml or k8s.yaml to Jenkinsfile or Kubernetes manifest
        files respectivly. Providing the '--keep-input', or '-k' flag ensures
        the input file is not deleted.

Examples
    $ infra-cli init
    $ infra-cli breakout -k k8s.yaml .
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
$ docker run -v ${PWD}:/cwd fareoffice/infra-cli breakout /cwd/k8s.yaml /cwd/deploy
```
