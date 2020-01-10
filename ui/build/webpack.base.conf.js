const path = require("path")
const MiniCssExtractPlugin = require("mini-css-extract-plugin")
const webpack = require("webpack")
const { VueLoaderPlugin } = require("vue-loader")

const resolve = dir => path.join(__dirname, "..", dir)
const apiEnv = process.env.API_ENV || "dev"

const postcssLoader = {
  loader: "postcss-loader",
  options: {
    ident: "postcss",
    plugins: loader => [
      require("tailwindcss")(),
      require("autoprefixer")(),
    ],
  },
}

/* eslint-disable */
const styleLoaders = process.env.NODE_ENV === "production" ?
  [
    MiniCssExtractPlugin.loader,
    "css-loader",
    postcssLoader,
  ]
  :
  [
    "vue-style-loader",
    "css-loader",
    postcssLoader,
  ]
/* eslint-enable */

module.exports = {
  entry: {
    app: "./src/main.js",
  },
  output: {
    path: resolve("../static/ui"),
    publicPath: "/",
  },
  resolve: {
    extensions: [".js", ".vue"],
    alias: {
      "@": resolve("src/pages"),
      "util": resolve("src/util"),
    },
    modules: [
      resolve("src"),
      resolve("src/components"),
      "node_modules",
    ],
  },
  module: {
    rules: [
      {
        test: /\.vue$/,
        loader: "vue-loader",
      },
      {
        test: /\.styl(us)?$/,
        use: styleLoaders.concat([
          {
            loader: "stylus-loader",
            options: {
              import: [
                path.join(__dirname, "..", "src", "config.styl"),
              ],
            },
          },
        ]),
      },
      {
        test: /\.css$/,
        use: styleLoaders,
      },
      {
        test: /\.js$/,
        use: [
          "babel-loader",
        ],
        exclude: /node_modules/,
      },
      {
        test: /\.pug$/,
        loader: "pug-plain-loader",
      },
      {
        test: /\.(png|jpe?g|gif|svg)(\?.*)?$/,
        loader: "url-loader",
        options: {
          limit: 10000,
          name: "imgs/[name].[hash:7].[ext]",
        },
      },
      {
        test: /\.(woff2?|eot|ttf|otf)(\?.*)?$/,
        loader: "url-loader",
        options: {
          limit: 10000,
          name: "fonts/[name].[hash:7].[ext]",
        },
      },
    ],
  },
  plugins: [
    new VueLoaderPlugin(),
    new webpack.DefinePlugin({
      "__API_ENV__": JSON.stringify(apiEnv),
    }),
  ],
}
