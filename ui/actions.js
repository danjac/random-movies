import _ from 'lodash';

import { pushPath } from 'redux-simple-router';

import { Actions, Alert } from './constants';
import * as WebAPI from './api';

// creates a standard FSA payload object
const createAction = (type, payload) => _.merge({ type }, { payload });

export const movieLoaded = movie => createAction(Actions.MOVIE_LOADED, movie);

export const addMessage = (status, message) => createAction(Actions.ADD_MESSAGE, { status, message });

export const dismissMessage = id => createAction(Actions.DISMISS_MESSAGE, id);

export const suggest = movie => createAction(Actions.NEW_SUGGESTION, movie);

export function getMovie(id) {
  return dispatch => {
    WebAPI
    .getMovie(id)
    .then(result => dispatch(movieLoaded(result.data)))
    .catch(() => {
      dispatch(addMessage(Alert.DANGER, "Sorry, no movie found"));
      dispatch(pushPath("/"));
    });
  }
}

export function addMovie(title) {
  return dispatch => {
    WebAPI
    .addMovie(title)
    .then(result => {
      const movie = result.data;
      if (result.status === 201) {
        dispatch(addMessage(Alert.SUCCESS, `New movie "${movie.Title}" added to your collection`));
        dispatch(pushPath(`/movie/${movie.imdbID}/`));
        dispatch(movieLoaded(movie));
      } else {
        dispatch(addMessage(Alert.INFO, `"${movie.Title}" is already in your collection`));
      }
    })
    .catch(() => {
      dispatch(addMessage(Alert.WARNING, `Sorry, couldn't find the movie "${title}"`));
    });
  }
}

export function deleteMovie(movie) {
  return dispatch => {
    WebAPI.deleteMovie(movie.imdbID);
    dispatch(addMessage(Alert.INFO, `Movie "${movie.Title}" deleted`));
    dispatch(pushPath("/"));
  }
}

export function getMovies() {
  return dispatch => {
    WebAPI.getMovies()
    .then(result => {
      dispatch(createAction(Actions.MOVIES_LOADED, result.data));
    });
  }
}

export function clearMovie() {
  return movieLoaded(null);
}

export function getRandomMovie() {
  return dispatch => {
    WebAPI
    .getRandomMovie()
    .then(result => dispatch(pushPath(`/movie/${result.data.imdbID}/`)));
  }
}

export function markSeen(movie) {
  WebAPI.markSeen(movie.imdbID);
  return createAction(Actions.MARK_SEEN);
}


