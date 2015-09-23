import { createStore, applyMiddleware } from 'redux';
import thunkMiddleware from 'redux-thunk';
import createLogger from 'redux-logger';
import reducer from './reducer';

const loggingMiddleware = createLogger({
  level: 'info',
  collapsed: true,
  logger: console
});

const createStoreWithMiddleware = applyMiddleware(
  thunkMiddleware,
  loggingMiddleware
)(createStore);

export default function configureStore(initialState) {
  return createStoreWithMiddleware(reducer, initialState);
}
