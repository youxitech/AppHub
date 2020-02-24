<template lang="pug">
.text-center.flex-1.mt-8
  template(v-if="version && app")
    .text-3xl {{ app && app.app.name }} {{ version.version.version}}
    .text-sm.mt-2.text-gray-600 bundleID: {{ app && app.app.bundleID }}

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
export default {
  data() {
    return {
      app: null,
      version: null,
      curIOSDeivceList: null,
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
      axios.get(`/apps/${ this.$route.params.id }/${ this.$route.params.version }`)
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
}
</script>

<style lang="stylus">
.version__hover:hover
  transform: scale(2)
  transition: 0.5s

.version__modal
  padding: 12px
  background: white
  height: 400px !important
</style>
