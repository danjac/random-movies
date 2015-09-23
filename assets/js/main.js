import React from 'react';
import { Provider } from 'react-redux';

import configureStore from './store';

const store = configureStore();


class Page extends React.Component {
  render() {
      return <div>OK</div>;
  }
}


class Container extends React.Component {
  render() {
    return (
    <div>
    <Provider store={store}>
      <Page />
    </Provider>
    </div>
    );
  }
}

React.render(<Container />, document.body);
