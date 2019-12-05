import React, {Component} from 'react';
import { Table, Checkbox, Button, Icon } from 'semantic-ui-react';
import { Link } from 'react-router-dom'
import HOSTS_COLLECTION from '../utils/Fixtures/hosts_collection.js';

class BasicDataTable extends Component {

  runTask = (hostname, docId) => {
    var { workspace } = this.props;
    this.props.history.push(`/ws/${workspace}/run-task`, { response: hostname, doc: docId});
  }

  render() {
    var { data, linkData, headers, cellKeys } = this.props;
    return (
      <div>
        <Table className="ui single line table">
          <Table.Header>
            <Table.Row>
              {headers && headers.map((element) => {
                var columnHeader = element.charAt(0).toUpperCase() + element.slice(1);
                return (
                  <Table.HeaderCell>{columnHeader}</Table.HeaderCell>
                )
              })}
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {data && data.map((listVal) => {
              var {workspace} = this.props;
              return (
                <Table.Row>
                  <Table.Cell><Link to={`/ws/${workspace}/${linkData[1]}/${listVal[linkData[2]]}`}>{listVal[linkData[0]]}</Link></Table.Cell>
                  {cellKeys.map((key) => {
                    return <Table.Cell>{listVal[key]}</Table.Cell>
                  })}
                </Table.Row>
              )
            })}
          </Table.Body>
        </Table>
      </div>
    )
  }
}
export default BasicDataTable;