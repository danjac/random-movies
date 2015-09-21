var m = require('mithril');

var components = require('./components'),
    Page = components.Page,
    MovieComponent = components.MovieComponent,
    TitlesComponent = components.TitlesComponent;

m.route.mode = "hash";

m.route(document.body, "/", {
   "/": new Page(MovieComponent),
   "/titles": new Page(TitlesComponent)
});


