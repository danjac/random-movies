var path = require('path');
var webpack = require('webpack');
var ExtractTextPlugin = require('extract-text-webpack-plugin');
var UglifyJsPlugin = webpack.optimize.UglifyJsPlugin;

var env = process.env.WEBPACK_ENV;

var jsLoaders = ['babel?stage=0&optional[]=runtime'];

var entry = ['./js/main.js'];

var plugins = [
  new ExtractTextPlugin('[name].css', {
    allChunks: true
  })
];

switch(process.env.WEBPACK_ENV) {
  case 'dev':
  jsLoaders.unshift('react-hot');
  entry.unshift('webpack-dev-server/client?http://localhost:8080');
  entry.unshift('webpack/hot/only-dev-server');
  plugins.unshift(new webpack.HotModuleReplacementPlugin());
  break;
  case 'prod':
  plugins.push(new UglifyJsPlugin({ minimize: true }));
  break;
}

module.exports = {
  context: path.join(__dirname, 'ui'),
  devtool: 'source-map',
  entry: entry,
  output: {
    path: path.join(__dirname, 'dist'),
    filename: "[name].js",
    publicPath: '/static/'
  },
  plugins: plugins,
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
        loaders: jsLoaders,
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
