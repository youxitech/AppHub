import Vue from "vue"
import router from "./router"
import store from "./store"
import "./global"
import "./style.styl"

window.Vue = Vue
Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  el: "#root",
  router,
  store,
  render: h => h(Vue.component("RouterView")),
})
