import React, { Component } from 'react';
import { Container, Grid, Message, Table, Header, Icon, Button, Divider } from 'semantic-ui-react';
import { withRouter } from 'react-router';
import { withFirebase } from '../utils/Firebase';

import './workspaces.css'

class Workspaces extends Component {

  constructor(props) {
    super(props);
    this.state = {
      data: [],
    }
  }

  pullWorkspace = () => {
    this.props.firebase.app.firestore().collection("ws").onSnapshot((querySnapshot) => {
      var stagingArray = []
      querySnapshot.forEach((doc) => {
        console.log(doc.id);
        stagingArray.push(doc.id);
      });
      this.setState({
        data: stagingArray
      })
    })
  }

  componentDidMount() {
    const { authUser } = this.props;
    if (!authUser) {
      this.props.history.push('/login');
    }
    this.pullWorkspace();
  }
  handleClick = (name) => {
    console.log(name);
    this.props.history.push(`/ws/${name}/hosts`);
  }

  render() {
    const {data} = this.state;

    return (
      <Container>
        <Grid textAlign='center' style={{ height: '100vh' }}  verticalAlign='middle'>
          <Grid.Column style={{ maxWidth: 450 }}>
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
                  {data && data.map((name) => {
                    return (
                      <Table.Row>
                        <Table.Cell style={{ textAlign: "center" }} className="view-source-name">
                          {name}
                    </Table.Cell>
                        <Table.Cell>
                          <Button className="view-source-button" onClick={() => this.handleClick(name)}>
                            Enter workspace
                      </Button>
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