import _ from 'lodash';

import { pushPath } from 'redux-simple-router';

import ActionTypes from './actionTypes';
import * as WebAPI from './api';

export function addMessage(status, message) {
  return {
    type: ActionTypes.ADD_MESSAGE,
    payload: {
      status,
      message,
      id: _.uniqueId()
    }
  };
}

export function dismissMessage(id) {
  return {
    type: ActionTypes.DISMISS_MESSAGE,
    payload: id
  };
}

export function getMovie(id) {
  return (dispatch) => {
    WebAPI
    .getMovie(id).then(result => dispatch({
      type: ActionTypes.MOVIE_LOADED,
      payload: result.data
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
    })
    .catch(() => {
      dispatch(addMessage("danger", `Sorry, couldn't find the movie "${title}"`));
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
        payload: result.data
      });
    });
  }
}

export function clearMovie() {
  return {
    type: ActionTypes.MOVIE_LOADED,
    payload: null
  };
}

export function getRandomMovie() {
  return (dispatch) => {
    WebAPI
    .getRandomMovie()
    .then(result => {
      dispatch({
        type: ActionTypes.MOVIE_LOADED,
        payload: result.data
      });
    });
  }
}

export function suggest(movie) {
  return {
    type: ActionTypes.NEW_SUGGESTION,
    payload: movie
  }
}
