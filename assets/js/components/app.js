import React, { PropTypes } from 'react';

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
    this.actions = bindActionCreators(actions, dispatch);
  }

  addMovie(event) {
    event.preventDefault();
    const node = this.refs.title.getInputDOMNode(),
          title = node.value.trim();

    if (title) {
      this.actions.addMovie(title, this.props.router);
      node.value = "";
    }
  }

  renderForm() {
    return (
      <form className="form form-horizontal" onSubmit={this.addMovie.bind(this)}>
        <Input type="text"
                ref="title"
                placeholder="Add another title" />
        <ButtonInput bsStyle="primary"
                     type="submit">Add</ButtonInput>
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


