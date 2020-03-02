<template lang="pug">
.w-screen.h-screen.flex.flex-col.items-center.bg-gray-200.overflow-auto(
  v-if="pkg"
)
  .mt-20(class="md:mt-20")
    .flex.justify-center.items-center
      img.w-12.h-12.rounded(
        :src="_getAsset('icon', pkg.app.platform, pkg.app.bundleID)"
      )
      .text-2xl.ml-4 {{ pkg.app.name }}

  .pkg__box(
    class="sm:py-12 sm:px-16 sm:flex"
  )
    .hidden(class="sm:block sm:w-40 sm:flex-shrink-0")
      img(:src="qrcode")
      a.block.w-full(
        :href="_getAsset('bundle', pkg.app.platform, pkg.app.bundleID, pkg.version.version, pkg.package.id)"
      )
        button.w-full.rounded-sm.mt-12.bg-blue-600.text-white.h-10(
        ) 下载

    .hidden(class="sm:block sm:w-px sm:bg-gray-400 sm:mx-8")

    div
      .text-xl.font-bold {{ pkg.app.platform }}
      .text-base.font-bold.mt-8 {{ pkg.package.name }}
      .pkg__info 版本：{{ pkg.version.version }}
      .pkg__info 环境：{{ pkg.package.env }}
      .pkg__info 渠道：{{ pkg.package.channel }}
      .pkg__info 时间：{{ pkg.package.createdAt | formatTime }}
      .pkg__info 大小：{{ pkg.package.size | bytesToSize }}
      a.block.w-full(
        class="sm:hidden"
        :href="_getAsset('bundle', pkg.app.platform, pkg.app.bundleID, pkg.version.version, pkg.package.id)"
      )
        button.w-full.rounded-sm.mt-12.bg-blue-600.text-white.h-10(
        ) 下载

  .mt-auto.mb-4
    .text-center.text-gray-600.leading-normal.text-sm(
      class="sm:inline"
    )
      span Powered by
      a.text-blue-600(href="https://www.youxishequ.com") &nbsp;Youxishequ.com
    .text-center.text-gray-600.leading-normal.text-sm(
      class="sm:inline sm:ml-2"
    )
      span Fork us on
      a.text-blue-600(href="https://github.com/youxitech/AppHub") &nbsp;Github
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
.pkg__qrcode
  background-image: url("/static/phone.png")
  background-size: cover
  width 300px
  height 450px

.pkg__info
  @apply text-gray-600 text-sm my-2

.pkg__box
  max-width: 80%
  @apply border border-gray-400 mt-8 bg-white \
    rounded-lg px-8 py-8 mx-6 mb-8 break-all
</style>
