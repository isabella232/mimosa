import React, { Component } from 'react';
import { Container, Divider, Header } from 'semantic-ui-react';
import {TaskDetail} from '../components';
import {NavMenu} from '../components';
import { withFirebase } from '../utils/Firebase';
import { withRouter } from 'react-router-dom';

class RunDetail extends Component {
  render() {
    const { authUser, firebase, history } = this.props;
    const { wsid, runid } = this.props.match.params;
    return (
      <div>
        <NavMenu authUser={authUser} workspace={wsid} activePath="hosts" />
        <Container>
          <Header as="h1">Run Context</Header>
          <Divider />
          <TaskDetail workspace={wsid} task={runid} history={history} firebase={firebase} />
        </Container>
      </div>
    )
  }
}
export default withRouter(withFirebase(RunDetail));