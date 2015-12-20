import { createStore, applyMiddleware, compose } from 'redux';
import thunkMiddleware from 'redux-thunk';
import createLogger from 'redux-logger';
import { devTools, persistState } from 'redux-devtools';

import reducer from './reducers';

const loggingMiddleware = createLogger({
  level: 'info',
  collapsed: true,
  logger: console
});

const createStoreWithMiddleware = compose(
  applyMiddleware(
    thunkMiddleware,
    loggingMiddleware
  ),
  devTools(),
  persistState(window.location.href.match(/[?&]debug_session=([^&]+)\b/))
)(createStore);

export default function configureStore(initialState) {
  return createStoreWithMiddleware(reducer, initialState);
}
