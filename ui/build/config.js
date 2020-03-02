module.exports = {
  port: 9200,
  proxyTable: {
    "**/v*.*": {
      target: 'http://192.168.10.88:3389/',
      changeOrigin: true,
      router: function (req) {
        if(!req.url.startsWith("/api")) {
          req.url = "http://localhost:9200/index.html"
        }
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
