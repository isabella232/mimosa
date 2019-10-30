import React from 'react';
import 'semantic-ui-css/semantic.min.css'
import MimosaHeader from './components/MimosaHeader';
import DataTable from './components/DataTable';
import { Component } from 'react';
import { Button,
         Form,
         Grid,
         Header,
         Segment,
         Container,
         Divider,
         Icon,
         Sidebar
        } from 'semantic-ui-react';
import firebase, { googleProvider } from './components/firebase.js';
class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      loggedIn: false,
    }
  }


  LoginForm = () => (
    <Grid textAlign='center' style={{ height: '100vh' }} verticalAlign='middle'>
      <Grid.Column style={{ maxWidth: 450 }}>
        <Header as='h2' color='teal' textAlign='center'>
          <Icon name="cocktail" />Log-in to your account
      </Header>
        <Form size='large'>
          <Segment stacked>
            <Form.Input onChange={this.set_email} fluid icon='user' iconPosition='left' placeholder='E-mail address' />
            <Form.Input onChange={this.set_password}
              fluid
              icon='lock'
              iconPosition='left'
              placeholder='Password'
              type='password'
            />

            <Button onClick={this.handle_email_login} color='teal' fluid size='large'>
              Login
            </Button>

            <Divider />
            <Button onClick={this.handle_google_login} color='teal' fluid size='large'>
              Login with Google
          </Button>

          </Segment>
        </Form>
        {/* <Message>
          New to us? <a href='#'>Sign Up</a>
        </Message> */}
      </Grid.Column>
    </Grid>
  )


  set_email = (e, { value }) => {
    this.state.email = value
  }

  set_password = (e, { value }) => {
    this.state.password = value
  }

  handle_google_login = () => {
    firebase.auth().signInWithPopup(googleProvider).then((result) => {
      // var token = result.credential.accessToken;
      // var user = result.user;
      this.setState({
        loggedIn: true,
      });
    }).catch((error) => {
      alert(error)
      this.setState({
        loggedIn: false,
      })
    });
  }

  handle_email_login = () => {
    firebase.auth().signInWithEmailAndPassword(this.state.email, this.state.password).then((result) => {
      // var token = result.credential.accessToken;
      // var user = result.user;
      this.setState({
        loggedIn: true,
      });
    }).catch((error) => {
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
      alert(error)
    });
  }
  render() {
    const { loggedIn } = this.state;
    return (
      <div className="App">
      <MimosaHeader />
        {loggedIn ? (
          <Container>
            <Divider />
            <Button onClick={this.handle_logout}>Logout</Button>
            <DataTable />
          </Container>
        ) : (
            <Container>
              <this.LoginForm />
            </Container>
          )
        }
      </div>
    );
  }
}

export default App;
