import React, { Component } from 'react';
import {NavMenu} from '../components';
import { withFirebase } from '../utils/Firebase';
import { withRouter } from 'react-router-dom';
import {Container, Divider, Header, Form, Button, Message, List} from 'semantic-ui-react';

class RunTask extends Component {
  render() {
    const { authUser } = this.props;
    const { wsid } = this.props.match.params;
    var hasHosts;
    if(this.props.location.state && this.props.location.state.response.length > 0) {
      hasHosts = this.props.location.state.response
    } else {
      hasHosts = false;
    }
    console.log(hasHosts);
    console.log('From data state ', hasHosts);
    return (
      <div>
        <NavMenu authUser={authUser} workspace={wsid} activePath="hosts" />
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
            {hasHosts ?
              <div>
                <List divided relaxed >
                  {hasHosts && hasHosts.map((singleHost) => {
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
              </div>
              :
              <div>
                <Message
                  icon='server'
                  header='Cannot run task - no hosts selected'
                >
                  <p>Please select hosts from Host list and select Run Task</p>
                </Message>
              </div>

            }
            
          </Form>
        </Container>
      </div>
    )
  }
}
export default withRouter(withFirebase(RunTask));