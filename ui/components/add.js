import React, { PropTypes } from 'react';

import {
  Input,
  Button,
  Glyphicon,
} from 'react-bootstrap';


export default class AddMovieForm extends React.Component {

  constructor(props) {
    super(props);
    this.addMovie = this.addMovie.bind(this);
  }

  addMovie(event) {
    event.preventDefault();
    const node = this.refs.title.getInputDOMNode();
    const title = node.value.trim();

    if (title) {
      node.value = '';
      this.props.addMovie(title);
    }
  }

  render() {
    return (
      <div className="container">
        <form className="form form-horizontal" onSubmit={this.addMovie}>
        <Input
          type="text"
          ref="title"
          placeholder="Add another title"
        />
        <Button
          bsStyle="primary"
          className="form-control"
          type="submit"
        ><Glyphicon glyph="plus" /> Add</Button>
        </form>
      </div>
    );
  }
}

AddMovieForm.propTypes = {
  addMovie: PropTypes.func.isRequired,
};
