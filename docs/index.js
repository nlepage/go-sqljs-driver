const { readFile } = require('fs/promises')
const { resolve } = require('path')
global.initSqlJs = require('sql.js')

require('./go_wasm_exec_1.16.7')

async function run() {
  const file = await readFile(resolve(__dirname, 'index.wasm'))
  const wasm = await WebAssembly.compile(file)
  const go = new Go()
  const instance = await WebAssembly.instantiate(wasm, go.importObject)
  await go.run(instance)
}

run().catch(e => {
  console.error(e)
  process.exit(1)
})
