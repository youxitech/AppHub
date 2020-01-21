module.exports = {
  port: 9200,
  proxyTable: {
    "!(/api)**/v*.*": {
      target: 'http://localhost:8880/',
      router: function (req) {
        req.url = "index.html"
      }
    },
    "/api": {
      changeOrigin: true,
      target: "http://192.168.10.88:3389",
    },
    "/data": {
      changeOrigin: true,
      target: "http://192.168.10.88:3389",
    },
  },
}
