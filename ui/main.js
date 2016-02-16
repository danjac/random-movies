import 'babel-polyfill';
import React from 'react';
import ReactDOM from 'react-dom';
import { browserHistory } from 'react-router';
import { Provider } from 'react-redux';

import configureRoutes from './routes';
import configureStore from './store';
import { suggest } from './actions';

const store = configureStore(browserHistory);
const routes = configureRoutes(browserHistory, store);

new WebSocket(`ws://${window.location.host}/api/suggest`).onmessage = event => {
  store.dispatch(suggest(JSON.parse(event.data)));
};

const Container = () => {
  return (
    <Provider store={store}>
      {routes}
    </Provider>
  );
};

ReactDOM.render(<Container />, document.getElementById('app'));
