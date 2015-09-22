var m = require('mithril');
var _ = require('lodash');

var Movie = require('./model').Movie;

function pageController(main) {

  newTitle = m.prop("");

  function addMovie(e) {
    e.preventDefault();
    var title = newTitle().trim();
    if (title) {
      newTitle("");
      Movie.addNew(title).then(function(movie) {
        m.route("/movie/" + movie.imdbID);
      });
    }
  }

  return function() {
    return {
      addMovie: addMovie,
      newTitle: newTitle,
      main: main
    };
  };

}


function pageView(ctrl) {
  return m("div.container", [
     m("h1", "Random movies"),
     m("form.form-horizontal", {onsubmit: ctrl.addMovie}, [
        m("div.form-group", [
          m("input.form-control.form-control-bg", {
            type: "text",
            placeholder: "Add another title",
            value: ctrl.newTitle(),
            onchange: m.withAttr("value", ctrl.newTitle)
          }),
          m("button.form-control.btn.btn-primary[type=submit]", "Add new")
        ])
     ]),
     m.component(ctrl.main)
  ]);
}


function Page(main) {

  return {
      controller: pageController(main),
      view: pageView
  };

}

function getInitial(title) {
    if (title.match(/$the\s/i)) {
      title = title.substring(4);
    }
    var upCase = title.charAt(0).toUpperCase();
    if (upCase.toLowerCase() !== upCase) { // ASCII letter
      return upCase;
    }
    return '-';
}

var ListComponent = {
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
    var initials = Object.keys(groups);
    initials.sort();

    return m("div.row", _.chunk(initials, 4).map(function(chunk) {
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

var MovieComponent = {

  controller: function(){

    var movie = m.prop({});

    function nextMovie(e) {
      e.preventDefault();
      Movie.getRandom().then(movie);
    }

    var imdbID = m.route.param("id");

    if (imdbID) {
      Movie.getMovie(imdbID).then(movie);
    } else {
      Movie.getRandom().then(movie);
    }

    return {
      'movie': movie,
      'nextMovie': nextMovie
    };
  },

  view: function(ctrl){
    var movie = ctrl.movie();

    function showButtons() {
      return m("div.btn-group", [
        m("button.btn.btn-primary", {onclick: ctrl.nextMovie}, "Show me another"),
        m("a.btn.btn-default[href=#/titles]", "See all titles")
      ]);
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

module.exports = {
  MovieList: Page(ListComponent),
  Movie: Page(MovieComponent)
};


