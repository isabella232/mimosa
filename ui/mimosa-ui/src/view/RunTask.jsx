import React, { Component } from 'react';
import {NavMenu} from '../components';
import { withFirebase } from '../utils/Firebase';
import { withRouter } from 'react-router-dom';
import {Container, Divider, Header, Form, Button, Message, Dropdown} from 'semantic-ui-react';

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
    const option = [
      {text: hasHosts[0], value: hasHosts[0]}
    ]
    return (
      <div>
        <NavMenu authUser={authUser} workspace={wsid} activePath="hosts" />
        <Container>
          <Header as="h1">Run Task</Header>
          <Divider />
          <Form>
            <Form.Field>
              <label>Task Name</label>
              <input value="facts"/>
            </Form.Field>
            <Form.Field disabled>
              <label>Params</label>
              <input placeholder="e.g. verbose"/>
            </Form.Field>
            <Form.Field disabled>
              <label>Note</label>
              <input placeholder="note about task run" />
            </Form.Field>
            <Header as="h4">Hosts</Header>
            {hasHosts ?
              <div>
                <Form.Dropdown
                  options={option}
                  defaultValue={option[0].value}
                />
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