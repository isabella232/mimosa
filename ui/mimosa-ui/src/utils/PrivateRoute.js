import React, { useContext } from 'react';
import { Route, Redirect } from 'react-router-dom';
import { AuthContext } from './auth'

const PrivateRoute = ({ component: RouteComponent, ...rest }) => {
  const {currentUser} = useContext(AuthContext);
  // Check if currentUser exists (i.e logged in)
  // if not redirect to login view
  return (
    <Route
    {...rest}
    render={routeProps => 
      !!currentUser ? (
        <RouteComponent {...routeProps} />
      ) : (
        <Redirect to={"/login"} />
      )
    }
    />
  );
}

export default PrivateRoute;