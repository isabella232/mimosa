import React from 'react';
import firebase, { firestore, provider } from './firebase.js';
import { Table, Button } from 'semantic-ui-react';

var db = firebase.firestore();

class DataTable extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      data: [{}],
    }
  }
  pullData = () => {
    var stagingArray = [];
    db.collection("hosts").get().then((querySnapshot) => {
      querySnapshot.forEach((doc) => {
        stagingArray.push(doc.data());
      });
      this.setState({
        data: stagingArray,
      });
    });
  }
  componentDidMount(){
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
    // let pulledData = this.pullData();
    this.setState({ data: fakeData,});
  }
  render() {
    var {data} = this.state;
    return (
      <Table className="table">
        <Table.Header>
          <Table.Row>
            <Table.HeaderCell>Id</Table.HeaderCell>
            <Table.HeaderCell>Domain name</Table.HeaderCell>
            <Table.HeaderCell>IP Address</Table.HeaderCell>
            <Table.HeaderCell>Source</Table.HeaderCell>
            <Table.HeaderCell>State</Table.HeaderCell>
            <Table.HeaderCell>Run CF</Table.HeaderCell>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {data && data.map((listVal, index) => {
            var rowState, showButton;
            if (listVal.state === 'terminated') {
              rowState=false;
              showButton=false;
            } else {
              rowState=true;
              showButton=true;
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
                    <Button inverted color='violet'>Run Bolt</Button>
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