<script setup lang="tsx">
import { reactive } from 'vue';
import { NButton, NPopconfirm } from 'naive-ui';
import { fetchGetDelegateOrderList } from '@/service/api';
import { useAppStore } from '@/store/modules/app';
import { defaultTransform, useNaivePaginatedTable, useTableOperate } from '@/hooks/common/table';
import { $t } from '@/locales';
import UserOperateDrawer from './modules/delegate-bill-drawer.vue';
import UserSearch from './modules/delegate-search.vue';

const appStore = useAppStore();

const searchParams: Api.Delegate.DelegateOrderListSearchParams = reactive({
  page: 1,
  limit: 10,
  transaction_id: '',
  status: '',
  typo: '',
  from_base58: '',
  to_base58: '',
  start: 0,
  end: 0
});

const { columns, columnChecks, data, getData, getDataByPage, loading, mobilePagination } = useNaivePaginatedTable({
  api: () => fetchGetDelegateOrderList(searchParams),
  transform: response => defaultTransform(response),
  onPaginationParamsChange: params => {
    searchParams.page = params.page;
    searchParams.limit = params.pageSize;
  },
  columns: () => [
    {
      type: 'selection',
      align: 'center',
      width: 48
    },
    {
      key: 'id',
      title: $t('common.id'),
      align: 'center',
      width: 64
    },
    {
      key: 'transaction_id',
      title: $t('page.transaction.delegate.order.transaction_id'),
      align: 'center',
      minWidth: 100
    },
    {
      key: 'userGender',
      title: $t('page.transaction.delegate.order.status'),
      align: 'center',
      width: 100
    },
    {
      key: 'nickName',
      title: $t('page.transaction.delegate.order.currency'),
      align: 'center',
      minWidth: 100
    },
    {
      key: 'status',
      title: $t('page.transaction.delegate.order.from_base58'),
      align: 'center',
      width: 100
    },
    {
      key: 'status',
      title: $t('page.transaction.delegate.order.to_base58'),
      align: 'center',
      width: 100
    },
    {
      key: 'userPhone',
      title: $t('page.transaction.delegate.order.received_amount'),
      align: 'center',
      width: 120
    },
    {
      key: 'userEmail',
      title: $t('page.transaction.delegate.order.received_sun'),
      align: 'center',
      minWidth: 200
    },
    {
      key: 'userEmail',
      title: $t('page.transaction.delegate.order.withdraw_time'),
      align: 'center',
      minWidth: 200
    },
    {
      key: 'userEmail',
      title: $t('page.transaction.delegate.order.description'),
      align: 'center',
      minWidth: 200
    },
    {
      key: 'operate',
      title: $t('common.operate'),
      align: 'center',
      width: 130,
      render: row => (
        <div class="flex-center gap-8px">
          <NButton type="primary" ghost size="small" onClick={() => edit(row.id)}>
            {$t('common.edit')}
          </NButton>
          <NPopconfirm onPositiveClick={() => handleDelete(row.id)}>
            {{
              default: () => $t('common.confirmDelete'),
              trigger: () => (
                <NButton type="error" ghost size="small">
                  {$t('common.delete')}
                </NButton>
              )
            }}
          </NPopconfirm>
        </div>
      )
    }
  ]
});

const {
  drawerVisible,
  operateType,
  editingData,
  handleAdd,
  handleEdit,
  checkedRowKeys,
  onBatchDeleted,
  onDeleted
  // closeDrawer
} = useTableOperate(data, 'id', getData);

async function handleBatchDelete() {
  // request
  console.log(checkedRowKeys.value);

  onBatchDeleted();
}

function handleDelete(id: number) {
  // request
  console.log(id);

  onDeleted();
}

function edit(id: number) {
  handleEdit(id);
}
</script>

<template>
  <div class="min-h-500px flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <UserSearch v-model:model="searchParams" @search="getDataByPage" />
    <NCard
      :title="$t('page.transaction.delegate.order.transaction_id')"
      :bordered="false"
      size="small"
      class="card-wrapper sm:flex-1-hidden"
    >
      <template #header-extra>
        <TableHeaderOperation
          v-model:columns="columnChecks"
          :disabled-delete="checkedRowKeys.length === 0"
          :loading="loading"
          @add="handleAdd"
          @delete="handleBatchDelete"
          @refresh="getData"
        />
      </template>
      <NDataTable
        v-model:checked-row-keys="checkedRowKeys"
        :columns="columns"
        :data="data"
        size="small"
        :flex-height="!appStore.isMobile"
        :scroll-x="962"
        :loading="loading"
        remote
        :row-key="row => row.id"
        :pagination="mobilePagination"
        class="sm:h-full"
      />
      <UserOperateDrawer
        v-model:visible="drawerVisible"
        :operate-type="operateType"
        :row-data="editingData"
        @submitted="getDataByPage"
      />
    </NCard>
  </div>
</template>

<style scoped></style>
