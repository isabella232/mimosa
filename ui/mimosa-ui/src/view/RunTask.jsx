import React from 'react';
import NavMenu from '../components/NavMenu';
import { withFirebase } from '../utils/Firebase';
import { withRouter } from 'react-router-dom';
import {Container, Divider, Header, Form, Button, List} from 'semantic-ui-react';

class RunTask extends React.Component {
  render() {
    const { authUser, firebase } = this.props;
    const { wsid } = this.props.match.params;
    const hosts = this.props.location.state.response
    console.log('From data state ', hosts);
    return (
      <div>
        <NavMenu authUser={authUser} activePath="hosts" />
        <Container>
          <Header as="h1">Run Task</Header>
          <Divider />
          <Form>
            <Form.Field>
              <label>Task Name</label>
              <input placeholder="Facts" />
            </Form.Field>
            <Form.Field>
              <label>Params</label>
              <input placeholder="e.g. verbose"/>
            </Form.Field>
            <Form.Field>
              <label>Note</label>
              <input placeholder="note about task run" />
            </Form.Field>
            <Header as="h4">Hosts</Header>
            <List divided relaxed >
              {hosts && hosts.map((singleHost) => {
                return (
                  <List.Item>
                    <List.Icon name='server' verticalAlign='middle' />
                    <List.Content>
                      <List.Header>{singleHost}</List.Header>
                    </List.Content>
                  </List.Item>
                )
              })}
            </List>
            <Divider />
            <Button type="submit">Run</Button>
          </Form>
        </Container>
      </div>
    )
  }
}
export default withRouter(withFirebase(RunTask));