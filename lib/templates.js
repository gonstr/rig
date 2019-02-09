const glob = require('glob')
const handlebars = require('handlebars')
const path = require('path')
const fs = require('fs-extra')

function write(values, outputDirPath, type) {
  const files = glob.sync(`${type.templateDirPath}/**/*`, {
    nodir: true
  })

  return Promise.all(
    files.map(async filePath => {
      const file = await fs.readFile(filePath, 'utf8')

      const template = handlebars.compile(file)

      const parsed = template(values)

      const newFilePath = path.join(outputDirPath, filePath.replace(type.templateDirPath, ''))

      await fs.ensureFile(newFilePath)
      await fs.writeFile(newFilePath, parsed)
    })
  )
}

module.exports = {
  write
}
