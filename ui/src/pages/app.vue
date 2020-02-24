<template lang="pug">
.flex-1.bg-gray-200.p-6.h-screen(v-if="app && version")
  app-info(:app="app.app")

  .flex.py-8
    .w-40.bg-white.rounded.shadow-lg.h-auto.overflow-auto.self-start
      .h-10.flex.items-center.justify-center.border-gray-200.border-solid.border-b.cursor-pointer(
        v-for="item in app.versions"
        :key="item.id"
        :class="{ 'bg-blue-400 text-white': item.version === curVersion }"
        @click="curVersion = item.version"
      ) {{ item.version }}

    .text-center.flex-1
      router-link.text-3xl.text-blue-500(
        :to="'/' + app.app.alias + '/' + version.version.version"
        class="hover:underline"
      ) {{ version.version.version }}

      .flex.justify-center.mt-8
        table.table-auto
          thead
            tr
              th.px-4.py-2 包名
              th.px-4.py-2(v-if="!isIOS") 渠道
              th.px-4.py-2 创建时间
              th.px-4.py-2 大小
              th.px-4.py-2(v-if="isIOS") 包类型
              th.px-4.py-2 二维码
              th.px-4.py-2 操作

          tbody
            tr(
              v-for="pkg of version.packages"
              :key="pkg.id"
            )
              td.border.px-4.py-2.text-center
                router-link(
                  class="text-blue-500 hover:text-blue-800"
                  target="_blank"
                  :to="`/pkg/${ pkg.id }`"
                ) {{ pkg.name }}
              td.border.px-4.py-2.text-center(v-if="!isIOS")
                router-link(
                  class="text-blue-500 hover:text-blue-800"
                  target="_blank"
                  :to="`/${ app.app.alias }/channel/${ pkg.channel }`"
                ) {{ pkg.channel }}
              td.border.px-4.py-2.text-center {{ pkg.createdAt | formatTime }}
              td.border.px-4.py-2.text-center {{ pkg.size | bytesToSize }}
              td.border.px-4.py-2.text-center(v-if="isIOS")
                div {{ pkg.iosPackageType }}
                .text-sm(
                  v-if="pkg.iosPackageType === 'ad-hoc'"
                )
                  span 含有
                  span.text-blue-500.cursor-pointer(
                    class="hover:underline"
                    @click="curIOSDeivceList = pkg.iosDeviceList, $modal.show('iosDeviceList')"
                  ) {{ pkg.iosDeviceList.length }}
                  span 台设备
              td.border.px-4.py-2.text-center
                img.w-20.h-20.version__hover(:src="pkg.qrcode")
              td.border.px-4.py-2.text-center
                a(:href="_getAsset('bundle', app.app.platform, app.app.bundleID, version.version.version, pkg.id)") 下载

  modal(name="iosDeviceList" classes="version__modal")
    .flex.flex-col.h-full
      .text-center.text-2xl ios设备列表
      .mt-4.overflow-auto.flex-1
        .text-left.text-red-400.leading-loose(
          v-for="item in curIOSDeivceList"
          :key="item"
        ) {{ item }}

      button.mt-6.border.w-24.rounded.h-8.bg-blue-700.text-white.self-center(
        @click="$modal.hide('iosDeviceList')"
      ) 确定
</template>

<script>
import AppInfo from "app-info"

export default {
  data() {
    return {
      curIOSDeivceList: null,
      curVersion: null,
      app: null,
      version: null,
    }
  },

  computed: {
    isIOS() {
      return this.app.app.platform === "ios"
    },
  },

  mounted() {
    this.fetchApp()
      .then(this.fetchVersion)
  },

  watch: {
    curVersion() {
      this.fetchVersion()
    },
  },

  methods: {
    fetchApp() {
      return axios.get(`/apps/${ this.$route.params.id }`)
        .then(res => {
          this.app = res.data
          this.curVersion = this.app.versions[0].version
        })
        .catch(_displayError)
    },

    fetchVersion() {
      axios.get(`/apps/${ this.$route.params.id }/${ this.curVersion }`)
        .then(res => {
          const v = res.data
          return Promise.all(v.packages.map(pkg => _util.idToQRCode(pkg.id)))
            .then(res => {
              v.packages.forEach((item, index) => {
                item.qrcode = res[index]
              })
              this.version = v
            })
        })
        .catch(_displayError)
    },
  },

  components: {
    AppInfo,
  },
}
</script>

<style lang="stylus">

</style>
