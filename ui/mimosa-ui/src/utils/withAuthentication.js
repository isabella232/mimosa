import React, { Component } from 'react';
import { withRouter } from 'react-router-dom';
import { withFirebase } from '../utils/Firebase';
import { compose } from 'recompose';

// This higher level component SHOULD allow us to redirect
// if a user is not signed in, however it seems to be rendered
// too late... so either it needs used earlier, or we readjust
// how we are handling routes
const withAuthentication = condition => Component => {
  class WithAuthentication extends Component {
    componentDidMount() {
      console.log('withAuth', condition);
      this.listener = this.props.firebase.auth.onAuthStateChanged(
        authUser => {
          if (!condition) {
            this.props.history.push('/login');
          }
        }
      )
      this.listener();
    }
    render() {
      return <Component {...this.props } />
    }
  }
  return compose(
    withRouter,
    withFirebase,
  )(WithAuthentication);
}

export default withAuthentication;