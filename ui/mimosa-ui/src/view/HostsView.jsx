import React from 'react';
import { Container, Divider, Header, Button, Icon } from 'semantic-ui-react';
import HostDataTable from '../components/HostDataTable';
import NavMenu from '../components/NavMenu';
import { withFirebase } from '../utils/Firebase';
import { withRouter } from 'react-router-dom';

class HostsView extends React.Component {
  render() {
    const { authUser, firebase } = this.props;
    const { wsid } = this.props.match.params;
    return (
      <div>
        <NavMenu authUser={authUser} activePath="hosts"/>
        <Container>
          <Header as="h1">Host</Header>
          <Divider />
          <HostDataTable workspace={wsid} firebase={firebase}/>
        </Container>
      </div>
    )
  }
}
export default withRouter(withFirebase(HostsView));