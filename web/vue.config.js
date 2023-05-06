const {defineConfig} = require('@vue/cli-service')
let webpack = require('webpack')
module.exports = defineConfig({
    transpileDependencies: true,
    configureWebpack: {
        // disable performance hints
        performance: {
            hints: false
        },
        plugins: [
            new webpack.optimize.MinChunkSizePlugin({minChunkSize: 10000})
        ]
    },

    publicPath: '/chat',
    outputDir: '../src/dist',
    crossorigin: "anonymous",
    devServer: {
        allowedHosts: ['127.0.0.1:5678'],
        port: 8888,
    }
})
