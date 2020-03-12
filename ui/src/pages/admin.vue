<template lang="pug">
.admin.h-screen.flex.flex-col
  header(class="w-full border-b border-gray-200 h-16 flex items-center p-6")
    a(href="/" class="block lg:mr-4") App Hub
    _button.uploader-trigger +
    _button.ml-auto(@click="logout") logout

  div(v-if="apps.length === 0") 请先上传 APP
  main(
    v-else
    class="w-full flex flex-1"
  )
    nav(class="w-40 border-b border-r")
      .flex.cursor-pointer(
        v-for="app in apps"
        :key="app.id"
        :class="'hover:bg-gray-300 p-3 w-40 ' + ($route.params.id === app.alias ? 'bg-gray-300' : '') "
        @click="() => onClickNav(app.alias)"
      )
        img.w-10.h-10.mr-2.rounded(:src="_getAsset('icon', app.platform, app.bundleID)")
        div
          div {{ app.name }}
          .text-xs.text-gray-500 {{ app.platform }}

    .flex-1.bg-gray-200.p-6
      router-view
</template>

<script>
import "@uppy/core/dist/style.css"
import "@uppy/dashboard/dist/style.css"
const Uppy = require("@uppy/core")
const Dashboard = require("@uppy/dashboard")
const XHRUpload = require("@uppy/xhr-upload")

export default {
  data() {
    return {
      apps: [],
    }
  },

  mounted() {
    this.fetchApps()
    this.initUploader()
  },

  watch: {
    $route() {
      this.fetchApps()
    },
  },

  methods: {
    fetchApps() {
      axios.get("/admin/apps")
        .then(res => {
          this.apps = res.data
          // 已选 app 或 没有 app 列表直接返回
          if(this.$route.params.id || this.apps.length === 0) return
          this.$router.push(`/admin/${ this.apps[0].alias }`)
        })
        .catch(_displayError)
    },

    initUploader() {
      const uppy = Uppy({ autoProceed: true })

      uppy
        .on("file-added", file => {
          let channel = ""
          let env = ""
          const { name } = file
          const nameArr = name.split(".")
          nameArr.splice(nameArr.length - 1, 1)
          const fileName = nameArr.join(".")
          const fileNameArr = fileName.split("-")
          if(fileNameArr.length === 3) {
            channel = fileNameArr[2]
            env = fileNameArr[1]
          }

          uppy.setFileMeta(file.id, {
            channel,
            env,
          })
        })
        .use(Dashboard, { trigger: ".uploader-trigger" })
        .use(XHRUpload, {
          endpoint: "/api/admin/upload",
          formData: true,
          fieldName: "file",
          metaFields: ["channel", "env"],
          headers: {
            "X-Admin-Token": _db.token,
          },
        })
        .on("complete", res => {
          if(res.failed.length > 0) {
            res.failed.forEach(item => {
              _showSuccess(`文件${ item.name }失败，原因：${ item.response.body.msg }`)
            })
            return
          }

          _showSuccess("全部上传成功")
          this.$router.push("/admin")
        })
    },

    onClickNav(alias) {
      this.$router.push(`/admin/${ alias }`)
    },

    logout() {
      _db.token = ""
      this.$router.push("/login")
    },
  },
}
</script>

<style lang="stylus">

</style>
