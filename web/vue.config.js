const {defineConfig} = require('@vue/cli-service')
let webpack = require('webpack')
module.exports = defineConfig({
    transpileDependencies: true,
    lintOnSave: false,   //关闭eslint校验
    configureWebpack: {
        // disable performance hints
        performance: {
            hints: false
        },
        plugins: [
            new webpack.optimize.MinChunkSizePlugin({minChunkSize: 10000})
        ]
    },

    publicPath: '/',
    outputDir: 'dist',
    crossorigin: "anonymous",
    devServer: {
        allowedHosts: ['127.0.0.1:5678'],
        port: 8888,
    }
})
