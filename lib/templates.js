const shell = require('shelljs')
const path = require('path')
const fs = require('fs-extra')
const os = require('os')

const files = require('./files')
const RigError = require('./error')

async function clone(host, owner, repo) {
  const rigDir = path.join(os.homedir(), '.rig')
  const ownerDir = path.join(rigDir, host, owner)
  const repoDir = path.join(ownerDir, repo)

  await fs.ensureDir(ownerDir)

  shell.cd(ownerDir)

  if (await fs.pathExists(repoDir)) {
    shell.cd(repoDir)
    console.log(`Updating ${repo}... `)

    let err
    if ((err = shell.exec(`git fetch --tags`, { silent: true }).stderr)) throw new RigError(err)
    if ((err = shell.exec(`git pull -q origin master`, { silent: true }).stderr)) {
      throw new RigError(err)
    }
    if ((err = shell.exec(`git clean -d -f`, { silent: true }).stderr)) throw new RigError(err)
  } else {
    const uri = `git@${host}:${owner}/${repo}`
    shell.exec(`git clone ${uri}`)
  }
}

async function parameters(host, owner, repo, template, version) {
  const repoDir = path.join(os.homedir(), '.rig', host, owner, repo)
  const templateDir = path.join(repoDir, template)

  if (!(await fs.pathExists(templateDir))) {
    throw new RigError(`Template '${repo}/${template}' does not exist.`)
  }

  shell.cd(repoDir)

  if (version) {
    const tag = `${template}#${version}`

    if (!shell.exec(`git tag -l ${tag}`, { silent: true }).stdout) {
      throw new RigError(`Could not find version '${version}' for template '${template}'.`)
    }

    let err
    if ((err = shell.exec(`git checkout -q tags/${tag}`).stderr)) throw new RigError(err)
  } else {
    let err
    if ((err = shell.exec('git checkout -q master').stderr)) throw new RigError(err)
  }

  return files.readFile(path.join(templateDir, 'parameters.yaml'))
}

module.exports = {
  clone,
  parameters
}
