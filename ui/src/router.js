import Vue from "vue"
import Router from "vue-router"
import db from "db"

import Login from "@/login"
import Admin from "@/admin"
import AdminApp from "@/admin-app"
import AdminVersion from "@/admin-version"
import Pkg from "@/pkg"
import Version from "@/version"
import App from "@/app"
import NotFound from "@/404"

Vue.use(Router)

export default new Router({
  mode: "history",
  routes: [
    // admin
    {
      path: "/",
      redirect: () => db.token ? "/admin" : "/login",
    },
    {
      path: "/login",
      component: Login,
    },
    {
      path: "/admin",
      component: Admin,
      children: [
        {
          path: ":id",
          component: AdminApp,
        },
        {
          path: ":id/version/:version",
          component: AdminVersion,
        },
      ],
    },
    // customer
    {
      path: "/pkg/:id",
      component: Pkg,
    },
    {
      path: "/version/:id",
      component: Version,
    },
    {
      path: "/:id",
      component: App,
    },
    {
      path: "*",
      component: NotFound,
    },
  ],
})
