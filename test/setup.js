require('babel/register')({
  stage: 0
});

var jsdom = require('jsdom');

var doc = jsdom.jsdom('<!doctype html><html><body></body></html>');
var win = doc.defaultView;

global.document = doc;
global.window = win;

propagateGlobal(win);

function propagateGlobal(window) {
  for(let key in window) {
    if (!window.hasOwnProperty(key)) continue;
    if (key in global) continue;
    global[key] = window[key];
  }
}
