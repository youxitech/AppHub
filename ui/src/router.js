import Vue from "vue"
import Router from "vue-router"

import App from "app"

import Home from "@/home"
import NotFound from "@/404"

Vue.use(Router)

export default new Router({
  mode: "history",
  routes: [
    {
      path: "/",
      component: App,
      children: [
        {
          path: "",
          redirect: "home",
        },
        {
          path: "/home",
          component: Home,
        },
      ],
    },
    {
      path: "*",
      component: NotFound,
    },
  ],
})
