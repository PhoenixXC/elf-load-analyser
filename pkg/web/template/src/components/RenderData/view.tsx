import React from 'react';

import { Panel } from 'rsuite';
import { Whisper, Popover } from 'rsuite';
import Markdown from 'react-markdown';
import gfm from 'remark-gfm';

import 'rsuite/lib/Popover/styles';
import './index.css';
import { ReportModelState } from '@/models/report';

// @ts-ignore
const renderListItem = ({ children }) => {
  let text = children[0].props.value;
  let idx = text.indexOf('=');
  let key = text.substring(0, idx), value = text.substring(idx, text.length);
  return (
    <li>
      <Whisper placement='top' trigger='click' speaker={
        <Popover title={key}><span>{value.substring(1, value.length)}</span></Popover>
      }>
        <span><b>{key}</b>{value}</span>
      </Whisper>
    </li>
  );
};

const ListRender = {
  listItem: renderListItem,
};

const RenderData = (props: ReportModelState) => {
  let data = props.data ? props.data : 'None data';
  return (
    <Panel className={'render-data'} header={props.name}>
      <Markdown source={data}
                renderers={ListRender}
                plugins={[gfm]}
                allowDangerousHtml={true} />
    </Panel>
  );
};

export default RenderData;
