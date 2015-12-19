import React from 'react';
import { Router } from 'react-router';
import createHashHistory from 'history/lib/createHashHistory';
import { syncReduxAndRouter } from 'redux-simple-router';
import { Provider } from 'react-redux';
import { DevTools, DebugPanel, LogMonitor } from 'redux-devtools/lib/react';

import getRoutes from './routes';
import configureStore from './store';

const history = createHashHistory();
const store = configureStore();

syncReduxAndRouter(history, store);

const debugPanel = window.__ENV__ === "dev" && (
  <DebugPanel top right bottom>
      <DevTools store={store} monitor={LogMonitor} />
  </DebugPanel>
) || "";

class Container extends React.Component {
  render() {
    return (
    <div>
    <Provider store={store}>
      {() => {
        return (
        <Router history={history}>
          {getRoutes(store)}
        </Router>
        );
      }}
    </Provider>
    {debugPanel}
    </div>
    );
  }
}

React.render(<Container />, document.body);
