const glob = require('glob')
const handlebars = require('handlebars')
const path = require('path')
const fs = require('fs-extra')

function write(values, outputDirPath, templateDirPath) {
  const files = glob.sync(`${templateDirPath}/**/*`, {
    nodir: true
  })

  return Promise.all(
    files.map(async filePath => {
      const file = await fs.readFile(filePath, 'utf8')

      const template = handlebars.compile(file)

      const processed = template(values)

      const newFilePath = path.join(outputDirPath, filePath.replace(templateDirPath, ''))

      await fs.ensureFile(newFilePath)
      await fs.writeFile(newFilePath, processed)
    })
  )
}

module.exports = {
  write
}
