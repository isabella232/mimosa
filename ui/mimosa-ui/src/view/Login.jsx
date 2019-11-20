import React, { Component } from 'react';
import { Header, Segment, Icon, Grid, Divider, Form, Button, Container } from 'semantic-ui-react';
import { withFirebase } from '../utils/Firebase';
import { withRouter } from 'react-router-dom';
import cookie from 'react-cookies';

const INITIAL_STATE = {
  email: '',
  password: '',
}

class Login extends Component {
  constructor(props) {
    super(props);
    this.state = {
      ...INITIAL_STATE,
    }
  }
  
  /**
   * Need to componentize this correctly
   * to allow for use of setState AND
   * currentUserCheck
   */
  set_email = (e, { value }) => {
    this.setState({
      email: value
    })
  }
  set_password = (e, { value }) => {
    this.setState({
      password: value
    })
  }
  googleLogin = () => {
    const googleProvider = this.props.firebase.googleProv
    this.props.firebase.auth.signInWithPopup(googleProvider).then((result) => {
      if (result && result.user.email) {
        console.log(result.user.email);
        cookie.save("userEmail", "loggedIn", { path: '/' });
        this.props.history.push('/ws')
      }
    }).catch((error) => {
      alert(error);
    });
  }


  emailLogin = () => {
    var { email, password } = this.state;
    this.props.firebase.auth.signInWithEmailAndPassword(email, password).then((result) => {
      if(result && result.user.email) {
        console.log("This one now! ", result.user.email);
        cookie.save("userEmail", "loggedIn", { path: '/' });
        cookie.loadAll();
        this.setState({ ...INITIAL_STATE });
        this.props.history.push('/ws')
      }
    }).catch((error) => {
      alert(error)
      this.setState({ ...INITIAL_STATE });
    });
  }

  render() {
    return (
      <Container>
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
                <Button onClick={this.emailLogin} color='teal' fluid size='large'>
                  Login
                </Button>

                <Divider />
                <Button onClick={this.googleLogin} color='teal' fluid size='large'>
                  Login with Google
                </Button>

              </Segment>
            </Form>
          </Grid.Column>
        </Grid>
      </Container>
    )
  }
}
export default withRouter(withFirebase(Login));