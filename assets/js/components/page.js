var m = require('mithril');

var Movie = require('../model').Movie;
var buttonLabel = require('./widgets').buttonLabel;

function controller(main) {

  newTitle = m.prop("");
  alerts = m.prop([]);

  function makeAlert(level, msg) {
      alerts().push({ level: level, msg: msg });
  }

  function dismissAlert(index, forceRedraw) {
    alerts().splice(index, 1);
    if (forceRedraw) {
      m.redraw(true);
    }
  }

  function addMovie(e) {
    e.preventDefault();
    var title = newTitle().trim();
    if (title) {
      newTitle("");
      Movie.addNew(title).then(function(movie) {
        makeAlert("success", "\"" + movie.Title + "\" has been added to the list!");
        m.route("/movie/" + movie.imdbID);
      });
    }
  }

  return function() {
    return {
      addMovie: addMovie,
      alerts: alerts,
      newTitle: newTitle,
      main: main,
      makeAlert: makeAlert,
      dismissAlert: dismissAlert
    };
  };

}



function view(ctrl) {

  function showAlert(alert, index) {
    // get around window timeout event issues
    var dismissAlert = function(forceRedraw) { return ctrl.dismissAlert.bind(ctrl, index, forceRedraw); };
    window.setTimeout(dismissAlert(true), 6000);

    return m("div.alert.alert-dismissable.alert-" + alert.level, {role: "alert"}, [
        m("button.close[type=button][aria-label=Close][data-dismiss=alert]",  {
          onclick: dismissAlert(false)
        },
        m("span[aria-hidden=true]", m.trust("&times;"))),
        alert.msg
    ]);
  }

  return m("div.container", [
     m("h1", "Random movies"),
      ctrl.alerts().map(showAlert),
      m("form.form-horizontal", {onsubmit: ctrl.addMovie}, [
        m("div.form-group", [
          m("input.form-control.form-control-bg", {
            type: "text",
            placeholder: "Add another title",
            value: ctrl.newTitle(),
            onchange: m.withAttr("value", ctrl.newTitle)
          }),
          m("button.form-control.btn.btn-primary[type=submit]", buttonLabel("plus", "Add new"))
        ])
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
