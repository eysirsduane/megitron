import { request } from '../request';

/**
 */
export function fetchGetDelegateOrderList(params?: Api.Delegate.DelegateOrderListSearchParams) {
  return request<Api.Delegate.DelegateOrderList>({
    url: '/delegate/order/list',
    method: 'get',
    params
  });
}

export function fetchGetDelegateBill(oid: number) {
  return request<Api.Delegate.DelegateBill>({
    url: '/delegate/bill',
    method: 'get',
    data: {
      oid
    }
  });
}
