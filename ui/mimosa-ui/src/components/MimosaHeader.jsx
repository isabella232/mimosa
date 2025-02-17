import React, { Component } from 'react';
import { Header, Segment, Icon } from 'semantic-ui-react';

class MimosaHeader extends Component {
  render() {
    return (
      <Segment inverted className="mimosaHeader">
        <Header inverted color="grey" as='h1' textAlign='center'>
          <Header.Content>
            <Icon name="cocktail" />
            mimosa
          </Header.Content>
        </Header>
      </Segment>
    )
  }
}
export default MimosaHeader;