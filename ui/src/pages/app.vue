<template lang="pug">
.flex-1.h-screen.text.flex.flex-col(v-if="app")
  .flex.h-16.bg-white.px-8.items-center
    img.w-10.h-10.mr-4(:src="_getAsset('icon', app.app.platform, app.app.bundleID)")
    .text-lg {{ app.app.name }}

    .flex.ml-auto.h-full
      .flex.px-4.text-base.text-gray-500.items-center.cursor-pointer(
        v-if="app.iosAlias !== ''"
        :class="{ 'text-blue-500 border-blue-500 border-t-2': app.app.platform === 'ios' }"
        @click="$router.replace(`/${ app.iosAlias }`)"
      ) iOS
      .flex.px-4.text-base.text-gray-500.items-center.cursor-pointer(
        v-if="app.androidAlias !== ''"
        :class="{ 'text-blue-500 border-blue-500 border-t-2': app.app.platform === 'android' }"
        @click="$router.replace(`/${ app.androidAlias }`)"
      ) Android

  .flex.px-8.flex-1.bg-gray-200
    .flex-1.flex.flex-col.mr-32
      .mt-8.mb-4
        span.text-2xl.mr-4 {{ app.app.platform }}
        span.text-sm.text-gray-600 {{ app.app.bundleID }}

      .app__filter
        .app__label 发布环境：
        .app__options
          .app__option(
            :class="{ 'app__option--selected': curEnv == null }"
            @click="curEnv = null, getApp()"
          ) 全部
          .app__option(
            v-for="item in app.envs"
            :key="item"
            :class="{ 'app__option--selected': curEnv === item }"
            @click="curEnv = item, getApp()"
          ) {{ item }}

      .app__filter
        .app__label 发布环境：
        .app__options
          .app__option(
            :class="{ 'app__option--selected': curChannel == null }"
            @click="curChannel = null, getApp()"
          ) 全部
          .app__option(
            v-for="item in app.channels"
            :key="item"
            :class="{ 'app__option--selected': curChannel === item }"
            @click="curChannel = item, getApp()"
          ) {{ item }}

      .flex-1.overflow-auto.mt-4
        .bg-white.rounded-sm.p-4.mb-4(
          v-for="item in pkgs"
        )
          router-link.text-base.font-bold.text-gray-900.cursor-pointer(
            class="hover:underline"
            :to="`/pkg/${ item.id }`"
          ) {{ item.name }}
          .mt-2.flex.items-center
            .w-56.mr-24
              .text-sm.text-gray-600 创建时间：{{ item.createdAt | formatTime }}
              .text-sm.text-gray-600 大小：{{ item.size | bytesToSize }} MB

            .flex-1.mr-4
              .app__tag {{ item.version }}
            .flex-1.mr-4
              .app__tag {{ item.env }}
            .flex-1.mr-4
              .app__tag {{ item.channel }}

            .w-10.h-10.mr-4.cursor-pointer(ref="popper")
              img.w-full.h-full(
                v-if="item.qrcode != null"
                :src="item.qrcode"
              )
            div(
              v-show="false"
              ref="display"
            )
              img(:src="item.qrcode")

            a.w-6.h-6.cursor-pointer(
              :href="_getAsset('bundle', app.app.platform, app.app.bundleID, item.version, item.id)"
            )
              img.w-full.h-full(src="/static/download.svg")

    .w-40.py-4.border-gray-400.border-l-2.mt-20.self-start.mr-20
      .h-8.text-sm.px-4.relative.flex.items-center(
        :class="{ 'text-blue-400': null === curVersion }"
        @click="curVersion = null, getApp()"
      )
        .cursor-pointer.app__version 全部
        .app__round(v-if="null === curVersion")
      .h-8.text-sm.px-4.relative.flex.items-center(
        v-for="version in app.versions"
        :key="version.id"
        :class="{ 'text-blue-400': version.id === curVersion }"
        @click="curVersion = version.id, getApp()"
      )
        .cursor-pointer.app__version {{ version.version }}
        .app__round(v-if="version.id === curVersion")
</template>

<script>
import QRCode from "qrcode"
import tippy from "tippy.js"
import "tippy.js/dist/tippy.css"
import "tippy.js/themes/light.css"

export default {
  data() {
    return {
      curEnv: null,
      curChannel: null,
      curVersion: null,
      pkgs: [],
    }
  },

  computed: {
    app() {
      return this.$store.state.app
    },
  },

  beforeRouteUpdate(to, from, next) {
    next()
    this.curEnv = null
    this.curChannel = null
    this.curVersion = null
    this.pkgs = []
    this.$store.dispatch("getAppInfo", this.$route.params.id)
      .then(() => this.getApp())
  },

  mounted() {
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
                  tippy(this.$refs.popper[index], {
                    content: this.$refs.display[index].innerHTML,
                    allowHTML: true,
                    theme: "light",
                  })
                })
              })
          })
        })
        .catch(_displayError)
    },
  },
}
</script>

<style lang="stylus">
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
  transform: translate(-60%, -50%)

.app__version:hover
  @apply underline text-blue-400

.app__filter
  @apply flex

.app__label
  @apply text-sm text-gray-700 flex-shrink-0

.app__options
  @apply flex flex-wrap

.app__option
  @apply text-sm text-gray-500 px-4 h-6 mx-4 mb-4 cursor-pointer

.app__option--selected
  @apply text-blue-600 bg-blue-200 rounded-sm

.app__tag
  @apply bg-gray-200 text-sm text-gray-600 h-6 px-2 rounded-sm flex items-center inline-block
</style>
