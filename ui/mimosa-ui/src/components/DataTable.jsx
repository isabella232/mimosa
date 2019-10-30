import _ from 'lodash';
import React from 'react';
import firebase, { firestore, provider, functions } from './firebase.js';
import { Table, Button, Modal } from 'semantic-ui-react';

var db = firebase.firestore();

class DataTable extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      data: [{}],
      index: [],
      outputArray: [],
      cloudFunctionResult: '',
    }
  }
  pullData = () => {
    var stagingArray = [],
        outputArray = [];
    db.collection("hosts").get().then((querySnapshot) => {
      querySnapshot.forEach((doc) => {
        stagingArray.push(doc.data());
        outputArray.push('');
      });
      console.log("sorting firebase data")
      stagingArray = _.sortBy(stagingArray, ["state", "source", "public_dns"])
      this.setState({
        data: stagingArray,
        output: outputArray,
      });
    });
  }
  callCloudFunction = (functionName, index) => {
    var cf = firebase.functions().httpsCallable(functionName);
    cf().then((result) => {
      console.log(result.data, index);
      var stagedOutput = this.state.outputArray;
      stagedOutput[index] = result.data;
      this.setState({
        cloudFunctionResult: result.data, 
        outputArray: stagedOutput,
      });
    }).catch((error) => {
      console.log(error);
    })
  }
  componentDidMount() {
    //fakeData to be used for styling, visual fixes, rather than hitting DB
    let fakeData = [
    //   {
    //     name: "onoijsaofjasmdfl;jasdofl;ask;dojasdfje",
    //     public_dns: "12234234590u320495u2039u4534",
    //     public_ip: "0.0.0.1",
    //     since: {
    //       seconds: 1234,
    //     },
    //     source: "me, myself and i",
    //     state: "running",
    //   },
    //   {
    //     name: "ksdfoijasd;fmas;odfj;ofj",
    //     public_dns: "123psdfosjdf4",
    //     public_ip: "0.0.0.1:/255",
    //     since: {
    //       seconds: 1234,
    //     },
    //     source: "vmpooler",
    //     state: "terminated",
    //   },
    //   {
    //     name: "asdc;amsd;kcnaskcn",
    //     public_dns: "1234",
    //     public_ip: "0.0.0.1",
    //     since: {
    //       seconds: 1234,
    //     },
    //     source: "bwabeabeaa",
    //     state: "running",
    //   },
    //   {
    //     name: "sdfasjo;fdjoais;djfo",
    //     public_dns: "1234",
    //     public_ip: "0.0.0.1",
    //     since: {
    //       seconds: 1234,
    //     },
    //     source: "sdfasdfasdfasdf",
    //     state: "terminated",
    //   }
    ]
    let pulledData = this.pullData();
    this.setState({ data: pulledData, });
  }
  render() {
    var { data, cloudFunctionResult, outputArray } = this.state;
    return (
      <Table className="table">
        <Table.Header>
          <Table.Row>
            <Table.HeaderCell>Id</Table.HeaderCell>
            <Table.HeaderCell>Domain name</Table.HeaderCell>
            <Table.HeaderCell>IP Address</Table.HeaderCell>
            <Table.HeaderCell>Source</Table.HeaderCell>
            <Table.HeaderCell>State</Table.HeaderCell>
            <Table.HeaderCell>Run Task</Table.HeaderCell>
            <Table.HeaderCell>Run Result</Table.HeaderCell>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {data && data.map((listVal, index) => {
            var rowState, showButton;
            if (listVal.state === 'terminated') {
              rowState = false;
              showButton = false;
            } else {
              rowState = true;
              showButton = true;
            }
            return (
              <Table.Row error={!rowState} positive={rowState} key={index}>
                <Table.Cell>{listVal.name}</Table.Cell>
                <Table.Cell>{listVal.public_dns}</Table.Cell>
                <Table.Cell>{listVal.public_ip}</Table.Cell>
                <Table.Cell>{listVal.source}</Table.Cell>
                <Table.Cell>{listVal.state}</Table.Cell>
                {showButton ? (
                  <Table.Cell>
                    <Button inverted color='violet' onClick={() => this.callCloudFunction('RunTask', index)}>Run Task</Button>
                  </Table.Cell>
                ) : (
                    <Table.Cell>
                      -
                  </Table.Cell>
                  )}
                {showButton ? (
                  <Table.Cell>
                    <p>{outputArray[index]}</p>
                  </Table.Cell>
                ) : (
                  <Table.Cell>
                    -
                  </Table.Cell>
                )}
              </Table.Row>
            )
          })}
        </Table.Body>
      </Table>
    )
  }
}
export default DataTable;