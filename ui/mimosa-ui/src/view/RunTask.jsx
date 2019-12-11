import React, { Component } from 'react';
import { NavMenu } from '../components';
import { withFirebase } from '../utils/Firebase';
import { withRouter } from 'react-router-dom';
import { Container, Divider, Header, Form, Button, Message, Icon } from 'semantic-ui-react';

class RunTask extends Component {
  constructor(props) {
    super(props);
    this.state = {
      isLoading: false,
      isTaskTriggered: false,
      isError: false
    }
  }
  // Call cloud function, since we don't expect result we don't do anything
  callCloudFunction = (functionName, hostid) => {
    var {wsid} = this.props.match.params;
    this.setState({
      isLoading: true,
      isTaskTriggered: true,
    });
    if (this.props.firebase.auth.currentUser) {
      this.props.firebase.auth.currentUser.getIdToken().then((idToken) => {
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
          this.setState({
            isLoading: false,
          });
        })
          .catch(error => {
            console.error('Error during Mimosa:', error);
            this.setState({
              isLoading: false,
              isError: true,
            });
          });
      }).catch(function (error) {
        console.error('Error during Mimosa:', error);
        this.setState({
          isLoading: false,
          isError: true,
        });
      });
    }
  }

  viewHost = (hostid) => {
    const { wsid } = this.props.match.params;
    this.props.history.push('/ws/' + wsid + '/host/' + hostid);
  }

  handleChange = (e, {name, value}) => this.setState({ [name]: value})

  render() {
    const { authUser } = this.props;
    const { wsid } = this.props.match.params;
    const { isLoading, isTaskTriggered, isError, name, param } = this.state;

    var hasHosts, docId;
    if (this.props.location.state && this.props.location.state.response.length > 0) {
      hasHosts = this.props.location.state.response
      docId = this.props.location.state.doc
    } else {
      hasHosts = false;
    }
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
            <Form.Input
              label="Task name"
              placeholder="facts"
              name="facts"
              value={name}
              onChange={this.handleChange}
            />
            <Form.Input
              label="Params"
              placeholder="e.g. verbose"
              name="param"
              value={param}
              onChange={this.handleChange}
            />
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
                {isTaskTriggered ? (
                  <Button
                    color={isError ? 'red' : 'teal'}
                    loading={isLoading}
                  >
                    {isError ? "Error" : 'Complete'}
                  </Button>
                ) : (
                  <Button
                    color="purple"
                    onClick={() => this.callCloudFunction('api/v1/runtask', docId)}
                  >
                    <Icon name="play circle outline" />
                    Run Task
                  </Button>
                )}
                <Button
                  color="teal"
                  onClick={() => this.viewHost(docId)}
                >
                  View Host Data
                </Button>
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