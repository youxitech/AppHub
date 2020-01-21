<template lang="pug">
div(v-if="pkgs.length > 0 && app != null")
  .flex.justify-center.mt-8
    table.table-auto
      thead
        tr
          th.px-4.py-2 包名
          th.px-4.py-2(v-if="inChannel") 版本
          th.px-4.py-2(v-if="!isIOS") 渠道
          th.px-4.py-2 创建时间
          th.px-4.py-2 大小
          th.px-4.py-2(v-if="isIOS") 包类型
          th.px-4.py-2 二维码
          th.px-4.py-2(v-if="!inAdmin") 操作
          th.px-4.py-2(v-if="inAdmin && !isIOS") 操作

      tbody
        tr(
          v-for="pkg of pkgsWithQRCode"
          :key="pkg.id"
        )
          td.border.px-4.py-2.text-center
            router-link(
              class="text-blue-500 hover:text-blue-800"
              target="_blank"
              :to="`/pkg/${ pkg.id }`"
            ) {{ pkg.name }}
          td.border.px-4.py-2.text-center(v-if="inChannel")
            router-link(
              class="text-blue-500 hover:text-blue-800"
              target="_blank"
              :to="`/${ id }/${ pkg.version }`"
            ) {{ pkg.version }}
          td.border.px-4.py-2.text-center(v-if="!isIOS")
            router-link(
              v-if="!inChannel"
              class="text-blue-500 hover:text-blue-800"
              target="_blank"
              :to="`/${ id }/channel/${ pkg.channel }`"
            ) {{ pkg.channel }}
            div(v-else) {{ pkg.channel }}
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
          td.border.px-4.py-2.text-center(v-if="!inAdmin")
            a(:href="_getAsset('bundle', app.app.platform, app.app.bundleID, pkg.version, pkg.id)") 下载
          td.border.px-4.py-2.text-center(v-if="inAdmin && !isIOS")
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
  props: {
    pkgs: {
      type: Array,
      required: true,
    },
  },

  data() {
    return {
      curIOSDeivceList: [],
      pkgsWithQRCode: [],
      newChannelName: "",
      changeChannelPkg: {},
    }
  },

  computed: {
    app() {
      return this.$store.state.app
    },

    isIOS() {
      return this.app.app.platform === "ios"
    },

    inAdmin() {
      return this.$route.path.startsWith("/admin")
    },

    inChannel() {
      return this.$route.params.channel != null
    },

    id() {
      return this.$route.params.id
    },
  },

  mounted() {
    this.getPkgsWithQRCode()
  },

  watch: {
    pkgs() {
      this.getPkgsWithQRCode()
    },
  },

  methods: {
    getPkgsWithQRCode() {
      Promise.all(this.pkgs.map(pkg => _util.idToQRCode(pkg.id)))
        .then(res => {
          this.pkgsWithQRCode = this.pkgs.map((item, index) => {
            return {
              qrcode: res[index],
              ...item,
            }
          })
        })
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
