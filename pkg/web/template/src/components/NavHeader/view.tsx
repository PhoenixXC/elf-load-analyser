import React from 'react';

import { Navbar, Nav, Icon } from 'rsuite';
import 'rsuite/lib/Nav/styles';
import 'rsuite/lib/Navbar/styles';
import 'rsuite/lib/Icon/styles';
import './index.css';

const NavHeader = () => {
  return (
    <Navbar appearance={'inverse'} id={'navbar'}>
      <Navbar.Header>
        <div className='title'>
          Analyse Report
        </div>
      </Navbar.Header>
      <Navbar.Body>
        <Nav pullRight href={'https://github.com/PhoenixXC/elf-load-analyser'}>
          <Nav.Item icon={<Icon icon='github' />}>Github</Nav.Item>
        </Nav>
      </Navbar.Body>
    </Navbar>
  );
};

export default NavHeader;
