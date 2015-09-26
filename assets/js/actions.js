import { createAction } from 'redux-actions';

import * as WebAPI from './api';

export const getMovie = createAction('GET_MOVIE', id => WebAPI.getMovie(id));

export const addMovie = createAction('ADD_MOVIE', async (id, router) => {
  const result = await WebAPI.addMovie(id);
  if (result && result.imdbID) {
    console.log(result, router);
    router.pushState(null, `/movie/${result.imdbID}/`);
  }
  return result;
});

export const resetMovie = createAction('RESET_MOVIE');

export const getMovies = createAction('GET_MOVIES', WebAPI.getMovies)

export const getRandomMovie = createAction('GET_RANDOM_MOVIE', WebAPI.getRandomMovie)
