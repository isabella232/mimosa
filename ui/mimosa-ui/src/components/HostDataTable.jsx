import _ from 'lodash';
import React from 'react';
import { Table, Checkbox, Button, Icon } from 'semantic-ui-react';
import { Link } from 'react-router-dom'
import firebase from 'firebase';

class HostDataTable extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      data: [{}],
      cap: undefined,
      hosts: [],
    }
    this.setHost = this.setHost.bind(this);
    this.setAllHost = this.setAllHost.bind(this);
    this.runTask = this.runTask.bind(this);
    this.props.firebase.auth.currentUser.getIdTokenResult().then((token) => {
      this.state.cap = token.claims.cap
    })

  }

  pullHostData = (workspace) => {
    this.props.firebase.auth.currentUser.getIdTokenResult().then((token) => {
      var stagingArray = [];
      // onSnapshot will update view if firestore updates
      this.props.firebase.app.firestore().collection("ws").doc(workspace).collection("hosts").onSnapshot((querySnapshot) => {
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

        // sort and store array in state
        // stagingArray = _.sortBy(stagingArray, ["state", "source", "public_dns"]);
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
        console.log(response.status)
        console.log(response.text())
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
    //fakeData to be used for styling, visual fixes, rather than hitting DB
    let fakeData = [
      {
        name: "onoijsaofjasmdfl;jasdofl;ask;dojasdfje",
        public_dns: "12234234590u320495u2039u4534",
        public_ip: "0.0.0.1",
        since: {
          seconds: 1234,
        },
        source: "me, myself and i",
        state: "running",
      },
      {
        name: "ksdfoijasd;fmas;odfj;ofj",
        public_dns: "123psdfosjdf4",
        public_ip: "0.0.0.1:/255",
        since: {
          seconds: 1234,
        },
        source: "vmpooler",
        state: "terminated",
      },
      {
        name: "asdc;amsd;kcnaskcn",
        public_dns: "1234",
        public_ip: "0.0.0.1",
        since: {
          seconds: 1234,
        },
        source: "bwabeabeaa",
        state: "running",
      },
      {
        name: "sdfasjo;fdjoais;djfo",
        public_dns: "1234",
        public_ip: "0.0.0.1",
        since: {
          seconds: 1234,
        },
        source: "sdfasdfasdfasdf",
        state: "terminated",
      }
    ]
    this.setState({
      // data: fakeData,
      hosts: [],
    });
    this.pullHostData(workspace);
  }
  setHost(e, data) {
    var { hosts } = this.state;
    console.log(hosts);
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
    console.log(this.state.hosts);
  }

  setAllHost(e, data) {
    var { data, hosts } = this.state;
    console.log(data);
  }

  runTask = () => {
    var { hosts } = this.state;
    console.log(hosts);
    this.props.history.push('/run-task', { response: hosts });
  }

  render() {
    var { data, cap, hosts } = this.state;
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
                <Checkbox className="all-hosts" disabled onChange={this.setAllHost} />
                Host Select
              </Table.HeaderCell>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {data && data.map((listVal) => {
              var rowState, showButton;
              if (listVal.state === 'terminated') {
                rowState = false;
                showButton = false;
              } else {
                rowState = true;
                showButton = true;
              }
              return (
                <Table.Row error={!rowState} positive={rowState}>
                  <Table.Cell>{listVal.name}</Table.Cell>
                  <Table.Cell>{listVal.hostname}</Table.Cell>
                  <Table.Cell>{listVal.ip}</Table.Cell>
                  <Table.Cell>{listVal.source}</Table.Cell>
                  <Table.Cell>{listVal.state}</Table.Cell>
                  <Table.Cell>
                    <Checkbox className="host-select" value={listVal.hostname} onChange={this.setHost} />
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