import React, {Component} from 'react';
import { Container, Divider, Header } from 'semantic-ui-react';
import {VulnDetail} from '../components';
import {NavMenu} from '../components';
import { withFirebase } from '../utils/Firebase';
import { withRouter } from 'react-router-dom';

class VulnDetailView extends Component {
  render() {
    const {authUser, firebase, history} = this.props;
    const { wsid, vulnid } = this.props.match.params;
    return (
      <div>
        <NavMenu authUser={authUser} workspace={wsid} activePath="hosts" />
        <Container>
          <Header as="h1">Vulnerability Details</Header>
          <Divider />
          <VulnDetail workspace={wsid} vuln={vulnid} history={history} firebase={firebase} />
        </Container>
      </div>
    )
  }
}

export default withRouter(withFirebase(VulnDetailView));