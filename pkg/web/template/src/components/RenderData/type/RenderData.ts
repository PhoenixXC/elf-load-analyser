// RenderData define
interface RenderDataType {
  name: string,
  status: number,
  desc: string,
  time: string,
  type: number,
  data: string | null,
  dataList: RenderDataType[] | null,
}

const ErrorStatus = 1;

export { RenderDataType, ErrorStatus };
