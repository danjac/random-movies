import {
  React,
  sinon,
  assert,
  expect,
  TestUtils
} from '../testHelper';

import { Movie } from '../../components/movie';

describe('Movie components', () => {
  let sandbox, component;

  const movie = {
    imdbID: "ttf10000",
    Title: "Ferris Bueller's Day Off",
    Year: 1985,
    Rating: "7.3",
    Plot: "..."
  };

  beforeEach(() => {
//sandbox = sinon.sandbox.create();
//   component = TestUtils.renderIntoDocument(<Movie movie={movie} />);
  });

  afterEach(() => {
//   sandbox.restore();
  });

  it('should show a movie title', () => {
//   const title = TestUtils.findRenderedDOMComponentWithTag(component, "h3");
//   expect(title).to.be.ok
  });

});
