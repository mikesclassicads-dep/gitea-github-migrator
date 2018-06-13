import * as React from "react";
import { Button, Icon } from "semantic-ui-react";
import { client } from "../configuration";

export interface IGitHubButtonState {
    loading: boolean
}

export class GitHubButton extends React.Component<{}, IGitHubButtonState> {

    constructor(props: {}) {
        super(props);
        this.state = {
            loading: false,
        };
    }

    render() {
        return <Button icon={true} loading={this.state.loading} disabled={this.state.loading} labelPosition='left' onClick={this.submitToGitHub}>
            <Icon name='github'/>
            Login with GitHub
        </Button>
    }

    submitToGitHub() {
        this.setState({loading: true})
        client.getGitHubRedirect()
        .then(res => {
            this.setState({loading: false})
            window.location.href = res.url
        })
    }
}