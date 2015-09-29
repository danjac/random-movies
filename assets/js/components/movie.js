import React, { PropTypes } from 'react';
import { Link } from 'react-router';

import {
  Input,
  Button,
  ButtonInput
} from 'react-bootstrap';

import { bindActionCreators } from 'redux';

import { connect } from 'react-redux';
import { pushState } from 'redux-router';

import * as actions from '../actions';

@connect(state => {
  return {
    movie: state.main.movie,
  };
})
export default class Movie extends React.Component {

  static propTypes = {
    dispatch: PropTypes.func.isRequired
  }

  constructor(props) {
    super(props);
    this.actions = bindActionCreators({ pushState, ...actions}, this.props.dispatch);
  }

  deleteMovie() {
    this.actions.deleteMovie(this.props.params.id);
    this.actions.pushState(null, "/all/");
  }

  render() {
    const movie = this.props.movie;
    if (!movie || !movie.imdbID) {
      return <div></div>;
    }
    return (
      <div className="row">
        <div className="col-md-3">
          <img className="img-responsive" src={movie.Poster} alt={movie.Title} />
        </div>
        <div className="col-md-9">
          <h2>{movie.Title}</h2>
          <dl className="dl-unstyled">
            <dt>Year</dt>
            <dd>{movie.Year}</dd>
            <dt>Actors</dt>
            <dd>{movie.Actors}</dd>
            <dt>Director</dt>
            <dd>{movie.Director}</dd>
          </dl>
          <p className="well">{movie.Plot}</p>
          <Button bsStyle="primary" onClick={this.actions.getRandomMovie.bind(this)}>Get another</Button>
          <Link className="btn btn-default" to="/all/">See all</Link>
          <Button bsStyle="danger" onClick={this.deleteMovie.bind(this)}>Delete</Button>
        </div>
      </div>
    );
  }

}


