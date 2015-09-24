import React from 'react';
import { Router, Route } from 'react-router';

import { App, Movie, MovieList } from './components';

export default function(history) {
  return (
    <Router history={history}>
      <Route component={App}>
        <Route path="/" component={Movie} />
        <Route path="/all/" component={MovieList} />
        <Route path="/movie/:id/" component={Movie} />
      </Route>
    </Router>
  );
}

