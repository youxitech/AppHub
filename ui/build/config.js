module.exports = {
  port: 9200,
  proxyTable: {
    "!(/api)**/v*.*": {
      target: 'http://localhost:9200/',
      router: function (req) {
        req.url = "index.html"
      }
    },
    "/api": {
      changeOrigin: true,
      target: "https://apphub.haibao6688.com",
    },
    "/data": {
      changeOrigin: true,
      target: "https://apphub.haibao6688.com",
    },
  },
}
