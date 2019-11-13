import React, { Component } from 'react';
import { Container, Divider, Header, Button, Icon } from 'semantic-ui-react';
import HostDataTable from '../components/HostDataTable';
import NavMenu from '../components/NavMenu';
import { withFirebase } from '../utils/Firebase';
import { withRouter } from 'react-router-dom';

class HostsView extends Component {
  render() {
    const { authUser, firebase, history } = this.props;
    const { wsid } = this.props.match.params;
    return (
      <div>
        <NavMenu authUser={authUser} workspace={wsid} activePath="hosts"/>
        <Container>
          <Header as="h1">Host</Header>
          <Divider />
          <HostDataTable workspace={wsid} history={history} firebase={firebase}/>
        </Container>
      </div>
    )
  }
}
export default withRouter(withFirebase(HostsView));