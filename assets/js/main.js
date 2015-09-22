var m = require('mithril');

var components = require('./components'),
    Movie = components.Movie,
    MovieList = components.MovieList;

m.route.mode = "hash";

m.route(document.body, "/", {
   "/": Movie,
   "/titles": MovieList,
   "/movie/:id": Movie
});


