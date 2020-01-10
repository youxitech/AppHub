<template lang="pug">
.pkg.w-screen.h-screen.flex.justify-center.items-center(v-if="pkg")
  .mr-48
    img.w-24.h-24.rounded(:src="_getAsset('icon', pkg.app.platform, pkg.app.bundleID)")
    .mt-5.text-2xl {{ pkg.app.name }}
    .mt-3.leading-loose.text-gray-500.text-sm
      p 版本: {{ pkg.version.version }} 大小: {{ pkg.package.size | bytesToSize }}
      p 发布日期: {{ pkg.package.createdAt | formatTime }}
  .pkg__qrcode
    img.mt-16.w-40.h-40.ml-6(:src="this.qrcode")
    .mt-12.ml-6.text-gray-500 请扫描二维码下载APP
    .ml-6.text-gray-500 适用于 {{ pkg.app.platform }} 系统
    _button.ml-6(v-if="global.isIos" @click="install") 安装
    a.ml-6(
      v-else
      :href="_getAsset('bundle', pkg.app.platform, pkg.app.bundleID, pkg.version.version, pkg.package.id)"
      :download="pkg.package.name"
    )
      _button 下载
</template>

<script>
import QRCode from "qrcode"

export default {
  data() {
    return {
      pkg: null,
      qrcode: "",
    }
  },

  mounted() {
    this.fetchPkg()
    QRCode.toDataURL(location.href)
      .then(url => {
        this.qrcode = url
      })
      .catch(_displayError)
  },

  methods: {
    fetchPkg() {
      return axios.get(`/packages/${ this.$route.params.id }`)
        .then(res => {
          this.pkg = res.data
        })
        .catch(_displayError)
    },

    install() {
      const url = ""
      const a = document.createElement("a")
      a.setAttribute("href", url)
      a.click()
    },
  },
}
</script>

<style lang="stylus">
.pkg
  background-image: url("/static/pkg-bg.png")
  background-size: cover

.pkg__qrcode
  background-image: url("/static/phone.png")
  background-size: cover
  width 300px
  height 450px

@media (max-width: 640px)
  .pkg
    flex-direction: column
</style>
