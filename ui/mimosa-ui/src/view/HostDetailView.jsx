import React, {Component} from 'react';
import { Container, Divider, Header } from 'semantic-ui-react';
import {HostDetail} from '../components';
import {NavMenu} from '../components';
import { withFirebase } from '../utils/Firebase';
import { withRouter } from 'react-router-dom';

class HostDetailView extends Component {
  render() {
    const {authUser, firebase, history} = this.props;
    const { wsid, hostid } = this.props.match.params;
    console.log(authUser)
    return (
      <div>
        <NavMenu authUser={authUser} workspace={wsid} activePath="hosts" />
        <Container>
          <Header as="h1">Host Details</Header>
          <Divider />
          <HostDetail workspace={wsid} host={hostid} history={history} firebase={firebase} />
        </Container>
      </div>
    )
  }
}

export default withRouter(withFirebase(HostDetailView));