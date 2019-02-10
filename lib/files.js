const fs = require('fs-extra')
const yaml = require('js-yaml')

async function readFile(path) {
  const content = await fs.readFile(path, 'utf8')

  if (path.endsWith('.json')) return JSON.parse(content)
  if (path.endsWith('.yaml')) return yaml.safeLoad(content)
  return content
}

async function writeFile(path, content) {
  await fs.ensureFile(path)
  if (path.endsWith('.yaml')) return fs.writeFile(path, yaml.safeDump(content))
  if (path.endsWith('.json')) return fs.writeFile(path, JSON.stringify(content, null, 2))
  return fs.writeFile(path, content)
}

function deleteFile(path) {
  return fs.unlink(path)
}

module.exports = {
  readFile,
  writeFile,
  deleteFile
}
