const handlebars = require('handlebars')
const { configFrom } = require('./config')

async function run(args, flags) {
  const config = configFrom(...args)
}

async function write() {}

module.exports = run
