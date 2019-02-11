const glob = require('glob')
const handlebars = require('handlebars')
const path = require('path')
const fs = require('fs-extra')
const files = require('../../files')
const validate = require('../../validation')
const { configFrom } = require('./config')

async function run(args) {
  const config = configFrom(...args)

  // Read input file content
  const input = await files.readFile(config.fileNamePath)

  // Validate input file
  const { valid, errors } = await validate(input, config.schemaFilePath)
  if (!valid) {
    return console.error(`Invalid file schema in file '${config.fileNamePath}': ${errors}`)
  }

  // Write template files
  await write(input, config.templateDirPath, config.outputDirPath)

  // Delete input files
  await files.deleteFile(config.fileNamePath)

  console.log('Done.')
}

function write(values, templateDirPath, outputDirPath) {
  const files = glob.sync(`${templateDirPath}/**/*`, {
    nodir: true
  })

  return Promise.all(
    files.map(async filePath => {
      const file = await fs.readFile(filePath, 'utf8')

      const template = handlebars.compile(file)

      const processed = template(values)

      const newFilePath = path.join(outputDirPath, filePath.replace(templateDirPath, ''))

      await fs.ensureFile(newFilePath)
      await fs.writeFile(newFilePath, processed)
    })
  )
}

module.exports = run
