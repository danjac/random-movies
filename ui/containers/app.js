import React, { PropTypes } from 'react';

import 'bootstrap/dist/css/bootstrap.min.css';

import {
  Input,
  Button,
  ButtonInput,
  Glyphicon,
  Alert
} from 'react-bootstrap';

import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import * as actions from '../actions';

class App extends React.Component {

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
      node.value = "";
      this.actions.addMovie(title);
    }
  }

  renderForm() {
    return (
      <form className="form form-horizontal" onSubmit={this.addMovie.bind(this)}>
        <Input
          type="text"
          ref="title"
          placeholder="Add another title" />
        <Button
          bsStyle="primary"
          className="form-control"
          type="submit"><Glyphicon glyph="plus" /> Add</Button>
        </form>
    );
  }

  render() {
    return (
      <div className="container">
        {this.props.messages.map(msg => {
        const dismissAlert = () => { this.actions.dismissMessage(msg.id); };
        return (
        <Alert key={msg.id}
               bsStyle={msg.status}
               onDismiss={dismissAlert}
               dismissAfter={3000}>
          <p>{msg.message}</p>
        </Alert>
          );
        })}
        <h1>Random movies</h1>
        {this.renderForm()}
        {this.props.children}
      </div>
    );
  }
}

export default connect(state => {
  return {
      messages: state.messages,
      router: state.router
  };
})(App);

