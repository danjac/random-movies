import _ from 'lodash';

import { combineReducers } from 'redux';
import { routeReducer } from 'redux-simple-router';

import { Actions } from './constants';

export function messagesReducer(state=[], action) {
  switch(action.type) {
    case Actions.DISMISS_MESSAGE:
      return _.reject(state, msg => action.payload === msg.id);
    case Actions.ADD_MESSAGE:
      return state.concat(action.payload);
  }
  return state;
}

export function movieReducer(state=null, action) {
  switch(action.type) {
    case Actions.MOVIE_LOADED:
      return action.payload;
    case Actions.MARK_SEEN:
      return _.isNull(state) ? state : Object.assign({}, state, { seen: true });
  }
  return state;
}

export function moviesReducer(state=[], action) {
  return action.type === Actions.MOVIES_LOADED ? action.payload: state;
}

export function suggestReducer(state=null, action) {
  return action.type === Actions.NEW_SUGGESTION ? action.payload: state;
}

export default combineReducers({
  movie: movieReducer,
  movies: moviesReducer,
  suggestion: suggestReducer,
  messages: messagesReducer,
  routing: routeReducer
});
