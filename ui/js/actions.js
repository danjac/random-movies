import { createAction } from 'redux-actions';

import ActionTypes from './actionTypes';
import * as WebAPI from './api';

export const getMovie = createAction(ActionTypes.GET_MOVIE, id => WebAPI.getMovie(id));

export const addMovie = createAction(ActionTypes.ADD_MOVIE, async (title, pushState) => {
  const result = await WebAPI.addMovie(title);
  if (result && result.data.imdbID) {
    pushState(null, `/movie/${result.data.imdbID}/`);
  }
  return result;
});

export const dismissMessage = createAction(ActionTypes.DISMISS_MESSAGE);

export const resetMovie = createAction(ActionTypes.RESET_MOVIE);

export const deleteMovie = createAction(ActionTypes.DELETE_MOVIE, WebAPI.deleteMovie);

export const getMovies = createAction(ActionTypes.GET_MOVIES, WebAPI.getMovies)

export const getRandomMovie = createAction(ActionTypes.GET_RANDOM_MOVIE, WebAPI.getRandomMovie)
