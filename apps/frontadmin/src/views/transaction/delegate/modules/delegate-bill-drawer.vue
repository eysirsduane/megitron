<script setup lang="ts">
import { ref } from 'vue';
import { fetchGetDelegateBill } from '@/service/api';
import { $t } from '@/locales';

defineOptions({
  name: 'DelegateBillDrawer'
});

const visible = defineModel<boolean>('visible', {
  default: false
});

type Model = Pick<
  Api.Delegate.DelegateBill,
  | 'order_id'
  | 'transaction_id'
  | 'currency'
  | 'from_base58'
  | 'to_base58'
  | 'delegated_amount'
  | 'description'
  | 'status'
>;

const model = ref(await getModel());

async function getModel(): Promise<Model> {
  // const data = ref({}) as Ref<Api.Delegate.DelegateBill>;

  const { data, error } = await fetchGetDelegateBill(1);
  if (!error) {
    return data;
  }

  return {
    order_id: 1,
    transaction_id: '',
    currency: '',
    from_base58: '',
    to_base58: '',
    delegated_amount: 1,
    description: '',
    status: ''
  };
}

function closeDrawer() {
  visible.value = false;
}
</script>

<template>
  <NDrawer v-model:show="visible" display-directive="show" :width="360">
    <NDrawerContent title="发货详情" :native-scrollbar="false" closable>
      <NForm ref="formRef" :model="model">
        <!--
 <NFormItem :label="$t('page.transaction.delegate.bill.order_id')" path="order_id">
          <NInput v-model:value="model.order_id" :placeholder="$t('page.transaction.delegate.bill.order_id')" />
        </NFormItem>
-->
        <NFormItem :label="$t('page.transaction.delegate.bill.transaction_id')" path="nickName">
          <NInput
            v-model:value="model.transaction_id"
            :placeholder="$t('page.transaction.delegate.bill.transaction_id')"
          />
        </NFormItem>
        <NFormItem :label="$t('page.transaction.delegate.bill.currency')" path="currency">
          <NInput v-model:value="model.currency" :placeholder="$t('page.transaction.delegate.bill.currency')" />
        </NFormItem>
        <NFormItem :label="$t('page.transaction.delegate.bill.from_base58')" path="from_base58">
          <NInput v-model:value="model.from_base58" :placeholder="$t('page.transaction.delegate.bill.from_base58')" />
        </NFormItem>
        <NFormItem :label="$t('page.transaction.delegate.bill.to_base58')" path="to_base58">
          <NInput v-model:value="model.from_base58" :placeholder="$t('page.transaction.delegate.bill.to_base58')" />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace :size="16">
          <NButton @click="closeDrawer">{{ $t('common.cancel') }}</NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>

<style scoped></style>
