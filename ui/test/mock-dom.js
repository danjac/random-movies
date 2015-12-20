import { jsdom } from 'jsdom';

const markup = '<html><body></body></html>';

export default function() {
  if (typeof document !== 'undefined') return;

  global.document = jsdom(markup || '');
  global.window = document.defaultView;

  global.navigator = {
    userAgent: 'node.js'
  };

}
