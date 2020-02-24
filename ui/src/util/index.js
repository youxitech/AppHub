import Vue from "vue"
import QRCode from "qrcode"

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

export const showSuccess = text => {
  Vue.notify({
    type: "success",
    text,
  })
}

export const showErr = text => {
  Vue.notify({
    type: "error",
    text,
  })
}

export const displayError = e => {
  console.error(e)
  let text
  try {
    text = typeof e === "string" ? e : (e.response.data.msg || "Some error occurred")
  } catch(error) {
    text = "Some error occurred"
  }

  showErr(text)
}

export const idToQRCode = id => QRCode.toDataURL(location.protocol + location.host + `/pkg/${ id }`)
