var m = require('mithril');

var components = require('./components'),
    Movie = components.Movie,
    Titles = components.Titles;

m.route.mode = "hash";

m.route(document.body, "/", {
   "/": Movie,
   "/titles": Titles
});


