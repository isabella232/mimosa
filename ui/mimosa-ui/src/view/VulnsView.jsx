import React, { Component } from 'react';
import { Container, Divider, Header, Button, Icon } from 'semantic-ui-react';
import {BasicDataTable} from '../components';
import {NavMenu} from '../components';
import { withFirebase } from '../utils/Firebase';
import { withRouter } from 'react-router-dom';
import _ from 'lodash';

class VulnsView extends Component {
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
    this.props.firebase.app.firestore().collection("ws").doc(workspace).collection("vulns")
      .onSnapshot(() => {
        this.setState({
          enableRefresh: true,
        })
      })
  }

  pullVulnData = (workspace) => {
    var stagingArray = [];
    // onSnapshot will update view if firestore updates
    this.props.firebase.app.firestore().collection("ws").doc(workspace).collection("vulns").get().then((querySnapshot) => {
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
      var sorted = _.orderBy(stagingArray, [function (obj) {
        return parseInt(obj.score, 10);
      }, 'count', 'name'], ['desc', 'desc', 'asc'])

      this.setState({
        data: sorted,
        enableRefresh: false,
      });
    });
  }

  componentDidMount() {
    const { workspace } = this.state;
    this.newDataAvailable(workspace);
    this.pullVulnData(workspace);
  }

  render() {
    const { authUser, firebase, history } = this.props;
    const { workspace, data, enableRefresh } = this.state;

    const headers = ["name", "score", "count"]; // table headers, ought to match data keys
    const linkData = [
      "name", // the name or identifier shown for the link
      "vuln", // the upper url path (i.e. host, vuln etc)
      "id" // specific variable used for detail path
    ];
    const cellKeys = ["score", "count"] // the remaining key values rendered in table
    return (
      <div>
        <NavMenu authUser={authUser} data={data} refresh={enableRefresh} workspace={workspace} activePath="hosts"/>
        <Container>
          <Header as="h1">
            Vulnerabilities
            <Button
              color="purple"
              disabled={!enableRefresh}
              style={{ float: "right"}}
              onClick={() => this.pullVulnData(workspace)}
            >
              Refresh&nbsp;
              <Icon name='refresh' />
            </Button>
          </Header>
          <Divider />
          <BasicDataTable 
            workspace={workspace}
            data={data}
            refresh={enableRefresh}
            history={history}
            firebase={firebase}
            headers={headers}
            linkData={linkData}
            cellKeys={cellKeys}
          />
        </Container>
      </div>
    )
  }
}
export default withRouter(withFirebase(VulnsView));