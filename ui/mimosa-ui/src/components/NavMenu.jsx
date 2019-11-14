import React, { Component } from 'react';
import { Menu, Button } from 'semantic-ui-react';
import {withRouter} from 'react-router';
import { withFirebase } from '../utils/Firebase';

import {Link} from 'react-router-dom'; 

class NavMenu extends Component {
  constructor(props) {
    super(props);
    this.state = {
      activeItem: '',
    }
  }

  handle_logout = () => {
    this.props.firebase.auth.signOut().then(() => {
      this.props.history.push('/login')
    }).catch((error) => {
      alert(error)
    });
  }

  //kinda irrelevent, need to re-adjust
  handleMenuNav = (e, { name }) => {
    this.setState({ activeItem: name })
  };

  goLogin = () => {
    this.props.history.push('/login');
  }

  render() {
    const {activeItem} = this.state;
    const {authUser, workspace} = this.props;
    if (!authUser) {
      this.props.history.push('/login');
    }
    var prefix = `/ws/${workspace}`
    // Random Nav Items commented out, can be deleted
    // or used for later work
    return (
      ///ws/:wsid/hosts
      <div>
        {authUser ?
          <Menu pointing vertical fixed inverted className="side-nav">
            <Menu.Item
              name='home'
              active={activeItem === 'home'}
              onClick={this.handleMenuNav}
              as={Link} to={prefix + '/home'}
            />
            <Menu.Item
              name='host list'
              active={activeItem === 'host list'}
              onClick={this.handleMenuNav}
              as={Link} to={prefix + '/hosts'}
            />
            <Menu.Item
              name='run context'
              active={activeItem === 'run context'}
              onClick={this.handleMenuNav}
              as={Link} to={prefix + '/run-context'}
            />
            <Menu.Item
              name='start run'
              active={activeItem === 'start run'}
              onClick={this.handleMenuNav}
              as={Link} to={prefix + '/run-task'}
            />
            <Button
              className='login'
              onClick={this.handle_logout}
              color='teal'>
              Logout
            </Button>
          </Menu> 
          :
          <Menu pointing vertical fixed inverted className="side-nav">
            <Menu.Item
              name='home'
              active={activeItem === 'home'}
              onClick={this.handleMenuNav}
              as={Link} to={prefix + '/home'}
            />
            <Button
              className='login'
              onClick={this.goLogin}
              color='teal'>
              Login
            </Button>
          </Menu>
        }
      </div>
    )
  }
}
export default withRouter(withFirebase(NavMenu));