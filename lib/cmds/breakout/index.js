const files = require('../../files')
const validate = require('../../validation')
const templates = require('./templates')
const { configFrom } = require('./config')

async function run(args, flags) {
  const config = configFrom(args[0], args[1], flags)

  // Read input file content
  const input = await files.readFile(config.fileNamePath)

  // Validate input file
  const { valid, errors } = await validate(input, config.schemaFilePath)
  if (!valid) {
    return console.error(`Invalid file schema in file '${config.fileNamePath}': ${errors}`)
  }

  // Write template files
  await templates.write(input, config.outputDirPath, config.templateDirPath)

  // Delete input files
  if (!config.keepInputFile) await files.deleteFile(config.fileNamePath)

  console.log('Done.')
}

module.exports = run
