<template lang="pug">
.bg-gray-200.p-6.h-screen.overflow-auto(v-if="app")
  app-info(:app="app.app")

  .text-3xl.text-center.mt-6 {{ channel }}
  .flex.justify-center.mt-8
    table.table-auto
      thead
        tr
          th.px-4.py-2 包名
          th.px-4.py-2 版本
          th.px-4.py-2(v-if="!isIOS") 渠道
          th.px-4.py-2 创建时间
          th.px-4.py-2 大小
          th.px-4.py-2(v-if="isIOS") 包类型
          th.px-4.py-2 二维码
          th.px-4.py-2 操作

      tbody
        tr(
          v-for="pkg of pkgs"
          :key="pkg.package.id"
        )
          td.border.px-4.py-2.text-center
            router-link(
              class="text-blue-500 hover:text-blue-800"
              target="_blank"
              :to="`/pkg/${ pkg.package.id }`"
            ) {{ pkg.package.name }}
          td.border.px-4.py-2.text-center
            router-link(
              class="text-blue-500 hover:text-blue-800"
              target="_blank"
              :to="`/${ app.app.alias }/${ pkg.version.version }`"
            ) {{ pkg.version.version }}
          td.border.px-4.py-2.text-center(v-if="!isIOS")
            div {{ pkg.package.channel }}
          td.border.px-4.py-2.text-center {{ pkg.package.createdAt | formatTime }}
          td.border.px-4.py-2.text-center {{ pkg.package.size | bytesToSize }}
          td.border.px-4.py-2.text-center(v-if="isIOS")
            div {{ pkg.package.iosPackageType }}
            .text-sm(
              v-if="pkg.package.iosPackageType === 'ad-hoc'"
            )
              span 含有
              span.text-blue-500.cursor-pointer(
                class="hover:underline"
                @click="curIOSDeivceList = pkg.package.iosDeviceList, $modal.show('iosDeviceList')"
              ) {{ pkg.package.iosDeviceList.length }}
              span 台设备
          td.border.px-4.py-2.text-center
            img.w-20.h-20.version__hover(:src="pkg.qrcode")
          td.border.px-4.py-2.text-center
            a(:href="_getAsset('bundle', app.app.platform, app.app.bundleID, pkg.version.version, pkg.package.id)") 下载

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

// route params: /:id/channel/:channel
//  id: app alias
//  channel: channel name
export default {
  data() {
    return {
      pkgs: [],
      curIOSDeivceList: null,
      app: null,
    }
  },

  computed: {
    channel() {
      return this.$route.params.channel
    },

    isIOS() {
      return this.app.app.platform === "ios"
    },
  },

  mounted() {
    this.fetchChannel()
    this.fetchApp()
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

    fetchChannel() {
      return axios.get(`/apps/${ this.$route.params.id }/channels/${ this.channel }`)
        .then(res => {
          const v = res.data.content
          return Promise.all(v.map(pkg => _util.idToQRCode(pkg.id)))
            .then(res => {
              v.forEach((item, index) => {
                item.qrcode = res[index]
              })
              this.pkgs = v
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
