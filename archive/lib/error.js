class RigError extends Error {
  constructor(...args) {
    super(...args)
    Error.captureStackTrace(this, RigError)
  }
}

module.exports = RigError
