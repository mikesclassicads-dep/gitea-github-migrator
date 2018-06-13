import * as React from "react";
import * as ReactDOM from "react-dom";
import { hot } from 'react-hot-loader';

import { WelcomePage } from "./components/WelcomePage";
import { HashRouter as Router, Route, Link} from "react-router-dom";
import { GitHubButton } from "./components/GitHubButton";

ReactDOM.render(
    <Router>
        <div>
            <Route exact path="/" component={WelcomePage}/>
            <Route path="/callback" component={GitHubButton}/>
        </div>
    </Router>,
    document.getElementById("example")
);

export default hot(module)(WelcomePage)