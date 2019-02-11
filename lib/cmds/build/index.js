const handlebars = require('handlebars')
const path = require('path')
const fs = require('fs-extra')

const files = require('../../files')
const validate = require('../../validation')
const { configFrom } = require('./config')

async function run(args, flags) {
  const config = configFrom(...args)

  // Read input file content
  const input = Object.assign(await files.readFile(config.fileNamePath), {
    cluster: config.cluster,
    imageTag: config.imageTag
  })

  // Validate input file
  const { valid, errors } = await validate(input, config.schemaFilePath)
  if (!valid) {
    return console.error(`Invalid file schema in file '${config.fileNamePath}': ${errors}`)
  }

  // Write template to stdout
  await write(input, config.templateFilePath)
}

async function write(values, templateFilePath) {
  const file = await fs.readFile(templateFilePath, 'utf8')

  const template = handlebars.compile(file)

  const processed = template(values)

  console.log(processed)
}

module.exports = run
