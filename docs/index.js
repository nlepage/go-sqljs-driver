global.initSqlJs = require('sql.js')
require('./go_wasm_exec_1.16.7')

async function run() {
  const go = new Go()
  const wasm = await WebAssembly.compile(fs.readFileSync('./index.wasm'))
  const instance = await WebAssembly.instantiate(wasm, go.importObject)
  await go.run(instance)
}

run().catch(e => {
  console.error(e)
  process.exit(1)
})
