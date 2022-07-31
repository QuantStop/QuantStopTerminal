module.exports = {
    /*devServer: {
        https: true,
        proxy: 'https://localhost/'
    }*/
    devServer: {
        proxy: {
            '^/api': {
                target: 'https://127.0.0.1/',
                ws: true,
                changeOrigin: true
            },

        }
    }
};