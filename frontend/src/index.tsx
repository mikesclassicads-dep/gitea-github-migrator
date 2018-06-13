import * as React from "react";
import * as ReactDOM from "react-dom";
import { hot } from "react-hot-loader";

import { HashRouter as Router, Route } from "react-router-dom";
import { GitHubButton } from "./components/GitHubButton";
import { WelcomePage } from "./components/WelcomePage";

ReactDOM.render(
  <Router>
    <div>
      <Route exact path="/" component={WelcomePage} />
      <Route path="/callback" component={GitHubButton} />
    </div>
  </Router>,
  document.getElementById("example")
);

export default hot(module)(WelcomePage);
