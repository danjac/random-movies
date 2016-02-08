import _ from 'lodash';

import { List } from 'immutable';

import { combineReducers } from 'redux';
import { routeReducer } from 'react-router-redux';

import { Actions } from './constants';
import { Movie, Message } from './records';

export function messagesReducer(state = new List(), action) {
  switch (action.type) {
    case Actions.DISMISS_MESSAGE:
      return state.filterNot(msg => msg.id === action.payload);
    case Actions.ADD_MESSAGE:
      return state.unshift(new Message(action.payload).set('id', _.uniqueId()));
    default:
      return state;
  }
}

export function movieReducer(state = new Movie(), action) {
  switch (action.type) {
    case Actions.MOVIE_LOADED:
      return new Movie(action.payload);
    case Actions.MARK_SEEN:
      return state.set('seen', true);
    default:
      return state;
  }
}

export function moviesReducer(state = new List(), action) {
  return action.type === Actions.MOVIES_LOADED ?
    new List(action.payload.map(obj => new Movie(obj))) : state;
}

export function suggestReducer(state = new Movie(), action) {
  return action.type === Actions.NEW_SUGGESTION ? new Movie(action.payload) : state;
}

export default combineReducers({
  movie: movieReducer,
  movies: moviesReducer,
  suggestion: suggestReducer,
  messages: messagesReducer,
  routing: routeReducer,
});
