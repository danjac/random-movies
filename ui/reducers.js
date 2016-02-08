import _ from 'lodash';

import { List, Record, fromJS } from 'immutable';

import { combineReducers } from 'redux';
import { routeReducer } from 'redux-simple-router';

import { Actions } from './constants';

const Message = Record({
  status: '',
  message: '',
  id: 0,
});

const Movie = Record({
  Title: '',
  Actors: '',
  Poster: '',
  Year: '',
  Plot: '',
  Director: '',
  imdbID: '',
  imdbRating: '',
  seen: false,
});


export function messagesReducer(state=List(), action) {
  switch(action.type) {
    case Actions.DISMISS_MESSAGE:
      return state.filterNot(msg => msg.id === action.payload);
    case Actions.ADD_MESSAGE:
      return state.unshift(new Message(action.payload).set('id', _.uniqueId()));
  }
  return state;
}

export function movieReducer(state=new Movie(), action) {
  switch(action.type) {
    case Actions.MOVIE_LOADED:
      return new Movie(action.payload);
    case Actions.MARK_SEEN:
      return state.set('seen', true);
  }
  return state;
}

export function moviesReducer(state=List(), action) {
    return action.type === Actions.MOVIES_LOADED ?
      new List(action.payload.map(obj => new Movie(obj))) : state;
}

export function suggestReducer(state=new Movie(), action) {
  return action.type === Actions.NEW_SUGGESTION ? new Movie(action.payload): state;
}

export default combineReducers({
  movie: movieReducer,
  movies: moviesReducer,
  suggestion: suggestReducer,
  messages: messagesReducer,
  routing: routeReducer
});
