import React from 'react';
import { Container, Divider, Message } from 'semantic-ui-react';
import NavMenu from '../components/NavMenu';
import { withFirebase } from '../utils/Firebase';
import { withRouter } from 'react-router-dom';

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
    console.log(this.props);
    this.pullTaskData(nodeId);
  }
  pullTaskData = (documentId) => {
    var taskResult = ''
    this.props.firebase.app.firestore().collection("hosts").doc(documentId).collection("tasks")
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
    const {authUser} = this.props;
    return (
      <div>
        <NavMenu authUser={authUser} activePath="task" />
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
export default withRouter(withFirebase(Home));