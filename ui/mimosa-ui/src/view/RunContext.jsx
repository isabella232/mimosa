import React, { Component } from 'react';
import { Container, Divider, Header } from 'semantic-ui-react';
import TaskDataTable from '../components/TaskDataTable';
import NavMenu from '../components/NavMenu';
import { withFirebase } from '../utils/Firebase';
import { withRouter } from 'react-router-dom';

class RunContext extends Component {
  render() {
    const { authUser, firebase, history } = this.props;
    const { wsid } = this.props.match.params;
    return (
      <div>
        <NavMenu authUser={authUser} workspace={wsid} activePath="hosts" />
        <Container>
          <Header as="h1">Run Context</Header>
          <Divider />
          <TaskDataTable workspace={wsid} history={history} firebase={firebase} />
        </Container>
      </div>
    )
  }
}
export default withRouter(withFirebase(RunContext));