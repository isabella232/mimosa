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
    const { wsid, runid } = this.props.match.params;
    console.log(this.props);
    this.pullTaskData(wsid, runid);
  }

  pullTaskData = (workspace, documentId) => {
    if (this.props.firebase.auth.currentUser) {
      this.props.firebase.auth.currentUser.getIdTokenResult().then((token) => {
        var taskResult = ''
        this.props.firebase.app.firestore().collection("ws").doc(workspace).collection("tasks").doc(documentId).get()
          .then(querySnapshot => {
            console.log(querySnapshot.data());
            var stringOutput = JSON.stringify(querySnapshot.data());
            this.setState({
              data: stringOutput
            })
          })
      });
    }
  }
  render() {
    const { data } = this.state;
    const { authUser } = this.props;
    const { wsid } = this.props.match.params;
    return (
      <div>
        <NavMenu authUser={authUser} workspace={wsid} activePath="task" />
        <Container>
          <Divider />
          <Message style={{overflowWrap: "break-word"}}>
            <Message.Header>Task Output</Message.Header>
              {data}
          </Message>
        </Container>
      </div>
    )
  }
}
export default withRouter(withFirebase(Home));