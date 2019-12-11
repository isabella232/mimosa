import React, { Component } from 'react';
import { Container, Grid, Message, Table, Header, Icon, Button } from 'semantic-ui-react';
import { withRouter } from 'react-router';
import { withFirebase } from '../utils/Firebase';
import cookie from 'react-cookies';

import './workspaces.css'

class Workspaces extends Component {

  constructor(props) {
    super(props);
    this.state = {
      data: [],
    }
  }

  pullWorkspace = (uid) => {
    this.props.firebase.app.firestore().collection("users").doc(uid).onSnapshot((querySnapshot) => {
      var workspace = querySnapshot.data().workspaces;
      this.setState({
        data: workspace
      })
    })
  }

  componentDidMount() {
    const { authUser } = !!cookie.load('userEmail');
    if (!cookie.load('userEmail')) {
      this.props.history.push('/login');
    }
    var cook = cookie.load('userEmail');

    this.pullWorkspace(cook.id);
  }
  handleClick = (name) => {
    this.props.history.push(`/ws/${name}/hosts`);
  }

  render() {
    const {data} = this.state;
    return (
      <Container>
        <Grid textAlign='center' style={{ height: '100vh' }}  verticalAlign='middle'>
          <Grid.Column style={{ maxWidth: 600 }}>
            <Message>
              <Header as='h2' color='teal' textAlign='center'>
                <Icon name="cube" />Select from your workspaces
              </Header>
              <Table>
                <Table.Header>
                  <Table.Row>
                    <Table.HeaderCell style={{ textAlign: "center" }}>Source Name</Table.HeaderCell>
                    <Table.HeaderCell></Table.HeaderCell>
                  </Table.Row>
                </Table.Header>
                <Table.Body>
                  {data && Object.entries(data).map(([key, value]) => {
                    return (
                      <Table.Row>
                        <Table.Cell style={{ textAlign: "center" }} className="view-source-name">
                          <Button className="view-source-button" onClick={() => this.handleClick(key)}>
                            {key}
                          </Button>
                        </Table.Cell>
                        <Table.Cell>
                          {value}
                        </Table.Cell>
                      </Table.Row>
                    )
                  })}
                </Table.Body>
              </Table>
            </Message>
          </Grid.Column>
        </Grid>
      </Container>
    )
  }
}

export default withRouter(withFirebase(Workspaces));