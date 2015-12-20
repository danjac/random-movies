import React from 'react';
import { Well } from 'react-bootstrap';
import { connect } from 'react-redux';

export default function NotFound(props) {
    return (
        <div className="container">
            <Well>Sorry, you are in the wrong place.</Well>
        </div>
    );
}
