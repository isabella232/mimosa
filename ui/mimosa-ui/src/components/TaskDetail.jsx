import _ from 'lodash';
import React, { Component } from 'react';
import { Table, Message, Grid, Icon } from 'semantic-ui-react';
import { Link } from 'react-router-dom'
import firebase from 'firebase';

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
    this.props.firebase.auth.currentUser.getIdTokenResult().then((token) => {
      var stagingArray = [];
      // onSnapshot will update view if firestore updates
      this.props.firebase.app.firestore().collection("ws").doc(workspace).collection("tasks").doc(task).onSnapshot((querySnapshot) => {
        // reset data to avoid duplication
        this.setState({
          data: [{}],
        });
        // iterate through docs, add id to doc
        // add doc to array
        // querySnapshot.forEach((doc) => {
        //   console.log(doc);
        // });
        console.log(querySnapshot.data());
        var data = querySnapshot.data();
        var type, name, status, count, nodes, result;
        if (data && data.items) {
          data.items.map((item) => {
            type = item.action;
            name = item.object;
            status = item.status;
            nodes = item.node;
            result = item.result
            count = data.node_count;
          })
        }

        this.setState({
          type: type,
          name: name,
          status: status,
          count: count,
          nodes: nodes,
          result: JSON.stringify(result)
        });

        // sort and store array in state
      });
    });
  }

  componentDidMount() {
    const { workspace, task } = this.props;
    this.setState({
      // data: fakeData,
      hosts: [],
    });
    this.pullTask(workspace, task);
  }

  render() {
    const { data, type, name, status, count, nodes, result} = this.state
    
    console.log(data);
    return (
      <div>
        <Grid columns='four' divided >
          <Grid.Row>
            <Grid.Column>
              <Message
                header='Run type'
                content={type}
              />
            </Grid.Column>
            <Grid.Column>
              <Message
                header='Task name'
                content={name}
              />
            </Grid.Column>
            <Grid.Column>
              <Message
                header='Run status'
                content={status}
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
            <Table.Row>
              <Table.Cell>{nodes}</Table.Cell>
              <Table.Cell>{result}</Table.Cell>
            </Table.Row>
          </Table.Body>
        </Table>
      </div >
    )
  }
}
export default TaskDetail;