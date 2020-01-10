import "@babel/polyfill"
import axios from "axios"
import * as util from "util"
import Vue from "vue"
import Notifications from "vue-notification"
import "normalize.css"
import db from "db"
import router from "./router"
import { format } from "date-fns"
import vmodal from "vue-js-modal"

// => Vue Config
Vue.config.warnHandler = (msg, vm, trace) => {
  if(msg.startsWith("Invalid component name")) return
  console.error(msg, trace)
}
Vue.use(Notifications)
Vue.use(vmodal)
Vue.mixin({
  methods: {
    _log: console.log,
    _getAsset: util.getAsset,
  },

  filters: {
    formatTime(time) {
      return format(new Date(time), "yyyy-MM-dd HH:mm")
    },

    bytesToSize(bytes) {
      const sizes = ["Bytes", "KB", "MB", "GB", "TB"]
      if(bytes === 0) return "0 Byte"
      const i = parseInt(Math.floor(Math.log(bytes) / Math.log(1024)))
      return Math.round(bytes / Math.pow(1024, i), 2) + " " + sizes[i]
    },
  },
})

// => Global Variables
window.axios = axios
window._util = util
window._displayError = util.displayError
window._db = db
window._global = {
  isDev: location.port !== "",
}

// => Axios Config
axios.defaults.baseURL = "/api"
axios.interceptors.request.use(config => {
  config.headers.common["X-Admin-Token"] = db.token
  return config
})
axios.interceptors.response.use(null, e => {
  if(e.response.status === 401) {
    router.push("/login")
  }
  return Promise.reject(e)
})

// => Global Components
const context = require.context("./components/_common/", true, /\.(vue|js)$/)

// key is the file path, e.g. `./_loading.vue`
context.keys().forEach(key => {
  // dir, e.g. `./_imgs-picker/index.vue`
  if(key.lastIndexOf("/") !== 1) {
    if(key.indexOf("index.vue") !== -1) {
      const name = key.slice(2, key.lastIndexOf("/"))
      Vue.component(name, context(key).default)
    }
  } else {
    const dotIndex = key.lastIndexOf(".")
    const name = key.slice(2, dotIndex)
    Vue.component(name, context(key).default)
  }
})
