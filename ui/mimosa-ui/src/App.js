import React from 'react';
import 'semantic-ui-css/semantic.min.css'
import MimosaHeader from './components/MimosaHeader';
import DataTable from './components/DataTable';
import { Component } from 'react';
import { Container, Divider, Button } from 'semantic-ui-react';
import firebase, { firestore, provider } from './components/firebase.js';
class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      loggedIn: false,
    }
  }

  handle_login = () => {
    firebase.auth().signInWithPopup(provider).then((result) => {
      var token = result.credential.accessToken;
      var user = result.user;
      this.setState({
        loggedIn: true,
      });
    }).catch((error) => {
      alert("Error during signin")
      this.setState({
        loggedIn: false,
      })
    });
  }

  handle_logout = () => {
    firebase.auth().signOut().then(() => {
      this.setState({
        loggedIn: false
      });
    }).catch((error) => {
      alert("Error logging out");
    });
  }
  render() {
    const { loggedIn } = this.state;
    return (
      <div className="App">
        <MimosaHeader />
        <Divider />
        {loggedIn ? (
          <Button onClick={this.handle_logout}>Logout</Button>
        ) : (
          <Button onClick={this.handle_login}>Login</Button>
        )}
        {loggedIn ? (
          <Container>
            <DataTable />
          </Container>
        ) : (
          <Container>
            <p>Please Sign in to view mimosa data</p>
          </Container>
        )}
      </div>
    );
  }
}

export default App;
