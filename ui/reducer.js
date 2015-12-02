import { combineReducers } from 'redux';
import { routeReducer } from 'redux-simple-router';

import ActionTypes from './actionTypes';

const initialState = {
  movie: null,
  movies: [],
  messages: []
};

function addMessage(messages, msg) {
  messages = messages.splice();
  messages.push(msg);
  return messages;
}

function removeMessage(messages, index) {
  messages = messages.splice();
  messages.splice(index, 1);
  return messages;
}


function mainReducer(state=initialState, action) {
  switch(action.type) {

    case ActionTypes.DISMISS_MESSAGE:
      return Object.assign({}, state, { messages: removeMessage(state.messages, action.payload) });

    case ActionTypes.RESET_MOVIE:
      return Object.assign({}, state, { movie: null });

    case ActionTypes.GET_MOVIES:
      return Object.assign({}, state, { movies: action.payload.data });

    case ActionTypes.GET_MOVIE:
    case ActionTypes.GET_RANDOM_MOVIE:
      return Object.assign({}, state, { movie: action.payload.data });

    case ActionTypes.DELETE_MOVIE:
      return Object.assign({}, state, { messages: addMessage(state.messages, { status: "info", msg: "Your movie has been deleted"}) });

    case ActionTypes.ADD_MOVIE:
      let msg;
      if (action.error) {
        msg = { status: "warning", msg: "We couldn't find that movie" };
      } else {
        msg = { status: "success", msg: "Your movie has been added" };
      }
      return Object.assign({}, state, { messages: addMessage(state.messages, msg) });

    default:
      return state;
  }
}

export default combineReducers({
  main: mainReducer,
  routing: routeReducer
});
