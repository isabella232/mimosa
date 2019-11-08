import React, {Component} from 'react';
import 'semantic-ui-css/semantic.min.css'
import './App.css'
import MimosaHeader from './components/MimosaHeader';
import Login from './view/Login';
import HostsView from './view/HostsView'
import Task from './view/Task';
import Home from './view/Home';
import {
  BrowserRouter as Router,
  Route,
} from "react-router-dom";
import { withFirebase } from './utils/Firebase';
import history from './utils/history';

// The router will only allow access to login for 
// users that have not logged in. 
// Also history is passed in to be used by components e.g. NavMenu
class App  extends Component {
  constructor(props) {
    super(props);
    this.state = {
      authUser: null,
    };
  }
  // get user and set as state if logged into firebase
  componentDidMount() {
    this.listener = this.props.firebase.auth.onAuthStateChanged(authUser => {
      authUser
      ? this.setState({ authUser })
      : this.setState({ authUser: null });
    });
  }
  // when leaving/unmounting App, remote the listener
  // done to avoid potential performance issues
  componentWillUnmount() {
    this.listener();
  }

  render() {
    const {firebase} = this.props;
    return (
        <div>
          <MimosaHeader />
          <Router history={history}>
            <div>
              <Route path="/login" render={() => <Login authUser={this.state.authUser} history={history} />} firebase={firebase} />
              <Route path="/home" authUser={this.state.authUser} render={() => <Home authUser={this.state.authUser} history={history} />} firebase={firebase} />
              <Route path="/hosts" render={() => <HostsView authUser={this.state.authUser} firebase={firebase} />} />
              <Route path="/:nodeId/task" render={() => <Task authUser={this.state.authUser} firebase={firebase} />} />
            </div>
          </Router>
        </div>
    )
  }
}
export default withFirebase(App);
