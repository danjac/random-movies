import React from 'react';
import { Route, Router } from 'react-router';
import createHashHistory from 'history/lib/createHashHistory';
import { syncReduxAndRouter } from 'redux-simple-router';
import { Provider } from 'react-redux';
import { DevTools, DebugPanel, LogMonitor } from 'redux-devtools/lib/react';

import { Movie, MovieList } from './components';
import { App } from './containers';

import * as actions from './actions';
import configureStore from './store';

const history = createHashHistory();
const store = configureStore();

syncReduxAndRouter(history, store);

function getRandomMovie() {
  store.dispatch(actions.getRandomMovie());
}

function getMovies() {
  store.dispatch(actions.getMovies());
}

function getMovie(location) {
  store.dispatch(actions.getMovie(location.params.id));
}

const debugPanel = window.__ENV__ === "dev!" && (
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
          <Route component={App} onEnter={() => console.log("onenter app")}>
            <Route path="/" component={Movie} onEnter={getRandomMovie} />
            <Route path="/all/" component={MovieList} onEnter={getMovies} />
            <Route path="/movie/:id/" component={Movie} onEnter={getMovie} />
          </Route>
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
