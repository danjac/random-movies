var m = require('mithril');

function buttonLabel(icon, text) {
  return [
    m("i.glyphicon.glyphicon-" + icon), " " + text
  ];
}

module.exports = {
  buttonLabel: buttonLabel
};


