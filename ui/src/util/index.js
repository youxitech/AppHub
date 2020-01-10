import Vue from "vue"

// type: bundle | icon
// platform: ios | android
// bundleId: string
// versionId: string | null 获取 icon 时可以为空
export const getAsset = (type, platform, bundleId, versionId, pkgId) => {
  switch(type) {
    case "icon":
      return `/data/${ platform }/${ bundleId }/icon.png`

    case "bundle":
      return `/data/${ platform }/${ bundleId }/${ versionId }/${ platform === "ios" ? `${ pkgId }.ipa` : `${ pkgId }.apk` }`

    default:
      throw new Error("no such type")
  }
}

export const displayError = e => {
  let text
  try {
    text = e.response.data.msg || "Some error occurred"
  } catch(error) {
    text = "Some error occurred"
  }

  Vue.notify({
    type: "error",
    text,
  })
}
