import React, { PropTypes } from 'react';
import { pushState } from 'redux-router';

import 'bootstrap/dist/css/bootstrap.min.css';

import {
  Input,
  Button,
  ButtonInput,
  Glyphicon
} from 'react-bootstrap';

import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import * as actions from '../actions';

@connect(state => {
  return {
      router: state.router
  };
})
export default class App extends React.Component {

  static propTypes = {
    dispatch: PropTypes.func.isRequired,
    children: PropTypes.node
  }

  constructor(props, context) {
    super(props, context);
    const { dispatch } = this.props;
      this.actions = bindActionCreators({
        pushState,
        ...actions
      }, dispatch);
  }

  addMovie(event) {
    event.preventDefault();
    const node = this.refs.title.getInputDOMNode(),
          title = node.value.trim();

    if (title) {
      node.value = "";
      this.actions.addMovie(title, this.actions.pushState);
    }
  }

  renderForm() {
    return (
      <form className="form form-horizontal" onSubmit={this.addMovie.bind(this)}>
        <Input type="text"
                ref="title"
                placeholder="Add another title" />
              <Button bsStyle="primary" type="submit"><Glyphicon glyph="plus" /> Add</Button>
        </form>
    );
  }

  render() {
    return (
      <div className="container">
        <h1>Random movies</h1>
        {this.renderForm()}
        {this.props.children}
      </div>
    );
  }
}


