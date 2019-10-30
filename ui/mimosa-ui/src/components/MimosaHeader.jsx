import React from 'react';
import { Header, Segment } from 'semantic-ui-react';


class MimosaHeader extends React.Component{
  render() {
    return (
      <Segment inverted>
        <Header inverted color="grey" as='h1' textAlign='center'>
          <Header.Content>Say hello to Mimosa</Header.Content>
        </Header>
      </Segment>
    )
  } 
}
export default MimosaHeader;