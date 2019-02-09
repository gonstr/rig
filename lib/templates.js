const glob = require('glob')
const handlebars = require('handlebars')
const path = require('path')

function write(values, outputDirPath, type) {
  const files = glob.sync(`${type.templateDirPath}/**/*`, {
    nodir: true
  })

  return Promise.all(
    files.map(async filePath => {
      const file = await fs.readFile(filePath, 'utf8')

      const template = handlebars.compile(file)

      await fs.writeFile(path.join(outputDirPath, filePath), template(values))
    })
  )
}

module.exports = {
  write
}
