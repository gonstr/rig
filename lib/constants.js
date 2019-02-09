function typeFrom(type) {
  switch (type) {
    case 'build':
      return TYPES.build
    case 'k8s':
      return TYPES.k8s
    default:
      throw new Error('Unknown type')
  }
}

const TYPES = {
  build: {
    code: 'build',
    defaultFilePath: 'build.yaml',
    schemaFilePath: './templates/build/schema.json',
    templateDirPath: './templates/build/files'
  },
  k8s: {
    code: 'k8s',
    defaultFilePath: 'k8s.yaml',
    schemaFilePath: './templates/k8s/schema.json',
    templateDirPath: './templates/k8s/files'
  }
}

module.exports = {
  typeFrom,
  TYPES
}
