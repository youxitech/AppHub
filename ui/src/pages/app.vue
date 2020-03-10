<template lang="pug">
.app(v-if="app")
  .app__header
    .app__header-wrap
      img.app__header-img(
        :src="_getAsset('icon', app.app.platform, app.app.bundleID)"
      )
      .app__header-name {{ app.app.name }}

      .app__tabs
        .app__tab(
          v-if="app.iosAlias !== ''"
          :class="{ 'app__tab--active': app.app.platform === 'ios' }"
          @click="$router.replace(`/${ app.iosAlias }`)"
        ) iOS
        .app__tab(
          v-if="app.androidAlias !== ''"
          :class="{ 'app__tab--active': app.app.platform === 'android' }"
          @click="$router.replace(`/${ app.androidAlias }`)"
        ) Android

  .app__body
    .app__main
      .app__top
        span.app__platform {{ PLATFORM[app.app.platform] }}
        span.app__bundle {{ app.app.bundleID }}

      .app__filter
        .app__label 发布环境：
        .app__options
          .app__option(
            :class="{ 'app__option--selected': curEnv == null }"
            @click="setEnv(null)"
          ) 全部
          .app__option(
            v-for="item in app.envs"
            :key="item"
            :class="{ 'app__option--selected': curEnv === item }"
            @click="setEnv(item)"
          ) {{ item }}

      .app__filter
        .app__label 发布渠道：
        .app__options
          .app__option(
            :class="{ 'app__option--selected': curChannel == null }"
            @click="setChannel(null)"
          ) 全部
          .app__option(
            v-for="item in app.channels"
            :key="item"
            :class="{ 'app__option--selected': curChannel === item }"
            @click="setChannel(item)"
          ) {{ item }}

      .app__pkg-list
        .app__pkg--none(v-if="pkgs.length === 0 && !isLoading")
          img.app__pkg-img--none(src="/static/pkg-none.png")
          .app__pkg-text--none 没有可用的安装包
          .app__pkg-text--none 更改一下筛选条件吧

        .app__pkg(
          v-for="item in pkgs"
          :key="item.id"
        )
          router-link.app__pkg-name(
            class="hover:underline"
            :to="`/pkg/${ item.id }`"
          ) {{ item.name }}
          .app__pkg-body
            .app__pkg-info
              div 创建时间：{{ item.createdAt | formatTime }}
              div 大小：{{ item.size | bytesToSize }}

            .app__tag-wrap
              .app__tag(
                @click="setVersion(app.versions.find(version => version.version === item.version).id)"
              ) {{ item.version }}
            .app__tag-wrap
              .app__tag(@click="setEnv(item.env)") {{ item.env }}
            .app__tag-wrap
              .app__tag(@click="setChannel(item.channel)") {{ item.channel }}

            .app__qrcode-wrap(ref="popper")
              img.app__qrcode(
                v-if="item.qrcode != null"
                :src="item.qrcode"
              )
            div(
              v-show="false"
              ref="display"
            )
              img(:src="item.qrcode")

            a.app__download-wrap(
              :href="_getAsset('bundle', app.app.platform, app.app.bundleID, item.version, item.id)"
            )
              img.app__download(src="/static/download.svg")

    .app__right
      .app__version(
        :class="{ 'app__version--active': null === curVersion }"
        @click="setVersion(null)"
      )
        .app__version-text 全部
        .app__round(v-if="null === curVersion")
      .app__version(
        v-for="version in app.versions"
        :key="version.id"
        :class="{ 'app__version--active': version.id === curVersion }"
        @click="setVersion(version.id)"
      )
        .app__version-text {{ version.version }}
        .app__round(v-if="version.id === curVersion")
</template>

<script>
import QRCode from "qrcode"
import tippy from "tippy.js"
import "tippy.js/dist/tippy.css"
import "tippy.js/themes/light.css"

const PLATFORM = {
  android: "Android",
  ios: "iOS",
}

export default {
  data() {
    return {
      curEnv: null,
      curChannel: null,
      curVersion: null,
      pkgs: [],
      isLoading: false,
      PLATFORM,
    }
  },

  computed: {
    app() {
      return this.$store.state.app
    },
  },

  beforeRouteUpdate(to, from, next) {
    next()
    this.curEnv = this.$route.query.env || null
    this.curChannel = this.$route.query.channel || null
    this.curVersion = Number(this.$route.query.version) || null
    this.pkgs = []
    this.$store.dispatch("getAppInfo", this.$route.params.id)
      .then(() => this.getApp())
  },

  mounted() {
    this.curEnv = this.$route.query.env || null
    this.curChannel = this.$route.query.channel || null
    this.curVersion = Number(this.$route.query.version) || null

    this.$store.dispatch("getAppInfo", this.$route.params.id)
      .then(() => this.getApp())
  },

  methods: {
    getApp() {
      const params = {}
      if(this.curEnv != null) {
        params.env = this.env
      }
      if(this.curChannel != null) {
        params.channel = this.curChannel
      }
      if(this.curVersion != null) {
        params.versionID = this.curVersion
      }

      this.pkgs.forEach(item => {
        item.tippyInstance.destroy()
      })
      this.pkgs = []
      this.isLoading = true

      axios.get(`/apps/${ this.$route.params.id }/packages`, { params })
        .then(res => {
          this.pkgs = res.data.map(item => ({
            ...item,
            version: this.app.versions
              .find(version => version.id === item.versionID)
              .version,
          }))
        })
        .then(() => {
          this.pkgs.map((pkg, index) => {
            QRCode.toDataURL(`${ location.origin }/pkg/${ pkg.id }`)
              .then(url => {
                this.$set(pkg, "qrcode", url)
                this.$nextTick(() => {
                  pkg.tippyInstance = tippy(this.$refs.popper[index], {
                    content: this.$refs.display[index].innerHTML,
                    allowHTML: true,
                    theme: "light",
                  })
                })
              })
          })
          this.isLoading = false
        })
        .catch(_displayError)
    },

    setEnv(env) {
      this.curEnv = env
      this.$router.replace({
        path: this.$route.path,
        query: {
          ...this.$route.query,
          env,
        },
      })
      this.getApp()
    },

    setChannel(channel) {
      this.curChannel = channel
      this.$router.replace({
        path: this.$route.path,
        query: {
          ...this.$route.query,
          channel,
        },
      })
      this.getApp()
    },

    setVersion(version) {
      this.curVersion = version
      this.$router.replace({
        path: this.$route.path,
        query: {
          ...this.$route.query,
          version,
        },
      })
      this.getApp()
    },
  },
}
</script>

<style lang="stylus">
.app
  height: 100vh
  width: 100vw
  flex: 1
  flex-direction: column
  display: flex
  background: $neutral-1
  overflow: auto
  overflow-x: auto

.app__header
  height: 64px
  background: white
  border-bottom: 1px solid rgba(16, 22, 26, 0.2)
  box-shadow: 0 2px 6px 0 rgba(16, 22, 26, 0.2)
  padding: 0 40px
  min-width: 1380px
  box-sizing: content-box

.app__header-wrap
  height: 100%
  display: flex
  align-items: center
  width: 1380px
  margin: 0 auto

.app__header-img
  width: 32px
  height: 32px
  border-radius: 2px
  margin-right: 16px

.app__header-name
  font-size: 20px
  color: $neutral-10

.app__tabs
  height: 100%
  margin-left: auto
  display: flex

.app__tab
  display: flex
  padding: 0 20px
  align-items: center
  cursor: pointer
  color: $neutral-7
  font-size: 14px

.app__tab--active
  color: $primary-6
  border-top: 2px solid $primary-6
  font-weight: bold

.app__body
  display: flex
  flex: 1
  padding: 0 40px
  width: 1380px
  margin: 0 auto
  box-sizing: content-box

.app__main
  flex: 1
  display: flex
  flex-direction: column
  margin-right: 150px

.app__top
  margin: 30px 0

.app__platform
  color: $neutral-10
  font-size: 28px
  font-weight: bold
  margin-right: 20px

.app__bundle
  color: $neutral-6
  font-size: 14px

.app__pkg-list
  margin-top: 30px
  flex: 1
  overflow: auto

.app__pkg
  background: white
  border-radius: 4px
  padding: 20px 20px 30px
  margin-bottom: 20px
  border: 1px solid rgba(16, 22, 26, 0.2)
  max-width: 958px
  box-sizing: border-box

.app__pkg-name
  font-size: 16px
  color: $neutral-10
  cursor: pointer
  font-weight: bold

.app__pkg-info
  width: 240px
  margin-right: 100px
  color: $neutral-9
  font-size: 14px
  line-height: 1

.app__pkg-info > div:last-child
  margin-top: 10px

.app__pkg-body
  margin-top: 10px
  display: flex

.app__tag-wrap
  flex: 1
  margin-right: 20px

.app__qrcode-wrap
  width: 32px
  height: 32px
  margin-right: 32px
  cursor: pointer
  align-self: center

.app__qrcode
  width: 100%
  height: 100%
  align-self: center

.app__download-wrap
  width: 32px
  height: 32px
  cursor: pointer
  align-self: center

.app__download
  margin: 6px
  width: 20px
  height: 20px

.app__round
  position: absolute
  width: 10px
  height: 10px
  box-sizing: border-box
  border: 3px solid #137CBD
  background: white
  border-radius: 999px
  position: absolute
  left: 0
  top: 50%
  transform: translate(-64%, -50%)

.app__version-text
  cursor: pointer

.app__version-text:hover
  text-decoration: underline
  color: $info-6

.app__filter
  display: flex

.app__label
  color: $neutral-10
  font-size: 14px
  flex-shrink: 0

.app__options
  display: flex
  flex-wrap: wrap

.app__option
  font-size: 14px
  color: $neutral-7
  padding: 0 16px
  height: 24px
  margin: 0 16px 16px
  display: flex
  align-items: center
  cursor: pointer

.app__option--selected
  background: $info-2
  color: $info-6
  border-radius: 4px
  font-weight: bold

.app__tag
  background: $neutral-2
  font-size: 14px
  color: $neutral-9
  height: 30px
  padding: 0 14px
  border-radius: 3px
  display: inline-flex
  align-items: center
  justify-content: center
  cursor: pointer

.app__right
  width: 250px
  padding: 8px 0
  border-left: 2px solid $neutral-3
  margin-top: 68px
  align-self: start
  box-sizing: border-box

.app__version
  height: 22px
  font-size: 14px
  margin: 12px 0
  padding: 0 16px
  position: relative
  display: flex
  align-items: center
  color: $neutral-9

.app__version--active
  color: $info-6
  font-weight: bold

.app__pkg--none
  height: 400px
  display: flex
  flex-direction: column
  justify-content: center
  align-items: center

.app__pkg-text--none
  color: $neutral-6
  font-size: 14px
  line-height: 20px
  margin-top: 10px

.app__pkg-text--none:last-child
  margin-top: 0
</style>
