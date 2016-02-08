import { createStore, applyMiddleware, compose } from 'redux';
import thunkMiddleware from 'redux-thunk';
import createLogger from 'redux-logger';
import { devTools, persistState } from 'redux-devtools';
import { syncHistory } from 'react-router-redux';

import reducer from './reducers';

const loggingMiddleware = createLogger({
  level: 'info',
  collapsed: true,
  logger: console,
});

export default function configureStore(history, initialState) {
  const routerMiddleware = syncHistory(history);

  const createStoreWithMiddleware = compose(
    applyMiddleware(
      routerMiddleware,
      thunkMiddleware,
      loggingMiddleware
    ),
    devTools(),
    persistState(window.location.href.match(/[?&]debug_session=([^&]+)\b/))
  )(createStore);

  const store = createStoreWithMiddleware(reducer, initialState);
  // routerMiddleware.listenForReplays(store);
  return store;
}
