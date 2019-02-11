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

function configFrom(fileNamePath, outputDirPath, templateType) {
  const fileConfig = FILE_CONFIGS[path.basename(fileNamePath)]
  if (!fileConfig) throw new Error(`Unknown file type: ${fileNamePath}`)

  const templateDirPath = templateType
    ? fileConfig.templateDirPaths[templateType]
    : fileConfig.templateDirPaths[fileConfig.defaultTemplateDirPath]
  if (!templateDirPath) throw new Error(`Unknown template type: ${templateType}`)

  return {
    schemaFilePath: fileConfig.schemaFilePath,
    templateDirPath,
    templateDirName: templateDirPath.split(path.sep).pop(),
    fileNamePath: fileNamePath,
    outputDirPath: outputDirPath || '.',
    doneMessage: fileConfig.doneMessage
  }
}

module.exports = {
  configFrom
}
