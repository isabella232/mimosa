import React, { Component } from 'react';
import { Container, Divider, Message } from 'semantic-ui-react';
import {NavMenu} from '../components';
import { withFirebase } from '../utils/Firebase';
import { withRouter } from 'react-router-dom';

class Home extends Component {
  constructor(props) {
    super(props);
    this.state = {
      data: '',
    }
  }
  //Collect the id from param route and use in firestore call
  componentDidMount() {
    const { nodeId } = this.props.match.params;
    console.log(this.props);
    this.pullTaskData(nodeId);
  }

  pullTaskData = (documentId) => {
    this.props.firebase.auth.currentUser.getIdTokenResult().then((token) => {
      var taskResult = ''
      this.props.firebase.app.firestore().collection("ws").doc(token.claims.defaultws).collection("tasks")
        .orderBy("timestamp", "desc").limit(1).get()
        .then(querySnapshot => {
          querySnapshot.forEach((doc) => {
            taskResult = JSON.stringify(doc.data(), null, 2);
          })
          this.setState({
            data: taskResult
          })
        })
    });
  }
  render() {
    const { data } = this.state;
    const { authUser } = this.props;
    return (
      <div>
        <NavMenu authUser={authUser} activePath="task" />
        <Container>
          <Divider />
          <Message>
            <Message.Header>Task Output</Message.Header>
            <pre>
              {data}
            </pre>
          </Message>
        </Container>
      </div>
    )
  }
}
export default withRouter(withFirebase(Home));