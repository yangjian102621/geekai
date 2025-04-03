const { defineConfig } = require("@vue/cli-service");
const path = require("path");
let webpack = require("webpack");
module.exports = defineConfig({
  transpileDependencies: true,
  lintOnSave: false, //关闭eslint校验
  productionSourceMap: false, //在生产模式中禁用 Source Map，既可以减少包大小，也可以加密源码

  configureWebpack: {
    // disable performance hints
    performance: {
      hints: false,
    },
    plugins: [new webpack.optimize.MinChunkSizePlugin({ minChunkSize: 10000 })],
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "src"),
      },
    },
  },

  publicPath: "/",

  outputDir: "dist",
  crossorigin: "anonymous",
  devServer: {
    client: {
      overlay: false // 关闭错误覆盖层
    },
    allowedHosts: "all",
    port: 8888,
    proxy: {
      "/static/upload/": {
        target: process.env.VUE_APP_API_HOST,
        changeOrigin: true,
      },
    },
  },
});
