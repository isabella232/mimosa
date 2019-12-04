import React, { Component } from 'react';
import { Table, Message, Grid, Icon, Button } from 'semantic-ui-react';
import { Link } from 'react-router-dom'
import HOST_DOCUMENT from '../utils/Fixtures/hosts_document';


class VulnDetail extends Component {
  constructor(props) {
    super(props);
    this.state = {
      cap: '',
      hostname: '-',
      ip: '-',
      name: '-',
      status: '-',
    };
    if (this.props.firebase.auth.currentUser) {
      this.props.firebase.auth.currentUser.getIdTokenResult().then((token) => {
        this.setState({
          cap: token.claims.cap
        })
      })
    }
  }
  runTask = (hostname, docId) => {
    console.log(hostname, docId);
    var { workspace } = this.props;
    this.props.history.push(`/ws/${workspace}/run-task`, { response: hostname, doc: docId});
  }


  pullHost = (workspace, vuln) => {
        console.log("this one over here", vuln);
        this.props.firebase.app.firestore().collection("ws").doc(workspace).collection("vulns").doc(vuln).onSnapshot((querySnapshot) => {
          this.setState({
            data: [{}]
          });

          // real firestore data, uncomment to use
          var data = querySnapshot.data();

          // fake fixture data, comment out to remove
          // var data = HOST_DOCUMENT;
          console.log("hello", data);

          this.setState({
            name: data.name,
            count: data.count,
            score: data.score,
            hosts: data.hosts,
          })
        })
  }

  componentDidMount() {
    const { workspace, vuln } = this.props;
    this.pullHost(workspace, vuln);
  }

  render() {
    const { name, count, score, hosts } = this.state
    const { workspace } = this.props;
    return (
      <div>
        <Grid columns='three'>
          <Grid.Row>
            <Grid.Column>
              <Message
                header='Name'
                content={name}
              />
            </Grid.Column>
            <Grid.Column>
              <Message
                header='Count'
                content={count}
              />
            </Grid.Column>
            <Grid.Column>
              <Message
                header='score'
                content={score}
              />
            </Grid.Column>
          </Grid.Row>
        </Grid>
        {/* <List divided relaxed> */}
          {/* <List.Item>
            <List.Icon name='globe' verticalAlign='middle' />
            <List.Content>
              <List.Header>Source</List.Header>
              <List.Description>{source}</List.Description>
            </List.Content>
          </List.Item> */}
          {/* <List.Item>
            <List.Icon name='clock outline' verticalAlign='middle' />
            <List.Content>
              <List.Header>Timestamp</List.Header>
              <List.Description>{time}</List.Description>
            </List.Content>
          </List.Item> */}
        {/* </List> */}
        <Table className="table" style={{ tableLayout: "fixed", width: "100%" }}>
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell>Hostname</Table.HeaderCell>
              <Table.HeaderCell>Name</Table.HeaderCell>
              <Table.HeaderCell></Table.HeaderCell>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {hosts && Object.keys(hosts).map((key) => {
              return (
                <Table.Row>
                  <Table.Cell>
                    <Link to={`/ws/${workspace}/host/${key}`}>
                      {hosts[key].hostname}
                    </Link>
                  </Table.Cell>
                  <Table.Cell>
                    {hosts[key].name}
                  </Table.Cell>
                  <Table.Cell>
                    <Button 
                      primary 
                      style={{ float: "right" }}
                      onClick={() => this.runTask(hosts[key].hostname, key)}
                    >
                      Run Task&nbsp;
                      <Icon name='bolt' />
                    </Button>
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

export default VulnDetail;