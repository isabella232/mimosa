import React from 'react';
import { Menu, Button } from 'semantic-ui-react';
import firebase from '../utils/firebase.js';
import {withRouter} from 'react-router';
import {Link} from 'react-router-dom'; 

class NavMenu extends React.Component {
  constructor(props) {
    super(props);

    /**
     * "activePath" is the item in left nav that will be
     * seen as active
     */
    const {activePath} = this.props
    this.state = {
      activeItem: activePath,
    }
  }

  /**
   * logout of firebase auth, this will
   * automatically update AuthContext
   * and redirect to Login view
   */
  handle_logout = () => {
    firebase.auth().signOut().then(() => {
    }).catch((error) => {
      alert(error)
    });
  }

  //kinda irrelevent, need to re-adjust
  handleMenuNav = (e, { name }) => {
    this.setState({ activeItem: name })
  };
  render() {
    const {activeItem} = this.state;
    // Random Nav Items commented out, can be deleted
    // or used for later work
    return (
      <div>
        <Menu pointing vertical fixed inverted className="side-nav">
          <Menu.Item
            name='home'
            active={activeItem === 'home'}
            onClick={this.handleMenuNav}
            as={Link} to='/home'
          />
          {/* <Menu.Item
            name='sources'
            active={activeItem === 'sources'}
            onClick={this.handleMenuNav}
            as={Link} to='/home'
          /> */}
          <Menu.Item
            name='hosts'
            active={activeItem === 'hosts'}
            onClick={this.handleMenuNav}
            as={Link} to='/hosts'
          />
          {/* <Menu.Item
            name='vulns'
            active={activeItem === 'vulns'}
            onClick={this.handleMenuNav}
            as={Link} to='/hosts'
          /> */}
          {/* <Menu.Item
            name='task'
            active={activeItem === 'task'}
            onClick={this.handleMenuNav}
            as={Link} to='/tasks/task'
          /> */}
          <Button
            className='login'
            onClick={this.handle_logout}
            color='teal'
          >Logout
              </Button>
        </Menu>
      </div>
    )
  }
}
export default withRouter(NavMenu);