const meow = require('meow')

const breakout = require('./cmds/breakout')

const m = meow(
  `
    Usage
      $ infra-cli init
      $ infra-cli breakout [-k keep-input] [FILE] [OUTPUTPATH]
 
    Commands
      init
        Creates build.yaml and k8s.yaml files in the current working dir.
      
      breakout
        Converts the build.yaml or k8s.yaml to their more complex Jenkinsfile
        and kubernetes manifest files respectivly. Providing the '--keep-input',
        or '-k' flag ensures the input file is not deleted.

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
      throw new Error('Not implemented yet')
    case 'breakout':
      return breakout(args, m.flags)
    default:
      m.showHelp()
  }
})()
