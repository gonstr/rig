const fs = require('fs-extra')
const yaml = require('js-yaml')

async function readFile(path) {
  const content = await fs.readFile(path, 'utf8')

  if (path.endsWith('.json')) return JSON.parse(content)
  if (path.endsWith('.yaml')) return yaml.safeLoad(content)
  return content
}

function deleteFile(path) {
  return fs.unlink(path)
}

module.exports = {
  readFile,
  deleteFile
}
