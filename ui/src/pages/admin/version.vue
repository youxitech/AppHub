<template lang="pug">
.admin-version.text-center(v-if="version != null && app != null")
  .text-3xl
    span {{ app && app.app.name }}
    router-link.text-blue-500(
      class="hover:underline"
      :to="`/${ $route.params.id }/${ $route.params.version }`"
      target="_blank"
    ) {{ version.version.version}}
  .text-sm.mt-2.text-gray-600 bundleID: {{ app && app.app.bundleID }}

  div
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
            th.px-4.py-2(v-if="!isIOS") 操作

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
              div {{ pkg.channel }}
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
            td.border.px-4.py-2.text-center(v-if="!isIOS")
              button(
                class="hover:text-teal-500"
                @click="changeChannel(pkg)"
              ) 修改渠道

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

    modal(name="newChannelName" @opened="afterModalOpen")
      .flex.items-center.justify-center.w-full.h-full.flex-col.p-6
        .text-2xl.mb-8 修改 {{ changeChannelPkg.name }} 渠道名
        _input(
          v-model="newChannelName"
          ref="channelInput"
          @enter="changeChannelName"
        )
        .flex.mt-10
          button.mr-10(class="mr-8 hover:text-teal-500" @click="changeChannelName") Confirm
          button(class="mr-8 hover:text-teal-500" @click="$modal.hide('newChannelName')") Cancel
</template>

<script>
export default {
  data() {
    return {
      version: null,
      app: null,
      curIOSDeivceList: [],
      newChannelName: "",
      changeChannelPkg: {},
    }
  },

  computed: {
    isIOS() {
      return this.app.app.platform === "ios"
    },
  },

  mounted() {
    this.fetchApp()
    this.fetchVersion()
  },

  methods: {
    fetchApp() {
      return axios.get(`/admin/apps/${ this.$route.params.id }`)
        .then(res => {
          this.app = res.data
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

    changeChannelName() {
      if(this.newChannelName === "") {
        _showErr("请输入渠道名")
        return
      }

      axios.patch(`/admin/packages/${ this.changeChannelPkg.id }`, {
        channel: this.newChannelName,
      })
        .then(() => {
          this.$modal.hide("newChannelName")
          this.changeChannelPkg.channel = this.newChannelName
          _showSuccess("修改成功")
        })
        .catch(_displayError)
    },

    changeChannel(pkg) {
      this.changeChannelPkg = pkg
      this.newChannelName = pkg.channel
      this.$modal.show("newChannelName")
    },

    afterModalOpen() {
      this.$nextTick(() => {
        this.$refs.channelInput.focus()
      })
    },
  },
}
</script>

<style lang="stylus">

</style>
