<template lang="pug">
.rounded.overflow-hidden.shadow-lg.p-6.flex.flex-col.bg-white
  .flex
    img.w-24.h-24.mr-2.rounded(:src="_getAsset('icon', app.platform, app.bundleID)")
    .ml-5
      .font-semibold {{ app.name }}
      .text-sm.text-gray-600.mt-3 Platform: {{ app.platform }}
      .text-sm.text-gray-600 BundleID: {{ app.bundleID }}
      .text-sm.text-gray-600 Download: {{ app.downloadCount }}
    .flex.ml-auto.items-center APP alias: {{ app.alias }}
</template>

<script>
export default {
  props: {
    app: {
      type: Object,
      required: true,
    },
  },
}
</script>

<style lang="stylus">
</style>
