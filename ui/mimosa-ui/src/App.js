import React, {Component} from 'react';
import 'semantic-ui-css/semantic.min.css'
import './App.css'
import MimosaHeader from './components/MimosaHeader';
import Login from './view/Login';
import HostsView from './view/HostsView'
import RunTask from './view/RunTask';
import Home from './view/Home';
import Workspaces from './view/Workspaces'
import RunContext from './view/RunContext.jsx'
import RunDetail from './view/RunDetail.jsx'
import HostDetailView from './view/HostDetailView.jsx';
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
              <Route exact path="/login" render={() => <Login authUser={this.state.authUser} history={history} />} firebase={firebase} />
              <Route exact path="/ws/:wsid/home" authUser={this.state.authUser} render={() => <Home authUser={this.state.authUser} history={history} firebase={firebase}  />}/>
              <Route exact path="/ws" render={() => <Workspaces authUser={this.state.authUser} firebase={firebase}/>} />
              <Route exact path="/ws/:wsid/hosts" render={() => <HostsView authUser={this.state.authUser} firebase={firebase} />} />
              <Route exact path="/ws/:wsid/host/:hostid" render={() => <HostDetailView authUser={this.state.authUser} firebase={firebase} />} />
              <Route exact path="/ws/:wsid/run-context" render={() => <RunContext authUser={this.state.authUser} firebase={firebase} />} />
              <Route exact path="/ws/:wsid/run-task" render={() => <RunTask authUser={this.state.authUser} firebase={firebase} />} />
              <Route exact path="/ws/:wsid/run/:runid" render={() => <RunDetail authUser={this.state.authUser} firebase={firebase} />} />
            </div>
          </Router>
        </div>
    )
  }
}
export default withFirebase(App);
