const Validator = require('jsonschema').Validator

const files = require('./files')

async function validate(content, schemePath) {
  const schema = await files.readFile(schemePath)

  const validator = new Validator()

  return validator.validate(content, schema)
}

module.exports = validate
