## rig

Command line tool for generating and managing infra configuration files.

```shell
Usage
  $ rig [command] <args>

Commands
  install
    Install a rig template from a github repo in the current
    working directory. Valid template arguments are:

    <owner>/<repo>/<template>[#<version>]

    For example:

    fareoffice/rig-templates/simple-app#1.0.0

  build
    Builds the installed template and prints it to stdout.

    --parameters, -p
      Template parameters in a comma separated key, value
      list. Parameters passed this way overrides parameters
      defined in the .rig file.

    --output-dir, -o
      Writes template to a directory instead of stdout.

Examples
  $ rig install fareoffice/rig-templates/simple-app
  $ rig install fareoffice/rig-templates/simple-app#1.0.0
  $ rig build
  $ rig build --parameters key=value,key=value --to-dir ./output
```

#### Install and use with npm

```shell
$ npm install -g git+ssh://git@github.com/fareoffice/rig.git
$ rig
```

#### Install and use with yarn

```shell
$ yarn global add git+ssh://git@github.com/fareoffice/rig.git
$ rig
```
