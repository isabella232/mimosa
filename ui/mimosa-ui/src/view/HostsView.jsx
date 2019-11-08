import React from 'react';
import { Container, Divider } from 'semantic-ui-react';
import HostDataTable from '../components/HostDataTable';
import NavMenu from '../components/NavMenu';
import { withFirebase } from '../utils/Firebase';
import { withRouter } from 'react-router-dom';

class HostsView extends React.Component {
  render() {
    const { authUser, firebase } = this.props;
    return (
      <div>
        <NavMenu authUser={authUser} activePath="hosts"/>
        <Container>
          <Divider />
          <HostDataTable firebase={firebase}/>
        </Container>
      </div>
    )
  }
}
export default withRouter(withFirebase(HostsView));