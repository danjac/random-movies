import React from 'react';
import { connect } from 'react-redux';
import { Router, Route, IndexRoute } from 'react-router';

import { bindActionCreators } from 'redux';

import * as actions from './actions';

import { App, Movie, MovieList, NotFound } from './components';

class Routes extends React.Component {

  constructor(props) {
    super(props);
    this.actions = bindActionCreators(actions, this.props.dispatch);
  }

  getMovie(location) {
    this.actions.getMovie(location.params.id);
  }

  render() {
    return (
      <Router history={this.props.history}>
        <Route path="/" component={App}>

          <IndexRoute component={MovieList}
                      onEnter={this.actions.getMovies.bind(this)} />

          <Route path="/movie/:id/"
                 component={Movie}
                 onEnter={this.getMovie.bind(this)}
                 onLeave={this.actions.clearMovie.bind(this)} />

          <Route path="*" component={NotFound} />
        </Route>
      </Router>

    );

  };
}
export default connect()(Routes);
