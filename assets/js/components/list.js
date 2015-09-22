var m = require('mithril');
var _ = require('lodash');

var Movie = require('../model').Movie;

function getInitial(title) {
    if (title.match(/^the\s/i)) {
      title = title.substring(4);
    }
    var upCase = title.charAt(0).toUpperCase();
    if (upCase.toLowerCase() !== upCase) { // ASCII letter
      return upCase;
    }
    return '-';
}

module.exports = {
  controller: function() {
    var movies = m.prop([]);
    Movie.getList().then(function(result) {
      result.sort(function(left, right) {
        return left.Title > right.Title ? 1 : (left.Title < right.Title ? -1 : 0);
      });
      movies(result);
    });
    return {
      movies: movies
    };
  },
  view: function(ctrl) {

    var groups = _.groupBy(ctrl.movies(), function(movie) {
      return getInitial(movie.Title);
    });
    return m("div.row", _.chunk(_.sortBy(Object.keys(groups)), 4).map(function(chunk) {
      return m("div.col-md-3", chunk.map(function(initial) {
        return [
          m("h3", initial),
          m("ul.list-unstyled", groups[initial].map(function(movie) {
              return m("li", m("a", {href: "#/movie/" + movie.imdbID}, movie.Title));
          }))
        ];
      }));
    }));
  }

};


