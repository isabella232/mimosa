import React, {Component} from 'react';
import {withRouter} from 'react-router';
import { Card, Container, Grid, Header, Icon } from 'semantic-ui-react';
import './notfound.css';
import cookie from 'react-cookies';

class NotFound extends Component {

  render() {
    if (!cookie.load('userEmail')) {
      this.props.history.push('/login');
    } else {
      this.props.history.push('/ws');
    }
    return (
      <Container>
        <Grid textAlign='center' style={{ height: '100vh' }} verticalAlign='middle'>
          <Grid.Column style={{ maxWidth: 450 }}>
            <Header as='h2' color='teal' textAlign='center'>
              <Icon name="frown outline" />We have a problem
            </Header>
            <Card centered className='notfound'>
              <Card.Content>
                <Card.Header>Sorry, page not found</Card.Header>
                <Card.Meta>404 not found</Card.Meta>
                <Card.Description>
                  The page you are trying to view does not exist. Please check your workspace, host id or run id in the url.
            </Card.Description>
              </Card.Content>
            </Card>
          </Grid.Column>
        </Grid>
      </Container>
    )
  }
}

export default withRouter(NotFound);