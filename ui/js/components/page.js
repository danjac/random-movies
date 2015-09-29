var m = require('mithril');

var Movie = require('../model').Movie;
var buttonLabel = require('./widgets').buttonLabel;

function controller(main) {

  newTitle = m.prop("");
  flashMessages = m.prop([]);

  function flash(status, msg) {
      flashMessages().push({ status: status, msg: msg });
  }

  function dismissFlash(index) {
    flashMessages().splice(index, 1);
    m.redraw();
  }

  function addMovie(e) {
    e.preventDefault();
    var title = newTitle().trim();
    if (title) {
      newTitle("");
      Movie.addNew(title).then(function(movie) {
        if (movie.Title) {
          flash("success", "\"" + movie.Title + "\" has been added to the list!");
          m.route("/movie/" + movie.imdbID);
        } else {
          flash("warning", "Sorry, no movie found with the title \"" + title + "\"");
        }
      });
    }
  }

  return function() {
    return {
      addMovie: addMovie,
      flashMessages: flashMessages,
      newTitle: newTitle,
      main: main,
      flash: flash,
      dismissFlash: dismissFlash
    };
  };

}



function view(ctrl) {

  function showFlashMessage(alert, index) {
    // get around window timeout event issues: has to be better way
    var dismissFlash = ctrl.dismissFlash.bind(ctrl, index);
    window.setTimeout(dismissFlash, 6000);

    return m("div.alert.alert-dismissable.alert-" + alert.status, {role: "alert"}, [
        m("button.close[type=button][aria-label=Close][data-dismiss=alert]",  {
          onclick: dismissFlash
        },
        m("span[aria-hidden=true]", m.trust("&times;"))),
        alert.msg
    ]);
  }

  return m("div.container", [
     m("h1", "Random movies"),
      ctrl.flashMessages().map(showFlashMessage),
      m("form.form-horizontal", {onsubmit: ctrl.addMovie}, [
        m("div.form-group", m("input.form-control.form-control-bg", {
            type: "text",
            placeholder: "Add another title",
            value: ctrl.newTitle(),
            onchange: m.withAttr("value", ctrl.newTitle)
          })),
        m("div.form-group", m("button.form-control.btn.btn-primary[type=submit]", buttonLabel("plus", "Add")))
     ]),
     m.component(ctrl.main, {parent: ctrl})
  ]);
}


module.exports = function(main) {

  return {
      controller: controller(main),
      view: view
  };

};
