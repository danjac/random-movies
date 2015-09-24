import { createStore, applyMiddleware } from 'redux';
import thunkMiddleware from 'redux-thunk';
import createLogger from 'redux-logger';
import promiseMiddleware from 'redux-promise';
import reducer from './reducer';

const loggingMiddleware = createLogger({
  level: 'info',
  collapsed: true,
  logger: console
});

const createStoreWithMiddleware = applyMiddleware(
  thunkMiddleware,
  loggingMiddleware,
  promiseMiddleware
)(createStore);

export default function configureStore(initialState) {
  return createStoreWithMiddleware(reducer, initialState);
}
