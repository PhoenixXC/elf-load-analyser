import React from 'react';

import { Panel, PanelGroup, Affix } from 'rsuite';
import 'rsuite/lib/PanelGroup/styles';
import 'rsuite/lib/Panel/styles';

import { NavHeader } from '../NavHeader';

import './index.css';

const PanelList = (props) => {
  const pannelNum = props.renderList;
}

const Container = (props) => {
  return (
    <div className={'container'}>
      <Affix>
        <NavHeader />
      </Affix>
      <div className='content'>
        <PanelGroup>
          <Panel header='系统环境'>
            <p>系统环境</p>
          </Panel>
          <Panel header='文件格式'>
            <p>文件格式</p>
          </Panel>
          <Panel header='加载过程'>
            <p>加载过程</p>
          </Panel>
        </PanelGroup>
      </div>
    </div>
  );
};

export default Container;
