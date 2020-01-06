const config = require("./config")
const opn = require("opn")
const express = require("express")
const webpack = require("webpack")
const proxyMiddleware = require("http-proxy-middleware")
const webpackConfig = require("./webpack.dev.conf")

const app = express()
const compiler = webpack(webpackConfig)

const devMiddleware = require("webpack-dev-middleware")(compiler, {
  publicPath: webpackConfig.output.publicPath,
  stats: "minimal",
})

const hotMiddleware = require("webpack-hot-middleware")(compiler, {
  log: false,
})

// force page reload when html-webpack-plugin template changes
compiler.hooks.compilation.tap("html-webpack-plugin-after-emit", () => {
  hotMiddleware.publish({action: "reload"})
})

// proxy api requests
const proxyTable = config.proxyTable
Object.keys(proxyTable).forEach(key => {
  let options = proxyTable[key]
  if(typeof options === "string") {
    options = {target: options}
  }
  app.use(proxyMiddleware(options.filter || key, options))
})

// handle fallback for HTML5 history API
app.use(require("connect-history-api-fallback")())

// serve webpack bundle output
app.use(devMiddleware)

// enable hot-reload and state-preserving
// compilation error display
app.use(hotMiddleware)

// serve pure static assets
app.use("/static", express.static("./static"))

// default port where dev server listens for incoming traffic
const port = process.env.PORT || config.port
const url = "http://localhost:" + port

console.log("> Starting dev server...")

// use `yarn dev -n` to disable browser auto open
const doNotAutoOpenBrowser = process.argv.includes("-n")

devMiddleware.waitUntilValid(() => {
  console.log(`> Listening at ${url}\n`)
  if(!doNotAutoOpenBrowser) {
    opn(url)
  }
})

const server = app.listen(port)
