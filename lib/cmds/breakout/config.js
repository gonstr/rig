const path = require('path')

const FILE_CONFIGS = {
  'build.yaml': {
    schemaFilePath: './templates/build/schema.json',
    templateDirPath: './templates/build/files'
  },
  'k8s.yaml': {
    schemaFilePath: './templates/k8s/schema.json',
    templateDirPath: './templates/k8s/files'
  }
}

function configFrom(fileNamePath, outputDirPath, flags) {
  const fileConfig = FILE_CONFIGS[path.basename(fileNamePath)]

  if (!fileConfig) throw new Error(`Unknown file type: ${fileNamePath}`)

  return Object.assign({}, fileConfig, {
    fileNamePath: fileNamePath,
    outputDirPath: outputDirPath || '.',
    keepInputFile: flags.k || false
  })
}

module.exports = {
  configFrom
}
