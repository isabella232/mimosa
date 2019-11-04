import React from 'react';
import 'semantic-ui-css/semantic.min.css'
import './App.css'
import MimosaHeader from './components/MimosaHeader';
import PrivateRoute from './utils/PrivateRoute';
import Login from './view/Login';
import HostsView from './view/HostsView'
import Task from './view/Task';
import Home from './view/Home';
import {
  BrowserRouter as Router,
  Route,
} from "react-router-dom";
import { AuthProvider } from './utils/auth';
import history from './utils/history';

// The router will only allow access to login for 
// users that have not logged in. 
// Also history is passed in to be used by components e.g. NavMenu
const App = () => {
  return (
    <div>
      <MimosaHeader />
      <AuthProvider>
        <Router history={history}>
          <div>
            <PrivateRoute path="/home" component={Home} />
            <PrivateRoute path="/hosts" component={HostsView} />
            <PrivateRoute path="/:nodeId/task" component={Task} />
            <Route path="/login" component={Login} />
            <Route path="/" component={Login} />
          </div>
        </Router>
      </AuthProvider>
    </div>
  )
}

export default App;
