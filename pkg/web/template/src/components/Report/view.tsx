import React, { FC } from 'react';
import { connect, ReportModelState, ConnectProps } from 'umi';

import { PanelGroup, Affix, Panel } from 'rsuite';
import 'rsuite/lib/PanelGroup/styles';
import 'rsuite/lib/Panel/styles';

import { NavHeader } from '@/components/NavHeader';
import { RenderData } from '@/components/RenderData';
import { TimelineData } from '@/components/TimelineData';

import './index.css';

interface ReportProps extends ConnectProps {
  reports: ReportModelState[];
}

const emptyState: ReportModelState = {
  id: '',
  name: '',
  status: 0,
  desc: '',
  time: '',
  type: 1,
  data: '',
  dataList: null,
};

const getStateById = (dataList: ReportModelState[], id: string) => {
  let val = dataList.find(s => s.id === id);
  return val ? val : emptyState;
};

const Report: FC<ReportProps> = (props) => {
  const envRenderData = getStateById(props.reports, '_ENV');
  const fileRenderData = getStateById(props.reports, '_ELF');
  const analyseRenderData = getStateById(props.reports, '_LOAD');
  return (
    <div className={'container'}>
      <Affix>
        <NavHeader />
      </Affix>
      <div className='content'>
        <PanelGroup>
          <RenderData {...envRenderData} />
          <RenderData {...fileRenderData} />
          <TimelineData {...analyseRenderData}/>
        </PanelGroup>
      </div>
    </div>
  );
};

export default connect(
  ({ reports }: { reports: ReportModelState[] }) => ({ reports }),
)(Report);
