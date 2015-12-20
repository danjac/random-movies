import _ from 'lodash';

import { pushPath } from 'redux-simple-router';

import { Actions, Alert } from './constants';
import * as WebAPI from './api';

function movieLoaded(movie) {
  return {
      type: Actions.MOVIE_LOADED,
      payload: movie
  };
}

export function addMessage(status, message) {
  return {
    type: Actions.ADD_MESSAGE,
    payload: {
      status,
      message,
      id: _.uniqueId()
    }
  };
}

export function dismissMessage(id) {
  return {
    type: Actions.DISMISS_MESSAGE,
    payload: id
  };
}

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
      dispatch({
        type: Actions.MOVIES_LOADED,
        payload: result.data
      });
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
  return {
    type: Actions.MARK_SEEN
  };
}

export function suggest(movie) {
  return {
    type: Actions.NEW_SUGGESTION,
    payload: movie
  }
}
