import React from 'react';

import { Timeline, Tooltip, Whisper } from 'rsuite';
import Markdown from 'react-markdown';

import './index.css';

import { RenderDataType, ErrorStatus } from '../RenderData/type/RenderData';

const Status = (props: RenderDataType) => {
  return (
    <Whisper placement={'top'} trigger={'hover'}
             speaker={<Tooltip>{props.desc}</Tooltip>}>
      <div className={'dot'} />
    </Whisper>
  );
};

const Item = (props: RenderDataType) => {
  const data = props.data ? props.data : 'No data';
  return (
    <Timeline.Item time={props.time} dot={<Status {...props} />}>
      <div className={'meta-item'}>{props.name}</div>
      <Markdown source={data} />
    </Timeline.Item>
  );
};

const TimelineData = (props: RenderDataType) => {
  const dataList = props.dataList;
  return (
    <Timeline align={'left'}>

    </Timeline>
  )
};

export default TimelineData;
