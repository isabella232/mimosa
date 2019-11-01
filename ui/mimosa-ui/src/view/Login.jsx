import React, {useContext} from 'react';
import { Header, Segment, Icon, Grid, Divider, Form, Button, Container } from 'semantic-ui-react';
import firebase, { googleProvider } from '../utils/firebase.js';
import { Redirect } from 'react-router-dom';
import { withRouter } from 'react-router-dom';
import {AuthContext} from '../utils/auth.js';


const Login = ({history}) => {
  var email = '',
      password = '';
  
  /**
   * Need to componentize this correctly
   * to allow for use of setState AND
   * currentUserCheck
   */
  const set_email = (e, { value }) => {
    email = value;
    // this.setState({
    //   email: value
    // })
  }
  const  set_password = (e, { value }) => {
    password = value;
    // this.setState({
    //   password: value
    // })
  }
  const handle_google_login = () => {
    firebase.auth().signInWithPopup(googleProvider).then((result) => {
      history.push("/hosts");
    }).catch((error) => {
      alert(error);
    });
  }

  const handle_email_login = () => {
    firebase.auth().signInWithEmailAndPassword(email, password).then(() => {
    }).catch((error) => {
      alert(error)
    });
  }

  const { currentUser } = useContext(AuthContext);
  if (currentUser) {
    return <Redirect to='/home' />;
  }

  return (
    <Container>
      <Grid textAlign='center' style={{ height: '100vh' }} verticalAlign='middle'>
        <Grid.Column style={{ maxWidth: 450 }}>
          <Header as='h2' color='teal' textAlign='center'>
            <Icon name="cocktail" />Log-in to your account
    </Header>
          <Form size='large'>
            <Segment stacked>
              <Form.Input onChange={set_email} fluid icon='user' iconPosition='left' placeholder='E-mail address' />
              <Form.Input onChange={set_password}
                fluid
                icon='lock'
                iconPosition='left'
                placeholder='Password'
                type='password'
              />
              <Button onClick={handle_email_login} color='teal' fluid size='large'>
                Login
          </Button>

              <Divider />
              <Button onClick={handle_google_login} color='teal' fluid size='large'>
                Login with Google
        </Button>

            </Segment>
          </Form>
        </Grid.Column>
      </Grid>
    </Container>
  )
}
export default withRouter(Login);