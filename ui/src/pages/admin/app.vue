<template lang="pug">
.rounded.overflow-hidden.shadow-lg.p-6.flex.flex-col.bg-white.h-full(v-if="app")
  .flex
    img.w-24.h-24.mr-2.rounded(:src="_getAsset('icon', app.app.platform, app.app.bundleID)")
    .ml-5
      router-link.font-semibold.text-blue-500(
        class="hover:underline"
        :to="`/${ app.app.alias }`"
        target="_blank"
      ) {{ app.app.name }}
      .text-sm.text-gray-600.mt-3 Platform: {{ app.app.platform }}
      .text-sm.text-gray-600 BundleID: {{ app.app.bundleID }}
      .text-sm.text-gray-600 Download: {{ app.app.downloadCount }}
    .flex.ml-auto.items-center
      .mr-5 APP alias: {{ app.app.alias }}
      button(
        class="mr-8 hover:text-teal-500"
        @click="$modal.show('changeAppId'), newAppId = app.app.alias"
      ) 修改

  table.table-auto.mt-10
    thead
      tr
        th.px-4.py-2 版本
        th.px-4.py-2 Package数目
        th.px-4.py-2 更新时间
        th.px-4.py-2 设置
    tbody
      tr(
        v-for="version, index of app.versions"
        :key="version.id"
      )
        td.border.px-4.py-2.text-center
          .flex.items-center.justify-center
            .rounded-full.h-3.w-3.mr-3(:class="index === 0 ? 'bg-teal-500' : 'bg-white'")
            router-link(
              class="text-blue-500 hover:text-blue-800"
              :to="`/admin/${ app.app.alias }/${ version.version }`"
              ) {{ version.version }}
        td.border.px-4.py-2.text-center {{ version.pacakgeCount }}
        td.border.px-4.py-2.text-center {{ version.updatedAt | formatTime }}
        td.border.px-4.py-2.text-center
          button.text-teal-500(
            class="hover:underline"
            @click="() => setDefaultVersion(version.id)"
          ) 设为默认版本
          button.text-red-500.ml-5(
            class="hover:underline"
            @click="() => onDeleteVersion(version.id)"
          ) 删除

  _button.mt-auto(@click="onDeleteApp") 删除此 APP

  modal(name="changeAppId" @opened="afterModalOpen")
    .flex.items-center.justify-center.w-full.h-full.flex-col.p-6
      _input(
        v-model="newAppId"
        ref="aliasInput"
        @enter="changeAppId"
      )
      .flex.mt-10
        button.mr-10(class="mr-8 hover:text-teal-500" @click="changeAppId") Confirm
        button(class="mr-8 hover:text-teal-500" @click="$modal.hide('changeAppId')") Cancel
</template>

<script>
export default {
  data() {
    return {
      newAppId: "",
      app: null,
    }
  },

  mounted() {
    this.fetchApp()
  },

  watch: {
    $route() {
      this.fetchApp()
    },
  },

  methods: {
    fetchApp() {
      return axios.get(`/admin/apps/${ this.$route.params.id }`)
        .then(res => {
          this.app = res.data
        })
        .catch(_displayError)
    },

    setDefaultVersion(id) {
      axios.post(`/admin/versions/${ id }/active`)
        .then(res => {
          const idx = this.app.versions.findIndex(i => i.id === id)
          this.app.versions.unshift(this.app.versions[idx])
          this.app.versions.splice(idx + 1, 1)
          _showSuccess("Success!")
        })
        .catch(_displayError)
    },

    changeAppId() {
      if(this.newAppId === "") {
        _showErr("请输入alias")
        return
      }

      axios.patch(`/admin/apps/${ this.app.app.id }`, {
        alias: this.newAppId,
      })
        .then(() => {
          this.$router.push("/admin/" + this.newAppId)
          this.$modal.hide("changeAppId")
          _showSuccess("修改成功")
        })
        .catch(_displayError)
    },

    onDeleteVersion(id) {
      _showConfirm("确认删除？")
        .then(() => {
          return axios.delete(`/admin/versions/${ id }`)
            .then(() => {
              this.app.versions.splice(this.app.versions.findIndex(i => i.id === id), 1)
            })
            .catch(_displayError)
        })
        .catch(() => {})
    },

    onDeleteApp() {
      _showConfirm("确认删除？")
        .then(() => {
          return axios.delete(`/admin/apps/${ this.app.app.id }`)
            .then(() => {
              this.$router.push("/admin")
            })
            .catch(_displayError)
        })
        .catch(() => {})
    },

    afterModalOpen() {
      this.$nextTick(() => {
        this.$refs.aliasInput.focus()
      })
    },
  },
}
</script>

<style lang="stylus">

</style>
