import React from 'react';

import { Panel, PanelGroup, Affix } from 'rsuite';
import 'rsuite/lib/PanelGroup/styles';
import 'rsuite/lib/Panel/styles';

import { NavHeader } from '../NavHeader/index';
import { RenderData } from '../RenderData/index';

import EnvData from './data/env.json';
import FileData from './data/fileFormat.json';

import './index.css';
import { RenderDataType } from '@/components/RenderData/type/RenderData';

const Container = () => {
  const envRenderData = EnvData as RenderDataType;
  const fileRenderData = FileData as RenderDataType;
  return (
    <div className={'container'}>
      <Affix>
        <NavHeader />
      </Affix>
      <div className='content'>
        <PanelGroup>
          <RenderData {...envRenderData} />
          <RenderData {...fileRenderData} />
          <Panel header='加载过程'>
            <p>加载过程</p>
          </Panel>
        </PanelGroup>
      </div>
    </div>
  );
};

export default Container;
