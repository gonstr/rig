function config(type, flags) {
  return {
    type,
    filePath: flags.f || type.defaultFilePath,
    keepInputFile: flags.k || false
  }
}

module.exports = config
