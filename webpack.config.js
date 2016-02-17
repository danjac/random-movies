var path = require('path');
var webpack = require('webpack');
var ExtractTextPlugin = require('extract-text-webpack-plugin');
var UglifyJsPlugin = webpack.optimize.UglifyJsPlugin;

require('es6-promise').polyfill();

var env = process.env.WEBPACK_ENV || 'dev';

var entry = ['babel-polyfill', './main.js'];

var plugins = [
  new ExtractTextPlugin('[name].css', {
    allChunks: true
  })
];

var jsloaders = ['babel-loader?presets[]=react,presets[]=es2015'];

switch(env) {
  case 'dev':
  entry.unshift('webpack-dev-server/client?http://localhost:8080');
  entry.unshift('webpack/hot/only-dev-server');
  plugins.unshift(new webpack.HotModuleReplacementPlugin());
  jsloaders.unshift('react-hot');
  break;
  case 'production':
  plugins.push(new DefinePlugin({ 'process.env': { NODE_ENV: 'production '} }));
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
    publicPath: 'http://localhost:8080/static/'
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
        test: /\.js$/,
        loaders: jsloaders,
        include: path.join(__dirname, 'ui'),
        exclude: path.join(__dirname, 'node_modules')
      }
    ]
  },
  resolve: {
    root: path.join(__dirname),
    extensions: ['', '.js'],
    modulesDirectories: ['./ui', './node_modules']
  }

};
