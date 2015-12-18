import _ from 'lodash';

import { pushPath } from 'redux-simple-router';

import ActionTypes from './actionTypes';
import * as WebAPI from './api';

export function addMessage(status, message) {
  return {
    type: ActionTypes.ADD_MESSAGE,
    message: {
      status,
      message,
      id: _.uniqueId()
    }
  };
}

export function dismissMessage(id) {
  return {
    type: ActionTypes.DISMISS_MESSAGE,
    id
  };
}

export function getMovie(id) {
  return (dispatch) => {
    WebAPI
    .getMovie(id).then(result => dispatch({
      type: ActionTypes.MOVIE_LOADED,
      movie: result.data
    }))
    .catch(() => {
      dispatch(addMessage("danger", "Sorry, no movie found"));
      dispatch(pushPath("/all/"));
    });
  }
}

export function addMovie(title) {
  return (dispatch) => {
    WebAPI
    .addMovie(title)
    .then(result => {
      dispatch(addMessage("success", "New movie added"));
      dispatch(pushPath(`/movie/${result.data.imdbID}/`));
    });
  }
}

export function deleteMovie(movie) {
  return (dispatch) => {
    dispatch(addMessage("info", "Movie deleted"));
    WebAPI.deleteMovie(movie)
  }
}

export function getMovies() {
  return (dispatch) => {
    WebAPI.getMovies()
    .then(result => {
      dispatch({
        type: ActionTypes.MOVIES_LOADED,
        movies: result.data
      });
    });
  }
}

export function getRandomMovie() {
  return (dispatch) => {
    WebAPI
    .getRandomMovie()
    .then(result => {
      dispatch({
        type: ActionTypes.MOVIE_LOADED,
        movie: result.data
      });
    });
  }
}
