import React from 'react';
import { connect } from 'react-redux';
import { Router, Route } from 'react-router';

import { bindActionCreators } from 'redux';

import * as actions from './actions';

import { Movie, MovieList } from './components';
import { App } from './containers';


class Routes extends React.Component {

  constructor(props) {
    super(props);
    this.actions = bindActionCreators(actions, this.props.dispatch);
  }

  getRandomMovie() {
    this.actions.getRandomMovie();
  }

  getMovies() {
    this.actions.getMovies();
  }

  getMovie(location) {
    this.actions.getMovie(location.params.id);
  }

  clearMovie() {
    this.actions.clearMovie();
  }

  render() {
    return (
      <Router history={this.props.history}>
        <Route component={App}>
          <Route path="/"
                 component={Movie}
                 onEnter={this.getRandomMovie.bind(this)} />

          <Route path="/all/"
                 component={MovieList}
                 onEnter={this.getMovies.bind(this)}
                 onLeave={this.clearMovie.bind(this)}/>

          <Route path="/movie/:id/"
                 component={Movie}
                 onEnter={this.getMovie.bind(this)}
                 onLeave={this.clearMovie.bind(this)} />
        </Route>
      </Router>

    );

  };
}

export default connect()(Routes);
