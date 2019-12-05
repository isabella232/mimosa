import React, { Component } from 'react';
import { Container, Divider, Header, Button, Icon } from 'semantic-ui-react';
import {HostDataTable} from '../components';
import {NavMenu} from '../components';
import { withFirebase } from '../utils/Firebase';
import { withRouter } from 'react-router-dom';

class HostsView extends Component {
  constructor(props) {
    super(props);
    const { wsid } = this.props.match.params;
    this.state = {
      data: [{}],
      cap: undefined,
      enableRefresh: false,
      workspace: wsid
    }
    if (this.props.firebase.auth.currentUser) {
      this.props.firebase.auth.currentUser.getIdTokenResult().then((token) => {
        this.setState({
          cap: token.claims.cap
        })
      })
    }
  }

  newDataAvailable = (workspace) => {
    this.props.firebase.app.firestore().collection("ws").doc(workspace).collection("hosts")
      .onSnapshot(() => {
          this.setState({
            enableRefresh: true,
          })
      })
  }

  pullHostData = (workspace) => {
    var stagingArray = [];
    // onSnapshot will update view if firestore updates
    this.props.firebase.app.firestore().collection("ws").doc(workspace).collection("hosts").get().then((querySnapshot) => {
      // reset data to avoid duplication
      this.setState({
        data: [{}],
      });
      // iterate through docs, add id to doc
      // add doc to array
      querySnapshot.forEach((doc) => {
        var rowData = doc.data();
        rowData["id"] = doc.id;
        stagingArray.push(rowData);
      });
      this.setState({
        data: stagingArray,
        enableRefresh: false,
      });
    });
  }

  componentDidMount() {
    const {workspace} = this.state;
    this.newDataAvailable(workspace);
    this.pullHostData(workspace);
  }

  render() {
    const { authUser, firebase, history } = this.props;
    const { workspace, data, enableRefresh } = this.state;

    return (
      <div>
        <NavMenu authUser={authUser} workspace={workspace} activePath="hosts"/>
        <Container>
          <Header as="h1">
            Host
            <Button
              color='purple'
              disabled={!enableRefresh}
              style={{ float: "right" }}
              onClick={() => this.pullHostData(workspace)}
            >
              Refresh&nbsp;
              <Icon name='refresh' />
            </Button>
          </Header>
          <Divider />
          <HostDataTable workspace={workspace} data={data} refresh={enableRefresh} history={history} firebase={firebase}/>
        </Container>
      </div>
    )
  }
}
export default withRouter(withFirebase(HostsView));