import React, { PropTypes } from 'react';
import { connect } from 'react-redux';
import { Router, Route, IndexRoute } from 'react-router';
import { bindActionCreators } from 'redux';

import * as actions from './actions';
import { App, Movie, MovieList, NotFound } from './components';

class Routes extends React.Component {

  constructor(props) {
    super(props);
    this.actions = bindActionCreators(actions, this.props.dispatch);
    this.getMovies = this.getMovies.bind(this);
    this.getMovie = this.getMovie.bind(this);
    this.clearMovie = this.clearMovie.bind(this);
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
        <Route path="/" component={App}>
          <IndexRoute
            component={MovieList}
            onEnter={this.actions.getMovies}
          />

        <Route
          path="/movie/:id/"
          component={Movie}
          onEnter={this.getMovie}
          onLeave={this.clearMovie}
        />

        <Route path="*" component={NotFound} />
       </Route>
      </Router>

    );
  }
}

Routes.propTypes = {
  history: PropTypes.object.isRequired,
  dispatch: PropTypes.func.isRequired,
};

export default connect()(Routes);
