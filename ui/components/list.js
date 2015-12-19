import React, { PropTypes } from 'react';
import _ from 'lodash';
import { Link } from 'react-router';
import { bindActionCreators } from 'redux';

import { connect } from 'react-redux';

import {
  Glyphicon,
  Grid,
  Row,
  Col
} from 'react-bootstrap';


import * as actions from '../actions';


function getInitial(title) {
  if (title.match(/^the\s/i)) {
    title = title.substring(4);
  }
  var upCase = title.charAt(0).toUpperCase();
  if (upCase.toLowerCase() !== upCase) { // ASCII letter
    return upCase;
  }
  return '-';
}

const ListItem = props => {
    const { movie } = props;
    return (
      <li>
        <Link to={`/movie/${movie.imdbID}/`}>{movie.Title}</Link> {movie.seen? <Glyphicon glyph="ok" /> : ''}
      </li>
    );
};

const InitialGroup = props => {
  const { initial, group } = props;
  return (
    <div>
      <h3>{initial}</h3>
      <ul className="list-unstyled">
        {group.map(movie => {
          return <ListItem key={movie.imdbID} movie={movie} />;
        })}
      </ul>
    </div>
  );
};

function normalizeTitle(title) {
    const lower = title.toLowerCase();
    ["the", "a", "an"].forEach(a => {
      if (lower.startsWith(a + " ")) {
        return lower.substring(a.length + 1);
      }
    });
    return lower;
}


class MovieList extends React.Component {

  static propTypes = {
    dispatch: PropTypes.func.isRequired
  }

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(actions, dispatch);
  }

  render() {

    const movies = _.sortBy(this.props.movies, movie => normalizeTitle(movie.Title));
    const groups = _.groupBy(movies, movie => getInitial(movie.Title));
    const cols = _.chunk(_.sortBy(Object.keys(groups)), 4);
    const rows = _.chunk(cols, 4);

    return (
      <div>
        {movies.length ? <h3>Total {movies.length} movies</h3> : ''}
        <Grid>
        {rows.map((row, index) => {
          return (
            <Row key={index}>
              {row.map((col, index) => {
                return (
                  <Col key={index} md={3}>
                    {col.map(initial => <InitialGroup key={initial} group={groups[initial]} initial={initial} />)}
                  </Col>
                );
              })}
            </Row>
          );
        })}
        </Grid>
      </div>
    );
  }

}

export default connect(state => {
  return {
      movies: state.movies
  };
})(MovieList);

