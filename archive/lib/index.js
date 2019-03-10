const meow = require('meow')

const install = require('./cmds/install')
const build = require('./cmds/build')
const RigError = require('./error')

const m = meow(
  `
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
`,
  {
    flags: {
      parameters: {
        type: 'string',
        alias: 'p'
      },
      'output-dir': {
        type: 'string',
        alias: 'o'
      }
    }
  }
)

;(async () => {
  const [cmd, ...args] = m.input
  try {
    switch (cmd) {
      case 'install':
        await install(args)
        break
      case 'build':
        await build(args)
        break
      default:
        m.showHelp()
    }
  } catch (err) {
    if (err instanceof RigError) {
      console.error(err.message)
    } else {
      throw err
    }
  }
})()
