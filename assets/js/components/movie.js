import React, { PropTypes } from 'react';
import { Link } from 'react-router';

import {
  Input,
  Button,
  ButtonInput
} from 'react-bootstrap';

import { bindActionCreators } from 'redux';

import { connect } from 'react-redux';

import * as actions from '../actions';

@connect(state => {
  return {
      movie: state.movie
  };
})
export default class Movie extends React.Component {

  static propTypes = {
    dispatch: PropTypes.func.isRequired
  }

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(actions, dispatch);
  }

  fetchMovie(props) {
    const id = props.params && props.params.id;
    if (id && (!this.props.movie || id !== this.props.movie.imdbID)) {
      this.actions.getMovie(id);
    } else {
      this.actions.getRandomMovie();
    }
  }

  componentDidMount() {
    this.fetchMovie(this.props);
  }

  componentWillReceiveProps(props) {
    if (props.params.id !== this.props.params.id) {
      this.fetchMovie(props);
    }
  }

  componentWillUnmount() {
    this.actions.resetMovie();
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
          <Button bgStyle="primary" onClick={this.actions.getRandomMovie.bind(this)}>Get another</Button>
          <Link className="btn btn-default" to="/all/">See all</Link>
        </div>
      </div>
    );
  }

}


