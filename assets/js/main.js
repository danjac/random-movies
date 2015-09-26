import React from 'react';
import { Route } from 'react-router';

import { ReduxRouter } from 'redux-router';

import { Provider } from 'react-redux';

import { App, Movie, MovieList } from './components';

import store from './store';

class Container extends React.Component {
  render() {
    return (
    <div>
    <Provider store={store}>
      {() => {
      return (
        <ReduxRouter>
          <Route component={App}>
            <Route path="/" component={Movie} />
            <Route path="/all/" component={MovieList} />
            <Route path="/movie/:id/" component={Movie} />
          </Route>
        </ReduxRouter>
        );
      }}
    </Provider>
    </div>
    );
  }
}

React.render(<Container />, document.body);
