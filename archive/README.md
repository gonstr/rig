## rig

Command line tool for downloading the building rig templates.

```
Usage
  $ rig [command] <args>

Commands
  install
    Install a rig template from a github repo to the current
    working directory.

  build
    Builds the installed template and prints it to stdout.
    Template parameters are read from the rig.yaml file and
    the -p argument.

    --parameters, -p
      Template parameters in a comma separated key=value
      list. Parameters passed this way overrides parameters
      defined in the rig.yaml file.

    --output-dir, -o
      Writes template to a directory instead of stdout.

Examples
  $ rig install fareoffice/rig-templates/simple-app
  $ rig install fareoffice/rig-templates/simple-app#1.0.0
  $ rig build
  $ rig build -p key1=value,key2=value --output-dir ./output
```

#### Install and use with npm

```shell
$ npm install -g git+ssh://git@github.com/gonstr/rig.git
$ rig
```

#### Install and use with yarn

```shell
$ yarn global add git+ssh://git@github.com/gonstr/rig.git
$ rig
```
