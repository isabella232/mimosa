import React, { Component } from 'react';
import { NavMenu } from '../components';
import { withFirebase } from '../utils/Firebase';
import { withRouter } from 'react-router-dom';
import { Container, Divider, Header, Form, Button, Message, Dropdown } from 'semantic-ui-react';

class RunTask extends Component {
  // Call cloud function, since we don't expect result we don't do anything
  callCloudFunction = (functionName, hostid) => {
    if (this.props.firebase.auth.currentUser) {
      var wsid = this.props.match.params.wsid;
      this.props.firebase.auth.currentUser.getIdToken().then(function (idToken) {
        // FIXME - ACCESS TOKEN SHOULD BE ADDED AS A BEARER TOKEN
        fetch('https://mimosa-esp-tfmdd2vwoq-uc.a.run.app/' + functionName + "?access_token=" + idToken, {
          method: 'POST',
          mode: 'cors',
          cache: 'no-cache',
          headers: {
            'Content-Type': 'application/json'
          },
          redirect: 'follow',
          referrer: 'no-referrer',
          body: JSON.stringify({ "workspace": wsid, "id": hostid })
        }).then(response => {
          // console.log(response.status)
          // console.log(response.text())
          this.props.history.push('/ws/' + wsid + '/host/' + hostid);
        })
          .catch(error => {
            console.error('Error during Mimosa:', error);
          });
      }).catch(function (error) {
        console.error('Error during Mimosa:', error);
      });
    }
  }

  render() {
    const { authUser } = this.props;
    const { wsid } = this.props.match.params;
    var hasHosts, docId;
    if (this.props.location.state && this.props.location.state.response.length > 0) {
      hasHosts = this.props.location.state.response
      docId = this.props.location.state.doc
    } else {
      hasHosts = false;
    }
    console.log(hasHosts);
    console.log('From data state ', hasHosts);
    const option = [
      { text: hasHosts, value: hasHosts }
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
              <input value="facts" />
            </Form.Field>
            <Form.Field disabled>
              <label>Params</label>
              <input placeholder="e.g. verbose" />
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
                <Button onClick={this.callCloudFunction('api/v1/runtask', docId)} type="submit">Run</Button>
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