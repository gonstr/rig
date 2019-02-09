const meow = require('meow')

const breakout = require('./breakout')
const config = require('./config')
const { typeFrom } = require('./constants')

const cli = meow(
  `
    Usage
      $ breakout <flags> [TYPE] [OUTPUTPATH]
 
    Types
      build               Converts a build yaml to a Jenkinsfile
      k8s                 Converts a k8s yaml to a Kubernetes manifests

    Options
      --file-path, -f     File path to file being converted
      --keep-input, -k    Don't delete the input file

    Examples
      $ breakout build -f build.yaml .
      $ breakout -k -f kube.yaml k8s ./output
`,
  {
    flags: {
      'file-path': {
        type: 'string',
        alias: 'f'
      }
    }
  }
)

;(() => {
  switch (cli.input[0]) {
    case 'build':
    case 'k8s':
      return breakout(config(typeFrom(cli.input[0]), cli.flags), cli.input[1])
    default:
      cli.showHelp()
  }
})()
