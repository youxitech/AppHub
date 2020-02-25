<template lang="pug">
.pkg.h-screen(v-if="pkg")
  .pkg__main.flex.items-center.justify-center.flex-col
    .flex.items-center.justify-center
      img.rounded(
        :src="_getAsset('icon', pkg.app.platform, pkg.app.bundleID)"
        style="height: 44px; width: 44px;"
      )
      div(style="font-size: 32px; margin-left: 20px") {{ pkg.app.name }}

    .pkg__card.mt-5.flex
      .pkg__card--lg.pr-12(class="hidden lg:flex items-center justify-center flex-col")
        img(:src="this.qrcode" style="width: 128px; height: 128px")
        a(
          :href="_getAsset('bundle', pkg.app.platform, pkg.app.bundleID, pkg.version.version, pkg.package.id)"
          :download="pkg.package.name"
        )
          _button.w-40.mt-6 下载

      div(class="lg:ml-12")
        p.mb-2.font-bold(style="font-size: 24px") {{ pkg.app.platform }}
        p.mb-2.font-bold(style="font-size: 16px") {{ pkg.package.name }}
        p 版本: {{ pkg.version.version }}
        p 渠道: {{ pkg.package.channel }}
        p 时间: {{ pkg.package.createdAt | formatTime }}
        p 大小: {{ pkg.package.size | bytesToSize }}
        _button.w-full.mt-6(v-if="global.isIos" class="lg:hidden" @click="install") 安装
        a(
          v-else
          :href="_getAsset('bundle', pkg.app.platform, pkg.app.bundleID, pkg.version.version, pkg.package.id)"
          :download="pkg.package.name"
          class="lg:hidden"
        )
          _button.w-full.mt-6 下载
  _copyright
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
      const url = `itms-services://?action=download-manifest&url=${ location.origin }/api/plist/${ this.pkg.package.id }`
      const a = document.createElement("a")
      a.setAttribute("href", url)
      a.click()
    },
  },
}
</script>

<style lang="stylus">
.pkg
  background: $background

.pkg__main
  height: calc(100vh - 80px)

.pkg__card
  background: #fff
  box-shadow: 0px 1px 1px 0px rgba(16,22,26,0.2),0px 0px 0px 1px rgba(16,22,26,0.1)
  border-radius: 28px
  padding: 50px
  p
    margin-bottom: 10px
    color: $text-gray

.pkg__card--lg
  border-right: 1px solid $border
</style>
