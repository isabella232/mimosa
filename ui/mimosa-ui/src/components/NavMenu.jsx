import React from 'react';
import { Menu, Button } from 'semantic-ui-react';
import {withRouter} from 'react-router';
import { withFirebase } from '../utils/Firebase';

import {Link} from 'react-router-dom'; 

class NavMenu extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      activeItem: '',
    }

    const { activePath } = this.props
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
    const {authUser} = this.props;
    console.log(authUser);
    if (!authUser) {
      this.props.history.push('/login');
    }
    // Random Nav Items commented out, can be deleted
    // or used for later work
    return (
      <div>
        {authUser ?
          <Menu pointing vertical fixed inverted className="side-nav">
            <Menu.Item
              name='home'
              active={activeItem === 'home'}
              onClick={this.handleMenuNav}
              as={Link} to='/home'
            />
            <Menu.Item
              name='hosts'
              active={activeItem === 'hosts'}
              onClick={this.handleMenuNav}
              as={Link} to='/hosts'
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
              as={Link} to='/home'
            />
            <Menu.Item
              disabled
              name='hosts'
              active={activeItem === 'hosts'}
              onClick={this.handleMenuNav}
              as={Link} to='/hosts'
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