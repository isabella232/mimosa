import React from 'react';
import { Container, Divider, Icon, Header } from 'semantic-ui-react';
import NavMenu from '../components/NavMenu';

class Home extends React.Component {
  render() {
    return (
      <div>
        <NavMenu activePath="home"/>
        <Container>
          <Header as="h1">
            <Icon name="cocktail" />
            Welcome to Mimosa
          </Header>
          <Divider />
          <p>Project Mimosa is Puppet working to do more cloudy...stuff</p>
        </Container>
      </div>
    )
  }
}
export default Home;