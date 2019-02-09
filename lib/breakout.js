const files = require('./files')
const validate = require('./validation')
const templates = require('./templates')

async function breakout(config, outputDirPath) {
  // Read input file content
  const input = await files.readFile(config.filePath)

  // Validate input file
  const { valid, errors } = await validate(input, config.type)
  if (!valid) return console.error(`Invalid file schema in file '${config.filePath}': ${errors}`)

  // Write template files
  await templates.write(input, outputDirPath, config.type)

  // Delete input files
  if (!config.keepInputFile) await files.deleteFile(config.filePath)

  console.log('Done.')
}

module.exports = breakout
