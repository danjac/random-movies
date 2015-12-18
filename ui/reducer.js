import _ from 'lodash';

import { combineReducers } from 'redux';
import { routeReducer } from 'redux-simple-router';

import ActionTypes from './actionTypes';

function messagesReducer(state=[], action) {
  switch(action.type) {
    case ActionTypes.DISMISS_MESSAGE:
      return _.reject(state, msg => action.id === msg.id);
    case ActionTypes.ADD_MESSAGE:
      return state.concat(action.message);
  }
  return state;
}

function movieReducer(state=null, action) {
  return action.type === ActionTypes.MOVIE_LOADED ? action.movie : state;
}

function moviesReducer(state=[], action) {
  return action.type === ActionTypes.MOVIES_LOADED ? action.movies: state;
}

export default combineReducers({
  movie: movieReducer,
  movies: moviesReducer,
  messages: messagesReducer,
  routing: routeReducer
});
