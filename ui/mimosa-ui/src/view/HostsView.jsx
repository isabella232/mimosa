import React from 'react';
import { Container, Divider } from 'semantic-ui-react';
import HostDataTable from '../components/HostDataTable';
import NavMenu from '../components/NavMenu';

class HostsView extends React.Component {
  render() {
    return (
      <div>
        <NavMenu activePath="hosts"/>
        <Container>
          <Divider />
          <HostDataTable />
        </Container>
      </div>
    )
  }
}
export default HostsView;