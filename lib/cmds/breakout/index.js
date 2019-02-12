const glob = require('glob')
const handlebars = require('handlebars')
const path = require('path')
const fs = require('fs-extra')
const readlineSync = require('readline-sync')

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

  // Select template type
  const templates = Object.keys(config.templateDirPaths)
  const index = readlineSync.keyInSelect(templates, 'Choose template', { cancel: false })
  const templateDirName = templates[index]
  const templateDirPath = config.templateDirPaths[templateDirName]

  if (
    readlineSync.keyInYN(
      `Generate template '${templateDirName}' to folder '${config.outputDirPath}'?`
    )
  ) {
    // Write template files
    await write(input, templateDirPath, config.outputDirPath)

    console.log(config.doneMessage)
  }
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
