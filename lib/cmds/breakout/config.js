const path = require('path')

const FILE_CONFIGS = {
  'build.yaml': {
    schemaFilePath: './templates/build.yaml/schema.json',
    templateDirPaths: {
      Jenkinsfile: './templates/build.yaml/Jenkinsfile'
    },
    defaultTemplateDirPath: 'Jenkinsfile',
    doneMessage: 'Done.'
  },
  'k8s.yaml': {
    schemaFilePath: './templates/k8s.yaml/schema.json',
    templateDirPaths: {
      Kustomize: './templates/k8s.yaml/kustomize',
      GoTemplates: './templates/k8s.yaml/go-templates'
    },
    defaultTemplateDirPath: 'Kustomize',
    doneMessage:
      'Done. You might need to adapt your build.yaml/Jenkinsfile to the manifest changes.'
  }
}

function configFrom(fileNamePath, outputDirPath) {
  const fileConfig = FILE_CONFIGS[path.basename(fileNamePath)]
  if (!fileConfig) throw new Error(`Unknown file type: ${fileNamePath}`)

  return {
    schemaFilePath: fileConfig.schemaFilePath,
    templateDirPaths: fileConfig.templateDirPaths,
    fileNamePath: fileNamePath,
    outputDirPath: outputDirPath || '.',
    doneMessage: fileConfig.doneMessage
  }
}

module.exports = {
  configFrom
}
