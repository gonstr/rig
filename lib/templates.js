const path = require('path')
const fs = require('fs-extra')
const os = require('os')
const { exec } = require('child_process')
const tmp = require('tmp')
const yaml = require('js-yaml')

const RigError = require('./error')

function shell(cmd, cwd) {
  return new Promise((resolve, reject) => {
    exec(cmd, { cwd }, (err, stdout, stderr) => {
      if (err) reject(new RigError(err.message))
      else resolve({ stdout, stderr })
    })
  })
}

async function clean(dir) {
  await shell('git fetch --tags', dir)
  await shell('git checkout master', dir)
  await shell('git clean -d -f', dir)
  await shell('git pull -q origin master', dir)
}

function createDirs(host, owner, repo, template) {
  const dirs = {
    rig: path.join(os.homedir(), '.rig')
  }

  if (owner) dirs.owner = path.join(dirs.rig, host, owner)
  if (repo) dirs.repo = path.join(dirs.owner, repo)
  if (template) dirs.template = path.join(dirs.repo, template)

  return dirs
}

async function checkoutVersion(host, owner, repo, template, version) {
  const dirs = createDirs(host, owner, repo, template)

  const { name } = tmp.dirSync({ unsafeCleanup: false })

  const tag = version ? `tags/${template}#${version}` : 'master'

  await shell(`git --work-tree=${name} checkout ${tag} -- ${template}`, dirs.repo)

  return name
}

async function clone(host, owner, repo) {
  const dirs = createDirs(host, owner, repo)

  await fs.ensureDir(dirs.owner)

  if (await fs.pathExists(dirs.repo)) {
    await clean(dirs.repo)
  } else {
    const uri = `git@${host}:${owner}/${repo}`
    await shell(`git clone ${uri}`, dirs.owner)
  }
}

async function parameters(host, owner, repo, template, version) {
  const dirs = createDirs(host, owner, repo, template)

  if (!(await fs.pathExists(dirs.template))) {
    throw new RigError(`Template '${repo}/${template}' does not exist.`)
  }

  await clean(dirs.repo)

  const versionDir = await checkoutVersion(host, owner, repo, template, version)

  const content = await fs.readFile(path.join(versionDir, template, 'parameters.yaml'))

  return yaml.safeLoad(content)
}

module.exports = {
  clone,
  parameters
}
