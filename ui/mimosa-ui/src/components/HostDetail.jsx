import React, { Component } from 'react';
import { Table, Message, Grid, List } from 'semantic-ui-react';
import { Link } from 'react-router-dom'
import HOST_DOCUMENT from '../utils/Fixtures/hosts_document';


class HostDetail extends Component {
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

    pullHost = (workspace, host) => {
      if (this.props.firebase.auth.currentUser) {
        this.props.firebase.auth.currentUser.getIdTokenResult().then((token) => {
          this.props.firebase.app.firestore().collection("ws").doc(workspace).collection("hosts").doc(host).onSnapshot((querySnapshot) => {
            this.setState({
              data: [{}]
            });

            // real firestore data, uncomment to use
            var data = querySnapshot.data();

            // fake fixture data, comment out to remove
            // var data = HOST_DOCUMENT;

            this.setState({
              hostname: data.hostname,
              ip: data.ip,
              name: data.name,
              status: data.state,
              source: data.source,
              time: data.timestamp,
            })
          })
        })
      }
    }

    pullTask = (workspace) => {
      if (this.props.firebase.auth.currentUser) {
        this.props.firebase.auth.currentUser.getIdTokenResult().then((token) => {
          this.props.firebase.app.firestore().collection("ws").doc(workspace).collection("tasks").onSnapshot((querySnapshot) => {
            var temp = []
            querySnapshot.forEach((doc) => {
              console.log(doc);
              var docData = doc.data();
              docData["id"] = doc.id;
              temp.push(docData);
            })
            console.log(temp);
            this.setState({
              tasks: temp
            })
          })
        })
      }
    }

    componentDidMount() {
      const { workspace, host } = this.props;
      this.pullHost(workspace, host);
      this.pullTask(workspace);
    }

    render() {
      const { hostname, ip, name, status, source, time, tasks } = this.state
      const {workspace} = this.props;
      return (
        <div>
          <Grid columns='four' divided>
            <Grid.Row>
              <Grid.Column>
                <Message
                  header='Hostname'
                  content={hostname}
                />
              </Grid.Column>
              <Grid.Column>
                <Message
                  header='Ip address'
                  content={ip}
                />
              </Grid.Column>
              <Grid.Column>
                <Message
                  header='Name'
                  content={name}
                />
              </Grid.Column>
              <Grid.Column>
                <Message
                  header='Status'
                  content={status}
                />
              </Grid.Column>
            </Grid.Row>
          </Grid>
          <List divided relaxed>
            <List.Item>
              <List.Icon name='globe' verticalAlign='middle' />
              <List.Content>
                <List.Header>Source</List.Header>
                <List.Description>{source}</List.Description>
              </List.Content>
            </List.Item>
            <List.Item>
              <List.Icon name='clock outline' verticalAlign='middle' />
              <List.Content>
                <List.Header>Timestamp</List.Header>
                <List.Description>{time}</List.Description>
              </List.Content>
            </List.Item>
          </List>
          <Table className="table" style={{tableLayout: "fixed", width: "100%"}}>
            <Table.Header>
              <Table.Row>
                <Table.HeaderCell style={{width: "20%"}}>Task</Table.HeaderCell>
                <Table.HeaderCell>Status</Table.HeaderCell>
              </Table.Row>
            </Table.Header>
            <Table.Body>
              {tasks && tasks.map((singleTask) => {
                console.log(singleTask.id);
                return (
                  <Table.Row>
                    <Table.Cell>
                      <Link to={`/ws/${workspace}/run/${singleTask.id}`}>
                        {singleTask.timestamp}
                      </Link>
                    </Table.Cell>
                    <Table.Cell style={{overflow: "hidden", textOverflow: "ellipsis"}}>{singleTask.error.Stderr}</Table.Cell>
                  </Table.Row>
                )
              })}
            </Table.Body>
          </Table>
        </div>
      )
    }
}

export default HostDetail;