import React from 'react';

import { Navbar, Nav, Icon, NavItemProps } from 'rsuite';
import 'rsuite/lib/Nav/styles';
import 'rsuite/lib/Navbar/styles';
import 'rsuite/lib/Icon/styles';
import './index.css';

const MyLink = React.forwardRef((props, ref) => {
  // @ts-ignore
  const { href, ...rest } = props;
  // @ts-ignore
  return (
    <a href={href} {...rest} />
  );
});

const NavLink = (props: JSX.IntrinsicAttributes & NavItemProps) =>
  <Nav.Item componentClass={MyLink} {...props} />;

const NavHeader = () => {
  return (
    <Navbar appearance={'inverse'} id={'navbar'}>
      <Navbar.Header>
        <div className='title'>
          Analyse Report
        </div>
      </Navbar.Header>
      <Navbar.Body>
        <Nav pullRight={true}>
          <NavLink href={'https://github.com/PhoenixXC/elf-load-analyser'}
                   target={'_blank'}
                   icon={<Icon icon='github' />}>Github</NavLink>
        </Nav>
      </Navbar.Body>
    </Navbar>
  );
};

export default NavHeader;
