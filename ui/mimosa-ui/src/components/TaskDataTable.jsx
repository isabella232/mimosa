import React, { Component } from 'react';
import { Table } from 'semantic-ui-react';
import { Link } from 'react-router-dom'

class TaskDataTable extends Component {
  constructor(props) {
    super(props);
    this.state = {
      data: [{}],
      cap: undefined,
      hosts: [],
    }
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
      this.props.firebase.app.firestore().collection("ws").doc(workspace).collection("tasks").onSnapshot((querySnapshot) => {
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
        // console.log(this.state.data);
      });
    });
  }

  componentDidMount() {
    const { workspace } = this.props;
    this.setState({
      hosts: [],
    });
    this.pullHostData(workspace);
  }

  render() {
    var { data } = this.state;
    const { workspace } = this.props;
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
        <Table className="table">
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell>ID</Table.HeaderCell>
              <Table.HeaderCell>Task</Table.HeaderCell>
              <Table.HeaderCell>Node Count</Table.HeaderCell>
              <Table.HeaderCell>Status</Table.HeaderCell>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {data && data.map((taskData) => {
              var task, status;
              if (taskData.items) {
                taskData.items.map((item) => {
                  task = item.object;
                  status = item.status;
                })
              }
              return (
                <Table.Row>
                  <Table.Cell><Link to={`/ws/${workspace}/run/${taskData.id}`}>{taskData.id}</Link></Table.Cell>
                  <Table.Cell>{task}</Table.Cell>
                  <Table.Cell>{taskData.node_count}</Table.Cell>
                  <Table.Cell>{status}</Table.Cell>
                </Table.Row>
              )
            })}
          </Table.Body>
        </Table>
      </div>
    )
  }
}
export default TaskDataTable;