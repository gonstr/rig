const gh = require('parse-github-url')

const RigError = require('../../error')

function configFrom(templateUri) {
  const { owner, host, name, branch, hash } = gh(templateUri)

  if (!owner || !name || !branch) throw new RigError('Unable to parse git uri')

  return {
    host,
    owner,
    repo: name,
    template: branch,
    version: hash ? hash.replace('#', '') : null
  }
}

module.exports = {
  configFrom
}
