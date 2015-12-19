import _ from 'lodash';

import { combineReducers } from 'redux';
import { routeReducer } from 'redux-simple-router';

import ActionTypes from './actionTypes';

function messagesReducer(state=[], action) {
  switch(action.type) {
    case ActionTypes.DISMISS_MESSAGE:
      return _.reject(state, msg => action.payload === msg.id);
    case ActionTypes.ADD_MESSAGE:
      return state.concat(action.payload);
  }
  return state;
}

function movieReducer(state=null, action) {
  return action.type === ActionTypes.MOVIE_LOADED ? action.payload : state;
}

function moviesReducer(state=[], action) {
  return action.type === ActionTypes.MOVIES_LOADED ? action.payload: state;
}

function suggestReducer(state=null, action) {
  return action.type === ActionTypes.NEW_SUGGESTION ? action.payload: state;
}

export default combineReducers({
  movie: movieReducer,
  movies: moviesReducer,
  suggestion: suggestReducer,
  messages: messagesReducer,
  routing: routeReducer
});
