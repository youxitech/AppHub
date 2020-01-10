<template lang="pug">
.dashboard.h-screen.flex.flex-col
  header(class="w-full border-b border-gray-200 h-16 flex items-center p-6")
    a(href="/" class="block lg:mr-4") App Hub
  main(class="w-full flex flex-1")
    nav(class="w-40 border-b border-r")
      .flex.cursor-pointer(
        v-for="app in apps"
        :key="app.id"
        :class="'hover:bg-gray-300 p-3 w-40 ' + (activeNavId === app.id ? 'bg-gray-300' : '') "
        @click="() => onClickNav(app.id)"
      )
        img.w-10.h-10.mr-2.rounded(:src="_getAsset('icon', app.platform, app.bundleID)")
        div
          div {{ app.name }}
          .text-xs.text-gray-500 {{ app.platform }}

    .flex-1.bg-gray-200.p-6(v-if="app")
      .rounded.overflow-hidden.shadow-lg.p-6.flex.flex-col.bg-white.h-full
        .flex
          img.w-24.h-24.mr-2.rounded(:src="_getAsset('icon', app.app.platform, app.app.bundleID)")
          .ml-5
            .font-semibold {{ app.app.name }}
            .text-sm.text-gray-600.mt-3 Platform: {{ app.app.platform }}
            .text-sm.text-gray-600 BundleID: {{ app.app.bundleID }}
            .text-sm.text-gray-600 Download: {{ app.app.downloadCount }}
          .flex.ml-auto.items-center
            .mr-5 APP ID: {{ app.app.alias }}
            button(class="mr-8 hover:text-teal-500" @click="$modal.show('changeAppId')") 修改

        table.table-auto.mt-10
          thead
            tr
              th.px-4.py-2 版本
              th.px-4.py-2 下载次数
              th.px-4.py-2 Build
              th.px-4.py-2 更新时间
              th.px-4.py-2 设置
          tbody
            tr(
              v-for="version, index of app.versions"
              :key="version.id"
            )
              td.border.px-4.py-2.text-center
                .flex.items-center.justify-center
                  .rounded-full.h-3.w-3.mr-3(:class="index === 0 ? 'bg-teal-500' : 'bg-white'")
                  a(class="text-blue-500 hover:text-blue-800" :href="`/version/${ version.id }`") {{ version.version }}
              td.border.px-4.py-2.text-center {{ version.downloadCount }}
              td.border.px-4.py-2.text-center {{ version.pacakgeCount }}
              td.border.px-4.py-2.text-center {{ version.updatedAt | formatTime }}
              td.border.px-4.py-2.text-center
                button(class="hover:text-teal-500" @click="() => setDefaultVersion(version.id)") 设为默认版本
                button(class="ml-3 hover:text-teal-500" @click="() => onDeleteVersion(version.id)") 删除

        _button.mt-auto(@click="onDeleteApp") 删除此 APP

  modal(name="changeAppId")
    .flex.items-center.justify-center.w-full.h-full.flex-col.p-6
      _input(v-model="newAppId")
      .flex.mt-10
        button.mr-10(class="mr-8 hover:text-teal-500" @click="changeAppId") Confirm
        button(class="mr-8 hover:text-teal-500" @click="$modal.hide('changeAppId')") Cancel
</template>

<script>
export default {
  data() {
    return {
      apps: [],
      activeNavId: null,
      app: null,
      newAppId: "",
    }
  },

  mounted() {
    this.fetchApps()
  },

  methods: {
    fetchApps() {
      axios.get("/admin/apps")
        .then(res => {
          this.apps = res.data
          this.activeNavId = this.apps[0].id
          return this.fetchApp()
        })
        .catch(_displayError)
    },

    fetchApp() {
      const alias = this.apps.find(i => i.id === this.activeNavId).alias
      return axios.get(`/admin/apps/${ alias }`)
        .then(res => {
          this.app = res.data
        })
        .catch(_displayError)
    },

    onClickNav(id) {
      if(this.activeNavId === id) return
      this.activeNavId = id
      this.fetchApp()
    },

    setDefaultVersion(id) {
      axios.post(`/admin/versions/${ id }/active`)
        .then(res => {
          const idx = this.app.versions.findIndex(i => i.id === id)
          this.app.versions.unshift(this.app.versions[idx])
          this.app.versions.splice(idx + 1, 1)
          this.$notify({
            type: "success",
            text: "Success!",
          })
        })
        .catch(_displayError)
    },

    changeAppId() {
      axios.patch(`/admin/apps/${ this.activeNavId }`, {
        alias: this.newAppId,
      })
        .then(() => {
          this.$notify({
            type: "success",
            text: "Success!",
          })
          this.app.app.alias = this.newAppId
          this.$modal.hide("changeAppId")
        })
        .catch(_displayError)
    },

    onDeleteVersion(id) {
      _showConfirm("确认删除？")
        .then(() => {
          return axios.delete(`/admin/versions/${ id }`)
            .then(() => {
              this.app.versions.splice(this.app.versions.findIndex(i => i.id === id), 1)
            })
            .catch(_displayError)
        })
        .catch(() => {})
    },

    onDeleteApp() {
      _showConfirm("确认删除？")
        .then(() => {
          return axios.delete(`/admin/apps/${ this.activeNavId }`)
            .then(() => {
              this.apps.splice(this.apps.findIndex(i => i.id === this.activeNavId), 1)
              this.activeNavId = this.apps[0].id
              this.fetchApp()
            })
            .catch(_displayError)
        })
        .catch(() => {})
    },
  },
}
</script>

<style lang="stylus">

</style>
