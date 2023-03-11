const cdn = {
  css: [
    '//lf3-cdn-tos.bytecdntp.com/cdn/expire-1-M/bulma/0.9.3/css/bulma.min.css',
    '//lf3-cdn-tos.bytecdntp.com/cdn/expire-1-M/font-awesome/4.7.0/css/font-awesome.min.css',
  ],
  js: [
    '//lf6-cdn-tos.bytecdntp.com/cdn/expire-1-M/vue/3.2.13/vue.global.min.js',
    '//cdn.staticfile.org/vue-router/4.1.6/vue-router.global.min.js'
  ]
}
 
module.exports = {
  publicPath: '/',
  productionSourceMap: process.env.NODE_ENV !== 'production',
  devServer: {
    port: 8077
  },
  chainWebpack: config => {
    config.resolve.symlinks(true)
  },
  configureWebpack: {
    externals: {
      'vue': 'Vue',
      'returnCitySN': 'returnCitySN'
    }
  },
  pages: {
    index: {
      entry: 'src/main.js',
      template: 'public/index.html',
      filename: 'index.html',
      title: '后台管理',
      chunks: ['chunk-vendors', 'chunk-common', 'index'],
      cdn: cdn
    }
  }
}