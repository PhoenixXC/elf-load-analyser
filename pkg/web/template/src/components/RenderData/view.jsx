import React from 'react';
import * as PropTypes from 'prop-types';

import { Panel } from 'rsuite';
import { Whisper, Popover } from 'rsuite';
import Markdown from 'react-markdown';
import gfm from 'remark-gfm';

import 'rsuite/lib/Popover/styles';
import './index.css';

const renderListItem = ({ children }) => {
  let text = children[0].props.value;
  return (
    <li>
      <Whisper placement='top' trigger='click' speaker={
        <Popover title=''>{text}</Popover>
      }>
        <span>{text}</span>
      </Whisper>
    </li>
  );
};

const ListRender = {
  listItem: renderListItem,
};

const RenderData = (props) => {
  return (
    <Panel className={'render-data'} header={props.name}>
      <Markdown source={props.render_data}
                renderers={ListRender}
                plugins={[gfm]}
                allowDangerousHtml={true} />
    </Panel>
  );
};

RenderData.propTypes = {
  name: PropTypes.string,
  render_data: PropTypes.string,
};

export default RenderData;
