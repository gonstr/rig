function configFrom(fileNamePath, cluster, imageTag) {
  if (cluster !== 'dev' && cluster !== 'int' && cluster !== 'prod') {
    throw new Error('Cluster must be one of: int, dev or prod.')
  }

  if (!imageTag) throw new Error('No image tag provided.')

  return {
    schemaFilePath: './templates/k8s.yaml/schema.json',
    templateFilePath: './templates/k8s.yaml/manifests.yaml',
    fileNamePath: fileNamePath,
    cluster,
    imageTag
  }
}

module.exports = {
  configFrom
}
