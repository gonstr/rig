const rl = require('readline-sync')
const path = require('path')

const files = require('../../files')

const cwd = path
  .dirname(path.join(__dirname, '../..'))
  .split(path.sep)
  .pop()

function run() {
  const repoName = rl.question(`Repository name (${cwd}): `) || cwd
  const appName = rl.question(`App name (${repoName}): `) || repoName
  const namespace = rl.question('Namespace: ')
  const servicePort = rl.question('Service port (80): ') || 80
  const readinessPath = rl.question('Readiness path (/health): ') || '/health'
  const cpuMin = rl.question('CPU min (500m): ') || '500m'
  const cpuMax = rl.question('CPU max (1000m): ') || '1000m'
  const memory = rl.question('Memory (800mi): ') || '800m'

  const input = {
    repoName,
    appName,
    namespace,
    servicePort,
    readinessPath,
    cpuMin,
    cpuMax,
    memory
  }

  const content = Object.assign({ version: 1 }, input)

  return files.writeFile('k8s.yaml', content)
}

module.exports = run
