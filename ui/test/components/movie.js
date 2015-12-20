import _ from 'lodash';
import React from 'react';
import TestUtils from 'react-addons-test-utils';
import jsdom from 'mocha-jsdom';
import { assert } from 'chai';

import { Movie } from '../../components/movie';

describe('Movie components', function() {

  jsdom({ skipWindowCheck: true });

  const movie = {
    imdbID: "ttf10000",
    Title: "Ferris Bueller's Day Off",
    Year: 1985,
    Rating: "7.3",
    Plot: "..."
  };

  it('should show a badge if seen', function() {
    const seenMovie = Object.assign({}, movie, { seen: true });
    const component = <Movie movie={seenMovie} dispatch={_.noop} />;
    const rendered = TestUtils.renderIntoDocument(component, 'div');
    const header = TestUtils.findRenderedDOMComponentWithTag(rendered, "h2");
    assert.include(header.textContent, "Seen it");
  });

  it('should show a movie title', function() {
    const component = <Movie movie={movie} dispatch={_.noop} />;
    const rendered = TestUtils.renderIntoDocument(component, 'div');
    const header = TestUtils.findRenderedDOMComponentWithTag(rendered, "h2");
    assert.include(header.textContent, "Ferris Bueller's Day Off");
  });

});
