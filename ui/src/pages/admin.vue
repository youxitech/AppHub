<template lang="pug">
.admin.h-screen.flex.flex-col
  header(class="w-full border-b border-gray-200 h-16 flex items-center p-6")
    a(href="/" class="block lg:mr-4") App Hub
    _button(@click="showUploader") +
  main(class="w-full flex flex-1")
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

  modal(name="uploader")
    .w-64.h-64.flex.flex-col.items-center
      .uploader
      .uploader__progress-bar.w-64
</template>

<script>
import Uppy from "@uppy/core"
import DragDrop from "@uppy/drag-drop"
import ProgressBar from "@uppy/progress-bar"
import XHRUpload from "@uppy/xhr-upload"

export default {
  data() {
    return {
      apps: [],
    }
  },

  mounted() {
    this.fetchApps()
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
          if(this.$route.params.id) return
          this.$router.push(`/admin/${ this.apps[0].alias }`)
        })
        .catch(_displayError)
    },

    showUploader() {
      this.$modal.show("uploader")
      // @FIX
      setTimeout(() => {
        const uppy = new Uppy({ debug: true, autoProceed: true })
        uppy
          .use(DragDrop, { target: ".uploader" })
          .use(ProgressBar, { target: ".uploader__progress-bar", hideAfterFinish: false })
          .use(XHRUpload, {
            endpoint: "/api/admin/upload",
            formData: true,
            fieldName: "file",
            headers: {
              "X-Admin-Token": _db.token,
            },
          })
          .on("upload-success", (file, response) => {
            this.$modal.hide("uploader")
            this.$notify({
              type: "success",
              text: "Success!",
            })
            this.$router.push(`/admin/${ response.body.app.alias }/version/${ response.body.version.id }?t=${ new Date().getTime() }`)
          })
          .on("upload-error", (file, error, response) => {
            _displayError(response.body.msg)
            this.$modal.hide("uploader")
          })
      }, 100)
    },

    onClickNav(alias) {
      this.$router.push(`/admin/${ alias }`)
    },
  },
}
</script>

<style lang="stylus">

</style>
