import { createAction } from 'redux-actions';

import ActionTypes from './actionTypes';
import * as WebAPI from './api';

export const getMovie = createAction(ActionTypes.GET_MOVIE, id => WebAPI.getMovie(id));

export const addMovie = createAction(ActionTypes.ADD_MOVIE, async (id, pushState) => {
  const result = await WebAPI.addMovie(id);
  if (result && result.imdbID) {
    console.log("new result", result);
    pushState(null, `/movie/${result.imdbID}/`);
  }
  return result;
});

export const dismissMessage = createAction(ActionTypes.DISMISS_MESSAGE, index => index);

export const resetMovie = createAction(ActionTypes.RESET_MOVIE);

export const deleteMovie = createAction(ActionTypes.DELETE_MOVIE, WebAPI.deleteMovie);

export const getMovies = createAction(ActionTypes.GET_MOVIES, WebAPI.getMovies)

export const getRandomMovie = createAction(ActionTypes.GET_RANDOM_MOVIE, WebAPI.getRandomMovie)
