module.exports = {
  devServer: {
    disableHostCheck: true,
    proxy: {
      '^/api': {
        target: 'http://localhost:3080/',
        ws: false,
        changeOrigin: true,
      },
      '^/auth': {
        target: 'http://localhost:3080/',
        ws: false,
        changeOrigin: true,
      },
      '^/static': {
        target: 'http://localhost:3080/',
        ws: false,
        changeOrigin: true,
      },
      '^/ws': {
        target: 'http://localhost:3080/ws',
        ws: true,
        changeOrigin: true,
      },
    },
  },

  transpileDependencies: ['vuetify'],

  publicPath: '/',
}
