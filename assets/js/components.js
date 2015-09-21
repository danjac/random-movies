var m = require('mithril');

var Movie = require('./model').Movie;

function Page(main) {

  newTitle = m.prop("");

  function addMovie(e) {
    e.preventDefault();
    var title = newTitle().trim();
    if (title) {
      newTitle("");
      Movie.addNew(title).then(function() {
        m.route("/", { title: title });
      });
    }
  }

  this.controller = function() {
    return {
      addMovie: addMovie,
      newTitle: newTitle
    };
  };

  this.view = function(ctrl) {
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
              m("button.form-control.btn.btn-primary[type=submit]", "Add")
            ])
         ]),
         m.component(main)
      ]);
  };

}

var TitlesComponent = {
  controller: function() {
    var titles = m.prop([]);
    Movie.getList().then(function(result) {
      result.sort();
      titles(result);
    });
    return {
      titles: titles
    };
  },
  view: function(ctrl) {
    return m("ul.list-unstyled", ctrl.titles().map(function(title) {
      return m("li", m("a", {href: "#/?title=" + title}, title));
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

    var title = m.route.param("title");
    if (title) {
      Movie.getMovie(title).then(movie);
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

    function showMovie() {
      if (!movie || !movie.Title || !movie.Poster) {
        return showButtons();
      }

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
    var imdbURL = "http://www.imdb.com/title/" +  movie.imdbID + "/";
    return showMovie();
  }
};

module.exports = {
  Page: Page,
  TitlesComponent: TitlesComponent,
  MovieComponent: MovieComponent
}


