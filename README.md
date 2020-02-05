# AppHub

## 基本说明

AppHub 用于内部分发 App，包上传到系统以后，移动端扫码即可下载安装。

使用一个简单的基于 token 的权限系统，输入 token 进入管理员界面。

UI 整体分为 Admin 和 Public 两大部分。Admin 用于管理操作，Public 用于用户查看各种信息。

**数据结构**

- App: bundle_id 和 platform(ios/android) 定义一个 App
- Version: App 下面会有多个版本
  - android: `[manifest.VersionName]` 定义版本
  - ios: `v[CFBundleShortVersionString].[CFBundleVersion]` 定义版本
  - 关于版本命名，可以看 [推荐的版本命名方式以及包命名方式](#推荐的版本命名方式以及包命名方式)
- Package: 同样的版本不同的包，主要用于 android 渠道包
  - Package 有一个属性叫做 channel 用于记录渠道，当上传包时，如果报名吻合 `vx.x.x.x-prod-[channel]` 结构，会自动解析渠道，否则，需要去 Admin 页面手动设置渠道。

## Public UI:

### App 首页：/[app_alias]

http://debug.haibao6688.com/ios

显示版本列表，当前最新版本的所有包

### Version 首页：/[app_alias]/[version]

http://debug.haibao6688.com/ios/v1.3.1.3

显示版本下面所有的包信息

### Package 首页：/pkg/[pkg_id]

http://debug.haibao6688.com/pkg/7643463104ba6778d999bf2622f9ab8e

显示 Package 的具体信息，PC 和 移动端都可以访问，PC 端可以下载 Package，移动端可以触发下载。

### 渠道首页：/[app_alias]/channel/[channel]

查看该渠道下的所有包，适合运营人员使用，当某个版本测试完毕需要提交上架时，负责不同运营渠道的运营人员打开对应的渠道页面即可找到需要的 APK。

## Admin UI:

### Dashboard: /admin/[app_alias]

http://debug.haibao6688.com/admin/ios

App 主要管理界面。

### Version: /admin/[app_alias]/[version]

某个版本的管理界面，可以修改渠道。

## 推荐的版本命名方式以及包命名方式

面向用户的版本始终为三位，`vx.x.x` 格式，开发时追加一位，为 4 位版本号，追踪内部迭代情况。

例如，开发 `v1.0.0` 版本，那么内部第一次提交测试的版本为 `v1.0.0.0`，第二次为 `v1.0.0.1`，直到比如说版本 `v1.0.0.22` 测试通过，上架，此时用户看到的版本号为 `v1.0.0`，但实际内部版本号为 `v1.0.0.22`。

对于 android，versionName 设置为 4 位版本号，即 `v1.0.0.22`，versionCode 基于三位版本号进行生成，生成算法有多种，例如简单的将每个版本号变成两位数拼接在一起，`v1.0.0` -> 100000。（这里要注意每个数字不能大于 99）

对于 ios，CFBundleShortVersionString 设置为前三位即 `1.0.0`，CFBundleVersion 设置为第四位即 `22`。

一个包除了包含版本信息以外，还需要环境信息（是内部测试包还是生产包）以及渠道信息（安卓特色）。

推荐的包命名方式为：`[version]-[env]-[channel].apk/ipa`，例如 `v1.0.0.0-staging-ios.ipa`。

## DEV NOTE

note:
- ngnix client_max_body_size

functions:
- preview could be protected by a password
- download all packages of a version in zip

- get device count when ad-hoc
- figure out ios package tyep: ad-hoc, appstore, etc..
- parse ipa, apk to get a verison
  - name
  - version
  - bundleid
- platform(ios/android) + bundleID define a app, app has many versions, app has a unique id
- update app info(icon, name) every time when we upload
- define a user version for display
  - android: versionCode define a version
  - ios: (versio + build) define a version
- each version has multiple packages
- each app has a download apge(always link to latest version)
- store download count
- external api for ci usage
- android support multiple channel packages

pages:
- index(login)
- app page(handle no app): display latest version
- upload
- version page(pc + mobile): list packages
- package page(pc + mobile): specific package
- 404

ios ipa type
Development Profile: Used to install an app on a registered device in debug mode
App Store Profile: A profile that is used to distribute a completed app to the App Store for sale
In-house Distribution Profile: Only available with the Enterprise developer account type, and is used for distributing apps to non-registered devices outside of the App Store. (Example: A company would use this profile type to distribute internal apps to their employees).
Ad-hoc Profile: A distribution profile for distributing an app to devices registered in the developer account

https://stackoverflow.com/questions/17584426/check-if-app-is-ad-hocdevapp-store-build-at-run-time

dev:
<key>get-task-allow</key>
<true/>

ad-hoc:
<key>ProvisionedDevices</key>
<array>
    <string>abcdef01234567890abcdef01234567890abacde</string>
    <string>1abcdef01234567890abcdef01234567890abacd</string>
    <string>2abcdef01234567890abcdef01234567890abacd</string>
</array>

in-house:
<key>ProvisionsAllDevices</key>

app-store:
no above key

we have a root data dir
- ios/android
  - [bundle_id]
    - icon.png
    - [version_id]
      - [id].ipa/apk

model
/[id]
app:
  - id: string primary key generated, editable
  - icon: path relative to root data dir
  - name: string
  - platform: 'ios' / 'android'
  - bundle id: string unique
  - install_password: string // empty means no password
  - download_count: int
  - created_at: ts
  - updated_at: ts

/version/[version]
version:
  - version: full generated string primary key unique
  - app_id: app id
  - androidVersionCode
  - androidVersionName
  - iosShortVersion
  - iosBundleVersion
  - created_at: ts
  - remark: string
  - download_count

/package/[package_id]
package:
  - id: md5 of package content
  - version_id: foreign key
  - download_count
  - name: extracted from filename, editable
  - size: int in bytes
  - created_at: ts
  - remark: string
