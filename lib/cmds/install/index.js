const path = require('path')
const _ = require('lodash')

const { configFrom } = require('./config')
const templates = require('../../templates')
const files = require('../../files')

async function run(args) {
  const { host, owner, repo, template, version } = configFrom(...args)

  await templates.clone(host, owner, repo)

  const params = await templates.parameters(host, owner, repo, template, version)

  const rig = {
    template: `${owner}/${repo}/${template}`,
    version,
    parameters: _.mapValues(params, spec => spec.default || '')
  }

  files.writeFile(path.join(process.cwd(), '.rig.yaml'), rig)
}

module.exports = run
