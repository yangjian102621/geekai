const {defineConfig} = require('@vue/cli-service')
let webpack = require('webpack')
module.exports = defineConfig({
    transpileDependencies: true,
    lintOnSave: false,   //关闭eslint校验
    productionSourceMap: false, //在生产模式中禁用 Source Map，既可以减少包大小，也可以加密源码
    configureWebpack: {
        // disable performance hints
        performance: {
            hints: false
        },
        plugins: [
            new webpack.optimize.MinChunkSizePlugin({minChunkSize: 10000})
        ]
    },

    publicPath: process.env.NODE_ENV === 'production' ? '/' : '/',

    outputDir: 'dist',
    crossorigin: "anonymous",
    devServer: {
        allowedHosts: ['127.0.0.1:5678'],
        port: 8888,
    }
})
