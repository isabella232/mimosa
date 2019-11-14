import React, {Component} from 'react';
import { Table, Checkbox, Button, Icon } from 'semantic-ui-react';
import { Link } from 'react-router-dom'
import HOSTS_COLLECTION from '../utils/Fixtures/hosts_collection.js';

class HostDataTable extends Component {
  constructor(props) {
    super(props);
    this.state = {
      data: [{}],
      cap: undefined,
      hosts: [],
    }
    this.setHost = this.setHost.bind(this);
    this.runTask = this.runTask.bind(this);
    if (this.props.firebase.auth.currentUser) {
      this.props.firebase.auth.currentUser.getIdTokenResult().then((token) => {
        this.setState({
          cap: token.claims.cap
        })
      })
    }
  }

  pullHostData = (workspace) => {
    this.props.firebase.auth.currentUser.getIdTokenResult().then((token) => {
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
        });
      });
    });
  }
  // Call cloud function, since we don't expect result we don't do anything
  callCloudFunction = (functionName, hostid) => {

    this.props.firebase.auth.currentUser.getIdToken().then(function (idToken) {
      // FIXME - ACCESS TOKEN SHOULD BE ADDED AS A BEARER TOKEN
      fetch('https://mimosa-esp-tfmdd2vwoq-uc.a.run.app/' + functionName + "?access_token=" + idToken, {
        method: 'POST',
        mode: 'cors',
        cache: 'no-cache',
        headers: {
          'Content-Type': 'application/json'
        },
        redirect: 'follow',
        referrer: 'no-referrer',
        body: JSON.stringify({ "workspace": "ws1", "id": hostid })
      }).then(response => {
        // console.log(response.status)
        // console.log(response.text())
      })
        .catch(error => {
          console.error('Error during Mimosa:', error);
        });

    }).catch(function (error) {
      console.error('Error during Mimosa:', error);
    });

  }

  componentDidMount() {
    const { workspace } = this.props;
    this.setState({
      data: HOSTS_COLLECTION, //comment out when not using fixture data
      hosts: [],
    });
    // pull the read data from firestore
    // this.pullHostData(workspace);
  }
  setHost(e, data) {
    var { hosts } = this.state;
    if (data.checked) {
      if (hosts && !hosts.includes(data.value)) {
        hosts.push(data.value);
      }
    } else {
      if (hosts && hosts.includes(data.value)) {
        let index = hosts.indexOf(data.value)
        hosts.splice(index, 1);
      }
    }
    this.setState({
      hosts: hosts,
    })
  }

  runTask = () => {
    var { hosts } = this.state;
    var { workspace } = this.props;
    this.props.history.push(`/ws/${workspace}/run-task`, { response: hosts });
  }

  render() {
    var { data } = this.state;
    /**
     * Iterate through firestore data and render table
     * the document ID is used in Task Output button
     * to pass it to the Task view
     *
     * Also Run Task and Task Output buttons will not
     * render unless host is running (should add other checks in future)
     */
    return (
      <div>
        <Button secondary>
          Refresh&nbsp;
              <Icon name='refresh' />
        </Button>
        <Button primary onClick={this.runTask}>
          Run Task&nbsp;
              <Icon name='bolt' />
        </Button>
        <Table className="ui single line table">
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell>Name</Table.HeaderCell>
              <Table.HeaderCell>Hostname</Table.HeaderCell>
              <Table.HeaderCell>IP</Table.HeaderCell>
              <Table.HeaderCell>Source</Table.HeaderCell>
              <Table.HeaderCell>State</Table.HeaderCell>
              <Table.HeaderCell>
                Host Select
              </Table.HeaderCell>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {data && data.map((listVal) => {
              var rowState;
              if (listVal.state === 'terminated') {
                rowState = false;
              } else {
                rowState = true;
              }
              var {workspace} = this.props;
              return (
                <Table.Row error={!rowState} positive={rowState}>
                  <Table.Cell><Link to={`/ws/${workspace}/host/${listVal.id}`}>{listVal.name}</Link></Table.Cell>
                  <Table.Cell>{listVal.hostname}</Table.Cell>
                  <Table.Cell>{listVal.ip}</Table.Cell>
                  <Table.Cell>{listVal.source}</Table.Cell>
                  <Table.Cell>{listVal.state}</Table.Cell>
                  <Table.Cell>
                    <Checkbox className="host-select" value={listVal.name} onChange={this.setHost} />
                  </Table.Cell>
                </Table.Row>
              )
            })}
          </Table.Body>
        </Table>
      </div>
    )
  }
}
export default HostDataTable;