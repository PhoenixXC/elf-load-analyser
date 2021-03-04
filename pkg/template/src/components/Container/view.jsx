import React from 'react';

import { Panel, PanelGroup, Affix } from 'rsuite';
import 'rsuite/lib/PanelGroup/styles';
import 'rsuite/lib/Panel/styles';

import { NavHeader } from '../NavHeader';

import './index.css';

const Container = () => {
  return (
    <div>
      <Affix>
        <NavHeader />
      </Affix>
      <div className='content'>
        <PanelGroup>
          <Panel header='Panel 1'>
            <p>panel 1</p>
          </Panel>
          <Panel header='Panel 2'>
            <p>panel 2</p>
          </Panel>
          <Panel header='Panel 3'>
            <p>panel 3</p>
          </Panel>
        </PanelGroup>
      </div>
    </div>
  );
};

export default Container;
