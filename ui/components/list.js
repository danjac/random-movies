import React, { PropTypes } from 'react';
import _ from 'lodash';
import { Link } from 'react-router';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { List } from 'immutable';

import {
  Grid,
  Row,
  Col,
} from 'react-bootstrap';


import * as actions from '../actions';
import { Movie } from '../records';

function stripArticle(title) {
  return title.match(/^the\s/i) ? title.substring(4) : title;
}

function getInitial(title) {
  var upCase = stripArticle(title).charAt(0).toUpperCase();
  if (upCase.toLowerCase() !== upCase) { // ASCII letter
    return upCase;
  }
  return '-';
}

const ListItem = props => {
  const { movie } = props;
  const link = <Link to={`/movie/${movie.imdbID}/`}>{movie.Title}</Link>;

  return (
    <li>
      {movie.seen ? <s>{link}</s> : <span>{link}</span>}
    </li>
  );
};

ListItem.propTypes = {
  movie: PropTypes.instanceOf(Movie).isRequired,
};


const InitialGroup = props => {
  const { initial, group } = props;
  const movies = _.sortBy(group, movie => stripArticle(movie.Title.toLowerCase()));
  return (
    <div>
      <h3>{initial}</h3>
      <ul className="list-unstyled">
        {movies.map(movie => {
          return <ListItem key={movie.imdbID} movie={movie} />;
        })}
      </ul>
    </div>
  );
};

InitialGroup.propTypes = {
  initial: PropTypes.string.isRequired,
  group: PropTypes.array.isRequired,
};

class MovieList extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(actions, dispatch);
  }

  render() {
    const { movies } = this.props;
    const groups = movies.groupBy(movie => getInitial(movie.Title)).toJS();
    const cols = _.chunk(_.sortBy(Object.keys(groups)), 4);
    const rows = _.chunk(cols, 4);

    return (
      <div>
        {movies.size ? <h3>Total {movies.size} movies</h3> : ''}
        <Grid>
        {rows.map((row, rowIndex) => {
          return (
            <Row key={rowIndex}>
              {row.map((col, colIndex) => {
                return (
                  <Col key={colIndex} md={3}>
                    {col.map(initial => {
                      return (
                        <InitialGroup
                          key={initial}
                          group={groups[initial]}
                          initial={initial}
                        />);
                    })}
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

MovieList.propTypes = {
  dispatch: PropTypes.func.isRequired,
  movies: PropTypes.instanceOf(List),
};

export default connect(state => {
  return {
    movies: state.movies,
  };
})(MovieList);
