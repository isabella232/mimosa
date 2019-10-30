import React from 'react';
import { Header, Segment, Icon } from 'semantic-ui-react';


class MimosaHeader extends React.Component {
  render() {
    return (
      <Segment inverted>
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