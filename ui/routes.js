import React from 'react';
import { Route } from 'react-router';
import * as actions from './actions';

import { Movie, MovieList } from './components';
import { App } from './containers';


export default function(store) {

  function getRandomMovie() {
    store.dispatch(actions.getRandomMovie());
  }

  function getMovies() {
    store.dispatch(actions.getMovies());
  }

  function getMovie(location) {
    store.dispatch(actions.getMovie(location.params.id));
  }

  function clearMovie() {
    store.dispatch(actions.clearMovie());
  }


  return (
      <Route component={App}>
        <Route path="/" component={Movie} onEnter={getRandomMovie} />
        <Route path="/all/" component={MovieList} onEnter={getMovies} onLeave={clearMovie}/>
        <Route path="/movie/:id/" component={Movie} onEnter={getMovie} onLeave={clearMovie} />
      </Route>
  );


}
