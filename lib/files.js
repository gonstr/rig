const { promisify } = require('util')
const fs = require('fs')
const fsReadFile = promisify(fs.readFile)
const fsUnlink = promisify(fs.unlink)
const yaml = require('js-yaml')

async function readFile(path) {
  const content = await fsReadFile(path, 'utf8')

  if (path.endsWith('.json')) return JSON.parse(content)
  if (path.endsWith('.yaml')) return yaml.safeLoad(content)
  return content
}

function deleteFile(path) {
  return fsUnlink(path)
}

module.exports = {
  readFile,
  deleteFile
}
