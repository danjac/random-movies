import { createAction } from 'redux-actions';

import * as WebAPI from './api';

export const getMovie = createAction('GET_MOVIE', id => WebAPI.getMovie(id));

export const addMovie = createAction('ADD_MOVIE', async (id, pushState) => {
  const result = await WebAPI.addMovie(id);
  if (result && result.imdbID) {
    pushState(null, `/movie/${result.imdbID}/`);
  }
  return result;
});

export const resetMovie = createAction('RESET_MOVIE');

export const deleteMovie = createAction('DELETE_MOVIE', WebAPI.deleteMovie);

export const getMovies = createAction('GET_MOVIES', WebAPI.getMovies)

export const getRandomMovie = createAction('GET_RANDOM_MOVIE', WebAPI.getRandomMovie)
