var m = require('mithril');
var _ = require('lodash');

var Movie = require('../model').Movie;
var buttonLabel = require('./widgets').buttonLabel;

module.exports = {

  controller: function(args){

    var movie = m.prop({});

    function nextMovie(e) {
      e.preventDefault();
      Movie.getRandom().then(movie);
    }

    function deleteMovie(e) {
      e.preventDefault();
      var mv = movie();
      if (window.confirm("Are you sure you want to remove \"" + mv.Title + "\"?")) {
      Movie.deleteMovie(mv.imdbID).then(function() {
          args.parent.flash("info", "\"" + mv.Title + "\" has been deleted");
          m.route("/titles/");
      });
      }
    }

    var imdbID = m.route.param("id");

    if (imdbID) {
      Movie.getMovie(imdbID).then(movie);
    } else {
      Movie.getRandom().then(movie);
    }

    return {
      movie: movie,
      nextMovie: nextMovie,
      deleteMovie: deleteMovie
    };
  },

  view: function(ctrl) {
    var movie = ctrl.movie();

    function showButtons() {
      var buttons = [
        m("button.btn.btn-primary", {onclick: ctrl.nextMovie}, buttonLabel("random", "Show me another")),
        m("a.btn.btn-default[href=#/titles]", buttonLabel("list", "See all titles")),
      ];
      if (movie.imdbID) {
        buttons.push(m("a.btn.btn-danger", {onclick: ctrl.deleteMovie}, buttonLabel("trash", "Remove this movie")));
      }

      return m("div.btn-group", buttons);
    }

    function showPoster() {
      if (movie.Poster == 'N/A') {
        return m("b", "No poster available");
      }
      return m("a[target=_blank]", {href: imdbURL},  m("img.img-responsive", {src: movie.Poster}));
    }

    if (!movie || !movie.Title || !movie.Poster) {
      return showButtons();
    }

    var imdbURL = "http://www.imdb.com/title/" +  movie.imdbID + "/";

    return m("div.row", [
      m("div.col-md-3", [
        showPoster()
      ]),
      m("div.col-md-9", [
        m("h2", m("a[target=_blank]", {href: imdbURL}, movie.Title)),
        m("dl.dl-unstyled", [
          m("dt", "Year"),
          m("dd", movie.Year),
          m("dt", "Actors"),
          m("dd", movie.Actors),
          m("dt", "Director"),
          m("dd", movie.Director),
        ]),
        m("p.well", movie.Plot),
        showButtons()
      ])
    ]);
  }
};


