var path = require('path');
var webpack = require('webpack');
var ExtractTextPlugin = require('extract-text-webpack-plugin');

module.exports = {
  context: path.join(__dirname, 'ui'),
  devtool: 'eval',
  entry: [
    'webpack-dev-server/client?http://localhost:8080',
    'webpack/hot/only-dev-server',
    './js/main.js'
  ],
  output: {
    path: path.join(__dirname, 'dist'),
    filename: "[name].js",
    publicPath: '/static/'
  },
  plugins: [
    new webpack.HotModuleReplacementPlugin(),
    new ExtractTextPlugin('[name].css', {
      allChunks: true
    })
  ],
  module: {
    loaders: [
      {
        test: /\.css$/,
        loader: ExtractTextPlugin.extract('style-loader', 'css-loader')
      },
      {
        test: /\.(png|woff|woff2|eot|ttf|svg)$/,
        loader: 'url-loader?limit=200000'
      },
      {
        test: /\.(js|jsx)$/,
        loaders: ['react-hot', 'babel?stage=0&optional[]=runtime'],
        include: path.join(__dirname, 'ui/js'),
        exclude: path.join(__dirname, 'node_modules')
      }
    ]
  },
  resolve: {
    extensions: ['', '.js'],
    modulesDirectories: ['ui', 'node_modules']
  }

};
