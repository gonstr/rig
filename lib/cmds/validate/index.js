const files = require('../../files')
const validate = require('../../validation')
const { configFrom } = require('./config')

async function run(args) {
  const config = configFrom(args[0])

  // Read file content
  const input = await files.readFile(config.fileNamePath)

  // Validate file
  const { valid, errors } = await validate(input, config.schemaFilePath)
  if (!valid) {
    return console.error(`Invalid file schema in file '${config.fileNamePath}': ${errors}`)
  }

  console.log('File is valid.')
}

module.exports = run
