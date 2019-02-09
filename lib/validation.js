const Validator = require('jsonschema').Validator

const files = require('./files')

async function validate(content, type) {
  const schema = await files.readFile(type.schemaFilePath)

  const validator = new Validator()

  return validator.validate(content, schema)
}

module.exports = validate
