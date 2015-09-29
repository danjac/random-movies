import { combineReducers } from 'redux';
import { routerStateReducer } from 'redux-router';

import ActionTypes from './actionTypes';

const initialState = {
  movie: null,
  movies: []
};

function mainReducer(state=initialState, action) {
  switch(action.type) {
    case ActionTypes.RESET_MOVIE:
      return Object.assign({}, state, { movie: null });
    case ActionTypes.GET_MOVIES:
      return Object.assign({}, state, { movies: action.payload });
    case ActionTypes.GET_MOVIE:
    case ActionTypes.GET_RANDOM_MOVIE:
      return Object.assign({}, state, { movie: action.payload });
    default:
      return state;
  }
}

export default combineReducers({
  main: mainReducer,
  router: routerStateReducer
});
