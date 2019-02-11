const path = require('path')

const FILE_CONFIGS = {
  'k8s.yaml': {
    schemaFilePath: './templates/k8s/schema.json',
    templateFilePath: './templates/k8s/manifests.yaml'
  }
}

function configFrom(fileNamePath, cluster, imageTag) {
  const fileConfig = FILE_CONFIGS[path.basename(fileNamePath)]

  if (!fileConfig) throw new Error(`Unknown file type: ${fileNamePath}`)
  if (cluster !== 'dev' && cluster !== 'int' && cluster !== 'prod') {
    throw new Error('Cluster must be one of: int, dev or prod.')
  }
  if (!imageTag) throw new Error('No image tag provided.')

  return Object.assign({}, fileConfig, {
    fileNamePath: fileNamePath,
    cluster,
    imageTag
  })
}

module.exports = {
  configFrom
}
