import React, {Component} from 'react';
import { Table, Button, Icon } from 'semantic-ui-react';
import { Link } from 'react-router-dom'
import HOSTS_COLLECTION from '../utils/Fixtures/hosts_collection.js';

class HostDataTable extends Component {

  runTask = (hostname, docId) => {
    var { workspace } = this.props;
    this.props.history.push(`/ws/${workspace}/run-task`, { response: hostname, doc: docId});
  }

  render() {
    var { data } = this.props;
    return (
      <div>
        <Table className="ui single line table">
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell>Name</Table.HeaderCell>
              <Table.HeaderCell>Source</Table.HeaderCell>
              <Table.HeaderCell>State</Table.HeaderCell>
              <Table.HeaderCell></Table.HeaderCell>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {data && data.map((listVal) => {
              var rowState;
              if (listVal.available === 'false') {
                rowState = false;
              } else {
                rowState = true;
              }
              var {workspace} = this.props;
              return (
                <Table.Row error={!rowState} positive={rowState}>
                  <Table.Cell><Link to={`/ws/${workspace}/host/${listVal.id}`}>{listVal.hostname}</Link></Table.Cell>
                  <Table.Cell>{listVal.source}</Table.Cell>
                  <Table.Cell>{listVal.available}</Table.Cell>
                  <Table.Cell>
                    <Button
                      primary
                      disabled={!rowState}
                      style={{ float: "right" }}
                      onClick={() => this.runTask(listVal.hostname, listVal.id)}
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
export default HostDataTable;