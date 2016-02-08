import React from 'react';
import ReactDOM from 'react-dom';
import { hashHistory } from 'react-router';
import { Provider } from 'react-redux';
import { DevTools, DebugPanel, LogMonitor } from 'redux-devtools/lib/react';

import Routes from './routes';
import configureStore from './store';
import { suggest } from './actions';

const store = configureStore(hashHistory);

const debugPanel = window.__ENV__ === "dev!!!" && (
  <DebugPanel top right bottom>
      <DevTools store={store} monitor={LogMonitor} />
  </DebugPanel>
) || "";

new WebSocket(`ws://${window.location.host}/api/suggest`).onmessage = event => {
  store.dispatch(suggest(JSON.parse(event.data)));
};

const Container = props => {
  return (
    <Provider store={store}>
      <Routes history={hashHistory} />
    </Provider>
  );
}

ReactDOM.render(<Container />, document.getElementById("app"));
