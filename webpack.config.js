var path = require('path');

module.exports = {
  entry: "./ui/js/main.js",
  output: {
    path: __dirname,
    filename: "main.js"
  },
  module: {
    loaders: [
      { test: /\.css$/, loader: "style!css" },
      {
        test: /\.(js|jsx)$/,
        loaders: ['react-hot', 'babel?stage=0&optional[]=runtime'],
        include: path.join(__dirname, 'ui/js')
      }
    ]
  }

};
