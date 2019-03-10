const path = require('path')
const _ = require('lodash')
const fs = require('fs-extra')
const yaml = require('js-yaml')

const { configFrom } = require('./config')
const templates = require('../../templates')

async function run(args) {
  const { host, owner, repo, template, version } = configFrom(...args)

  console.log(`Installing ${repo}/${template}${version ? `#${version}` : ''}... `)

  await templates.clone(host, owner, repo)

  const params = await templates.parameters(host, owner, repo, template, version)

  const rig = {
    template: `${owner}/${repo}/${template}`,
    version,
    parameters: _.mapValues(params, spec => spec.default || '')
  }

  await fs.writeFile(path.join(process.cwd(), 'rig.yaml'), yaml.safeDump(rig))

  console.log(`Done. Add static parameters to rig.yaml then build the template with 'rig build'.`)
}

module.exports = run
