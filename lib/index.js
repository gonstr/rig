const meow = require('meow')

const init = require('./cmds/init')
const validate = require('./cmds/validate')
const build = require('./cmds/build')
const breakout = require('./cmds/breakout')

const m = meow(
  `
  Usage
    $ infra-cli init
    $ infra-cli validate [FILE]
    $ infra-cli build [FILE] [CLUSTER] [IMAGETAG]
    $ infra-cli breakout [FILE] [OUTPUTPATH]

  Commands
    init
      Creates build.yaml and k8s.yaml files in the current working dir.

    validate
      Validates a build.yaml or k8s.yaml file.
    
    build
      Converts a k8s.yaml to Kubernetes manifests and prints them to stdout.

    breakout
      Converts build.yaml or k8s.yaml to Jenkinsfile or Kubernetes manifest
      files respectivly and writes them to file system.

  Examples
    $ infra-cli init
    $ infra-cli build k8s.yaml int v1.2.3
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
    case 'build':
      return build(args, m.flags)
    case 'breakout':
      return breakout(args)
    default:
      m.showHelp()
  }
})()
