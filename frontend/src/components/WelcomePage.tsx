import * as React from "react";
import { Grid, Header, Image, Segment } from 'semantic-ui-react'
import { GitHubButton } from "./GitHubButton";


export class WelcomePage extends React.Component {
    render() {
        return  <div className='login-form'>
            {/*
      Heads up! The styles below are necessary for the correct render of this example.
      You can do same with CSS, the main idea is that all the elements up to the `Grid`
      below must have a height of 100%.
    */}
            <style>{`
      body > div,
      body > div > div,
      body > div > div > div.login-form {
        height: 100%;
      }
    `}</style>
            <Grid textAlign='center' style={{ height: '100%' }} verticalAlign='middle'>
                <Grid.Column style={{ maxWidth: 450 }}>
                    <div>
                        <Header as='h2' color='green' textAlign='center'>
                            <Image src='assets/img/gitea.svg' /> Gitea Migrator
                        </Header>
                        <Segment attached>
                            Gitea Migrator helps you to move your GitHub repositories to your Gitea instance. To proceed, please sign in
                            via GitHub.
                        </Segment>
                        <Segment attached>
                            <GitHubButton/>
                        </Segment>
                    </div>
                </Grid.Column>
            </Grid>
        </div>;
    }
}