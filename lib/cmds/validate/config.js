const path = require('path')

const FILE_CONFIGS = {
  'build.yaml': {
    schemaFilePath: './templates/build/schema.json'
  },
  'k8s.yaml': {
    schemaFilePath: './templates/k8s/schema.json'
  }
}

function configFrom(fileNamePath) {
  const fileConfig = FILE_CONFIGS[path.basename(fileNamePath)]

  if (!fileConfig) throw new Error(`Unknown file type: ${fileNamePath}`)

  return Object.assign({}, fileConfig, {
    fileNamePath: fileNamePath
  })
}

module.exports = {
  configFrom
}
