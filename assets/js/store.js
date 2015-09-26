import { createStore, applyMiddleware, compose } from 'redux';
import { reduxReactRouter } from 'redux-router';
import createHashHistory from 'history/lib/createHashHistory';
import thunkMiddleware from 'redux-thunk';
import createLogger from 'redux-logger';
import promiseMiddleware from 'redux-promise';

import reducer from './reducer';

const loggingMiddleware = createLogger({
  level: 'info',
  collapsed: true,
  logger: console
});

export default compose(
  applyMiddleware(
    thunkMiddleware,
    loggingMiddleware,
    promiseMiddleware
  ),
  reduxReactRouter({
    createHistory: createHashHistory
  })
)(createStore)(reducer);
