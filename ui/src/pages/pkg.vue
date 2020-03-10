<template lang="pug">
.pkg(v-if="pkg")
  .pkg__top
    .pkg__top-main
      img.pkg__top-img(
        :src="_getAsset('icon', pkg.app.platform, pkg.app.bundleID)"
      )
      .pkg__top-title {{ pkg.app.name }}

  .pkg__box
    .pkg__box-left
      img(:src="qrcode")
      a.pkg__pc-download-btn-wrap(
        :href="_getAsset('bundle', pkg.app.platform, pkg.app.bundleID, pkg.version.version, pkg.package.id)"
      )
        button 下载

    .pkg__box-sep

    div
      .pkg__title {{ PLATFORM[pkg.app.platform] }}
      .pkg__name {{ pkg.package.name }}
      .pkg__info 版本：{{ pkg.version.version }}
      .pkg__info 环境：{{ pkg.package.env }}
      .pkg__info 渠道：{{ pkg.package.channel }}
      .pkg__info 时间：{{ pkg.package.createdAt | formatTime }}
      .pkg__info 大小：{{ pkg.package.size | bytesToSize }}
      a.pkg__mobile-btn-wrap(
        :href="_getAsset('bundle', pkg.app.platform, pkg.app.bundleID, pkg.version.version, pkg.package.id)"
      )
        button 下载

  .pkg__footer
    .pkg__footer-text
      span Powered by
      a.pkg__footer-link(
        href="https://www.youxishequ.com"
        target="_blank"
      ) &nbsp;Youxishequ.com
    .pkg__footer-text
      span Fork us on
      a.pkg__footer-link(
        href="https://github.com/youxitech/AppHub"
        target="_blank"
      ) &nbsp;Github
</template>

<script>
import QRCode from "qrcode"

const PLATFORM = {
  android: "Android",
  ios: "iOS",
}

export default {
  data() {
    return {
      pkg: null,
      qrcode: "",
      PLATFORM,
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
  width: 100vw
  min-height: 100vh
  display: flex
  flex-direction: column
  align-items: center
  background: $neutral-1
  overflow: auto

.pkg__top
  margin-top: 12vh

.pkg__top-main
  display: flex
  justify-content: center
  align-items: center

.pkg__top-img
  width: 44px
  height: 44px
  border-radius: 2px

.pkg__top-title
  font-size: 32px
  color: $neutral-10
  margin-left: 20px

.pkg__info
  color: $neutral-9
  font-size: 14px
  margin: 10px 0
  line-height: 1

.pkg__mobile-btn-wrap
  display: block
  margin-top: 40px
  width: 100%
  height: 40px

  button
    width: 100%
    height: 100%
    border-radius: 3px
    background: rgba(19, 124, 189, 1)
    box-shadow: 0 0 0 1px rgba(16, 22, 26, 0.4), 0 -1px 1px 0 rgba(16, 22, 26, 0.2)
    color: white

.pkg__box
  border: 1px solid rgba(16, 22, 26, 0.2)
  background: white
  border-radius: 28px
  padding: 12vw
  margin: 5vh 14vw
  word-break: break-all
  max-width: 74vw

.pkg__box-left
  display: none

.pkg__box-sep
  display: none

.pkg__title
  color: $neutral-10
  font-size: 24px
  font-weight: bold

.pkg__name
  color: $neutral-10
  font-size: 16px
  margin-top: 18px
  font-weight: bold

.pkg__footer
  margin-top: auto
  height: 12vh
  box-sizing: content-box
  text-align: center
  width: 100%
  display: flex
  flex-direction: column
  justify-content: center

.pkg__footer-text
  text-align: center
  color: $neutral-7
  font-size: 14px

.pkg__footer-link
  color: $info-6

@media only screen and (min-width: 640px)
  .pkg__top
    margin-top: 18vh

  .pkg__box
    padding: 60px 90px
    display: flex

  .pkg__box-left
    display: flex
    width: 200px
    flex-shrink: 0
    flex-direction: column
    height: 100%
    align-self: flex-start

    img
      width: 128px
      height: 128px
      margin: 0 auto

  .pkg__pc-download-btn-wrap
    display: block
    width: 100%
    height: 40px
    margin-top: auto

    button
      width: 100%
      height: 100%
      border-radius: 3px
      background: rgba(19, 124, 189, 1)
      box-shadow: 0 0 0 1px rgba(16, 22, 26, 0.4), 0 -1px 1px 0 rgba(16, 22, 26, 0.2)
      color: white

  .pkg__box-sep
    display: block
    width: 1px
    background: $neutral-3
    height: 100%
    margin: 0 60px

  .pkg__mobile-btn-wrap
    display: none

  .pkg__footer
    flex-direction: row
    align-items: center

  .pkg__footer-text
    display: inline

  .pkg__footer-text:first-child
    margin-right: 8px
</style>
