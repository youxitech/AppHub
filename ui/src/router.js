import Vue from "vue"
import Router from "vue-router"
import db from "db"

import Login from "@/login"
import Dashboard from "@/dashboard"
import Pkg from "@/pkg"
import Version from "@/version"
import NotFound from "@/404"

Vue.use(Router)

export default new Router({
  mode: "history",
  routes: [
    {
      path: "/",
      redirect: () => db.token ? "/dashboard" : "/login",
    },
    {
      path: "/login",
      component: Login,
    },
    {
      path: "/dashboard",
      component: Dashboard,
    },
    {
      path: "/pkg/:id",
      component: Pkg,
    },
    {
      path: "/version/:id",
      component: Version,
    },
    {
      path: "*",
      component: NotFound,
    },
  ],
})
