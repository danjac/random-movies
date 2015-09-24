import React from 'react';

import { Provider } from 'react-redux';
import createHashHistory from 'history/lib/createHashHistory';

import routes from './routes';
import configureStore from './store';

const store = configureStore();
const history = createHashHistory();

class Container extends React.Component {
  render() {
    return (
    <div>
    <Provider store={store}>
      {() => routes(history)}
    </Provider>
    </div>
    );
  }
}

React.render(<Container />, document.body);
