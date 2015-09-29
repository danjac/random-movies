import React from 'react';
import { Route } from 'react-router';
import { ReduxRouter } from 'redux-router';
import { Provider } from 'react-redux';
import { DevTools, DebugPanel, LogMonitor } from 'redux-devtools/lib/react';

import { App, Movie, MovieList } from './components';

import * as actions from './actions';
import configureStore from './store';

const store = configureStore();

function getRandomMovie() {
  store.dispatch(actions.getRandomMovie());
}

function getMovies() {
  store.dispatch(actions.getMovies());
}

function getMovie(location) {
  store.dispatch(actions.getMovie(location.params.id));
}

function resetMovie() {
  store.dispatch(actions.resetMovie());
}

class Container extends React.Component {
  render() {
    return (
    <div>
    <Provider store={store}>
      {() => {
      return (
        <ReduxRouter>
          <Route component={App}>
            <Route path="/" component={Movie} onEnter={getRandomMovie} onLeave={resetMovie} />
            <Route path="/all/" component={MovieList} onEnter={getMovies} />
            <Route path="/movie/:id/" component={Movie} onEnter={getMovie} onLeave={resetMovie} />
          </Route>
        </ReduxRouter>
        );
      }}
    </Provider>
    <DebugPanel top right bottom>
        <DevTools store={store} monitor={LogMonitor} />
    </DebugPanel>
    </div>
    );
  }
}

React.render(<Container />, document.body);
