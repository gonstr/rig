const meow = require('meow')

const init = require('./cmds/init')
const validate = require('./cmds/validate')
const breakout = require('./cmds/breakout')

const m = meow(
  `
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
`,
  {
    flags: {
      'keep-input': {
        type: 'boolean',
        alias: 'k'
      }
    }
  }
)

;(() => {
  const [cmd, ...args] = m.input
  switch (cmd) {
    case 'init':
      return init()
    case 'validate':
      return validate(args)
    case 'breakout':
      return breakout(args, m.flags)
    default:
      m.showHelp()
  }
})()
