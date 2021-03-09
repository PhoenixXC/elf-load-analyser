import React from 'react';

import { Panel, Timeline, Tooltip, Whisper } from 'rsuite';
import Markdown from 'react-markdown';

import 'rsuite/lib/Timeline/styles';
import 'rsuite/lib/Tooltip/styles';
import './index.css';

import { ReportModelState, ErrorStatus } from '@/models/report';

const Status = (props: ReportModelState) => {
  return (
    <Whisper placement={'top'} trigger={'hover'}
             speaker={<Tooltip>{props.desc}</Tooltip>}>
      <div className={`dot ${props.status === ErrorStatus ? 'error' : ''}`} />
    </Whisper>
  );
};

const Item = (props: ReportModelState) => {
  const data = props.data ? props.data : 'No data';
  return (
    <Timeline.Item dot={<Status {...props} />}>
      <div className={'meta-item'}>
        <span className={'time'}> {props.time}</span>
        <span className={'name'}> {props.name}</span>
      </div>
      <Markdown source={data} />
    </Timeline.Item>
  );
};

const TimelineData = (props: ReportModelState) => {
  const dataListItems = props.dataList
    ? props.dataList.map((item) => <Item {...item} key={item.name} />)
    : <div>Analyse Data not found</div>;
  return (
    <Panel header={props.name} className={"footer"}>
      <Timeline endless={true}>{dataListItems}</Timeline>
    </Panel>
  );
};

export default TimelineData;
