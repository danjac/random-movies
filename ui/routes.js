import React from 'react';
import { Router, Route, IndexRoute } from 'react-router';

import { App, Movie, MovieList, NotFound } from './components';

export default function (history) {
  const scrollUp = () => window.scrollTo(0, 0);

  return (
    <Router history={history}>
      <Route path="/" component={App} onEnter={scrollUp}>
        <IndexRoute component={MovieList} />
        <Route path="/movie/:id/" component={Movie} />
        <Route path="*" component={NotFound} />
     </Route>
    </Router>
  );
}
