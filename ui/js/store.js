import { createStore, applyMiddleware, compose } from 'redux';
import { reduxReactRouter } from 'redux-router';
import createHashHistory from 'history/lib/createHashHistory';
import thunkMiddleware from 'redux-thunk';
import createLogger from 'redux-logger';
import promiseMiddleware from 'redux-promise';
import { devTools, persistState } from 'redux-devtools';

import reducer from './reducer';

const loggingMiddleware = createLogger({
  level: 'info',
  collapsed: true,
  logger: console
});

const createAppStore = compose(
  applyMiddleware(
    thunkMiddleware,
    loggingMiddleware,
    promiseMiddleware
  ),
  reduxReactRouter({
    createHistory: createHashHistory
  }),
  devTools(),
  persistState(window.location.href.match(/[?&]debug_session=([^&]+)\b/))
)(createStore);

export default function configureStore(initialState) {
  return createAppStore(reducer, initialState);
}
