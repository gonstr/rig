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
