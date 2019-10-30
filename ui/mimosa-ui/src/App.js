import React from 'react';
import 'semantic-ui-css/semantic.min.css'
import MimosaHeader from './components/MimosaHeader';
import DataTable from './components/DataTable';
import { Component } from 'react';
import { Container, Divider, Button } from 'semantic-ui-react';
import firebase, { googleProvider } from './components/firebase.js';
class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      loggedIn: false,
    }
  }

  handle_login = () => {
    firebase.auth().signInWithPopup(googleProvider).then((result) => {
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

  handle_email_login = () => {
    firebase.auth().signInWithEmailAndPassword("alice@example.com", "alicealice").then((result) => {
      // var token = result.credential.accessToken;
      // var user = result.user;
      this.setState({
        loggedIn: true,
      });
    }).catch((error) => {
      alert("Error during signin")
      alert(error)
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
          <Container>
            <Button onClick={this.handle_logout}>Logout</Button>
            <DataTable />
          </Container>
        ) : (
            <Container>
              <Button onClick={this.handle_login}>Login with Google</Button>
              <Button onClick={this.handle_email_login}>Login with Email</Button>
            </Container>
          )}
      </div>
    );
  }
}

export default App;
