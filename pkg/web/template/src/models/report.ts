import { Effect, Reducer, Subscription, request } from 'umi';

export interface ReportModelState {
  id: string;
  name: string,
  status: number,
  desc: string,
  time: string,
  type: number,
  data: string | null,
  dataList: ReportModelState[] | null,
}

export const ErrorStatus = 1;

export interface ReportModelType {
  namespace: string;
  state: ReportModelState[];
  effects: {
    fetch: Effect;
  };
  reducers: {
    save: Reducer<ReportModelState[]>;
  };
  subscriptions: { setup: Subscription }
}

const initState: ReportModelState = {
  id: '',
  name: 'init...',
  status: 0,
  desc: 'init...',
  time: 'init...',
  type: 1,
  data: 'init...',
  dataList: null,
};

const ReportModel: ReportModelType = {
  namespace: 'reports',

  state: [initState, initState, initState],

  effects: {
    * fetch({ type, payload }, { call, put, select }) {
      const data = yield request('/api/report');
      yield put({
        type: 'save',
        payload: [...data],
      });
    },
  },

  reducers: {
    save(state, action) {
      return [
        ...action.payload,
      ];
    },
  },

  subscriptions: {
    setup({ dispatch, history }) {
      return history.listen(({ pathname }) => {
        if (pathname === '/') {
          dispatch({ type: 'fetch' });
        }
      });
    },
  },

};

export default ReportModel;
