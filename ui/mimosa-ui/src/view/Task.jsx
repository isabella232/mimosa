import React from 'react';
import { Container, Divider, Message } from 'semantic-ui-react';
import NavMenu from '../components/NavMenu';
import firebase from '../utils/firebase.js'

const db = firebase.firestore();

class Home extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      data: '',
    }
  }
  //Collect the id from param route and use in firestore call
  componentDidMount () {
    const { nodeId } = this.props.match.params;
    this.pullTaskData(nodeId);
  }
  pullTaskData = (documentId) => {
    var taskResult = ''
    db.collection("hosts").doc(documentId).collection("tasks")
      .orderBy("timestamp", "desc").limit(1).get()
      .then(querySnapshot => {
        querySnapshot.forEach((doc) => {
          taskResult = JSON.stringify(doc.data());
        })
        this.setState({
          data: taskResult
        })
      })
    
  }
  render() {
    const {data} = this.state;
    return (
      <div>
        <NavMenu activePath="task" />
        <Container>
          <Divider />
          <Message>
            <Message.Header>Task Output</Message.Header>
            <code>
              {data}
            </code>
          </Message>
        </Container>
      </div>
    )
  }
}
export default Home;