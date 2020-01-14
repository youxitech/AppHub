<template lang="pug">
.h-screen(v-if="version")
  .p-10.text-3xl {{ version.version.version }}
  .flex.justify-center
    table.table-auto
      thead
        tr
          th.px-4.py-2 包名
          th.px-4.py-2 创建时间
          th.px-4.py-2 大小
          th.px-4.py-2 下载次数
          th.px-4.py-2 二维码
      tbody
        tr(
          v-for="pkg of version.packages"
          :key="pkg.id"
        )
          td.border.px-4.py-2.text-center
            a(class="text-blue-500 hover:text-blue-800" :href="`/pkg/${pkg.id}`") {{ pkg.name }}
          td.border.px-4.py-2.text-center {{ pkg.createdAt | formatTime }}
          td.border.px-4.py-2.text-center {{ pkg.size | bytesToSize }}
          td.border.px-4.py-2.text-center {{ pkg.downloadCount }}
          td.border.px-4.py-2.text-center
            img(:src="pkg.qrcode")
</template>

<script>
import QRCode from "qrcode"

export default {
  props: {
    id: {
      type: String,
      default: "",
    },

    alias: {
      type: String,
      default: "",
    },
  },

  data() {
    return {
      version: null,
    }
  },

  mounted() {
    console.log("versioin")
    this.fetchVersion()
  },

  methods: {
    fetchVersion() {
      console.log("fetch")
      axios.get(`/${ this.alias || this.$route.params.id }/${ this.id || this.$route.params.version }`)
        .then(res => {
          console.log(res.data)
          this.version = res.data
          Promise.all(this.version.packages.map(pkg => QRCode.toDataURL(location.host + `/pkg/${ pkg.id }`)))
            .then(res => {
              this.version.packages = this.version.packages.map((item, index) => {
                return {
                  qrcode: res[index],
                  ...item,
                }
              })
            })
        })
        .catch(_displayError)
    },
  },
}
</script>

<style lang="stylus">

</style>
