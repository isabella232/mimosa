import React, { Component } from 'react';
import { Table, Message, Grid } from 'semantic-ui-react';
import RUN_DOCUMENT from '../utils/Fixtures/run_document';
import { Link } from 'react-router-dom';

class TaskDetail extends Component {
  constructor(props) {
    super(props);
    this.state = {
      type: '',

    }
    if (this.props.firebase.auth.currentUser) {
      this.props.firebase.auth.currentUser.getIdTokenResult().then((token) => {
        this.setState({
          cap: token.claims.cap
        })
      })
    }
  }

  pullTask = (workspace, task) => {
    if (this.props.firebase.auth.currentUser) {
      this.props.firebase.auth.currentUser.getIdTokenResult().then((token) => {
        // onSnapshot will update view if firestore updates
        this.props.firebase.app.firestore().collection("ws").doc(workspace).collection("tasks").doc(task).onSnapshot((querySnapshot) => {
          // reset data to avoid duplication
          this.setState({
            data: [{}],
          });

          // real firestore data, uncomment to use
          var data = querySnapshot.data();

          // fake fixture data, comment to remove
          // var data = RUN_DOCUMENT;

          var keys = Object.keys(data.hosts),
              count = keys.length,
              hosts = data.hosts;
          var stagedHosts = [];
          keys.forEach(key => {
            hosts[key].docid = key;
            stagedHosts.push(hosts[key]);
          })
          this.setState({
            name: data.name,
            user: data.displayname,
            time: data.timestamp,
            count: count,
            hosts: stagedHosts,
          });
        });
      });
    }
  }

  componentDidMount() {
    const { workspace, task } = this.props;
    this.setState({
      hosts: [],
    });
    this.pullTask(workspace, task);
  }

  render() {
    const { name, user, time, count, hosts } = this.state;
    const { workspace } = this.props;
    
    return (
      <div>
        <Grid columns='four' divided >
          <Grid.Row>
            <Grid.Column>
              <Message
                header='Task name'
                content={name}
              />
            </Grid.Column>
            <Grid.Column>
              <Message
                header='User'
                content={user}
              />
            </Grid.Column>
            <Grid.Column>
              <Message
                header='Timestamp'
                content={time}
              />
            </Grid.Column>
            <Grid.Column>
              <Message
                header='Node count'
                content={count}
              />
            </Grid.Column>
          </Grid.Row>
        </Grid>
        <Table className="table">
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell>Node</Table.HeaderCell>
              <Table.HeaderCell>Result</Table.HeaderCell>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {hosts && hosts.map((host) => {
              console.log('hello single host ', host)
              return (
                <Table.Row>
                  <Table.Cell>
                    <Link to={`/ws/${workspace}/host/${host.docid}`}>{host.hostname}</Link>
                  </Table.Cell>
                  <Table.Cell>
                    {host.status}
                  </Table.Cell>
                </Table.Row>
              )
            })}
          </Table.Body>
        </Table>
      </div >
    )
  }
}

export default TaskDetail;