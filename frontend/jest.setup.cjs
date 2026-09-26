// jsdom не отдаёт TextEncoder/TextDecoder, а react-router обращается к ним при импорте
const { TextDecoder, TextEncoder } = require('util')

Object.assign(globalThis, { TextDecoder, TextEncoder })
