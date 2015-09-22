var Movie = require('./movie');
var MovieList = require('./list');
var Page = require('./page');

module.exports = {
  Movie: Page(Movie),
  MovieList: Page(MovieList)
};
