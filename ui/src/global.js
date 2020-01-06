import "@babel/polyfill"
import qs from "qs"
import axios from "axios"
import * as util from "util"
import displayError from "error"
import Vue from "vue"

// => Vue Config
Vue.config.warnHandler = (msg, vm, trace) => {
  if(msg.startsWith("Invalid component name")) return
  console.error(msg, trace)
}

Object.defineProperty(Vue.prototype, "$query", {
  get() {
    return qs.parse(window.location.search.slice(1))
  },
})

// => Global Variables
window.axios = axios
window._util = util
window._displayError = displayError
window._global = {
  isDev: location.port !== "",
}

// => Axios Config
axios.defaults.baseURL = "/api"
axios.interceptors.request.use(config => {
  // add custom headers
  return config
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
