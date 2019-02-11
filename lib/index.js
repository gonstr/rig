const meow = require('meow')

const init = require('./cmds/init')
const validate = require('./cmds/validate')
const genk8s = require('./cmds/gen-k8s')
const breakout = require('./cmds/breakout')

const m = meow(
  `
  Usage
    $ infra-cli init
    $ infra-cli validate [file]
    $ infra-cli gen-k8s [file] [cluster] [image tag]
    $ infra-cli breakout [file] [output path] <template type>

  Commands
    init
      Creates build.yaml and k8s.yaml files in the current working directory.

    validate
      Validates a build.yaml or k8s.yaml file.

    gen-k8s
      Generates kubernetes manifests from a k8s.ysml file and prints them to
      stdout.

    breakout
      Converts build.yaml or k8s.yaml to Jenkinsfile or Kubernetes manifest
      files respectivly and writes them to the specified output path. Output
      template can be specified as an optional third argument. Valid templates
      are:

      k8s.yaml:
        Kustomize (default) - converts k8s.yaml to kustomize yaml files.
        GoTemplates - converts k8s.yaml to go-template yaml files.

      build.yaml
        Jenkinsfile - converts build.yaml to Jenkinsfile.

  Examples
    $ infra-cli init
    $ infra-cli validate build.yaml
    $ infra-cli gen-k8s k8s.yaml int v1.2.3
    $ infra-cli breakout k8s.yaml .
`
)

;(() => {
  const [cmd, ...args] = m.input
  switch (cmd) {
    case 'init':
      return init()
    case 'validate':
      return validate(args)
    case 'gen-k8s':
      return genk8s(args, m.flags)
    case 'breakout':
      return breakout(args)
    default:
      m.showHelp()
  }
})()
