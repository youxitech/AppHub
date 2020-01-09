import Vue from "vue"
import router from "./router"
import store from "./store"
import "./global"
import "./style.styl"
import Entry from "./entry"

window.Vue = Vue
Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  el: "#root",
  router,
  store,
  render: h => h(Entry),
})
