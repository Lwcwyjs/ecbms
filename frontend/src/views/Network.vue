<template>
  <div class="network">
    <el-card class="mb-20">
      <template #header>
        <div class="card-header">
          <span>网络接口列表</span>
          <el-button type="primary" size="small" @click="fetchInterfaces">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>
      <el-table :data="interfaces" border>
        <el-table-column prop="name" label="接口名称" width="150" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'up' ? 'success' : 'danger'">
              {{ row.status === 'up' ? '已连接' : '已断开' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="hwaddr" label="MAC地址" width="200" />
        <el-table-column prop="ip_address" label="IP地址" />
        <el-table-column prop="netmask" label="子网掩码" />
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="openConfigDialog(row)">
              配置
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card>
      <template #header>
        <div class="card-header">
          <span>网络测试</span>
        </div>
      </template>
      <el-form :inline="true" :model="pingForm">
        <el-form-item label="目标地址">
          <el-input v-model="pingForm.host" placeholder="请输入IP或域名" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handlePing" :loading="pingLoading">
            <el-icon><Position /></el-icon>
            Ping测试
          </el-button>
        </el-form-item>
      </el-form>
      <el-input
        v-if="pingResult"
        type="textarea"
        :rows="8"
        :model-value="pingResult"
        readonly
      />
    </el-card>

    <el-dialog
      v-model="configDialogVisible"
      title="网络配置"
      width="500px"
      @close="resetConfigForm"
    >
      <el-form :model="configForm" label-width="100px">
        <el-form-item label="接口名称">
          <el-input :model-value="configForm.interface_name" readonly />
        </el-form-item>
        <el-form-item label="配置模式">
          <el-radio-group v-model="configForm.mode">
            <el-radio value="dhcp">DHCP自动获取</el-radio>
            <el-radio value="static">静态IP</el-radio>
          </el-radio-group>
        </el-form-item>
        <template v-if="configForm.mode === 'static'">
          <el-form-item label="IP地址">
            <el-input v-model="configForm.ip_address" placeholder="例如: 192.168.1.100" />
          </el-form-item>
          <el-form-item label="子网掩码">
            <el-input v-model="configForm.netmask" placeholder="例如: 255.255.255.0" />
          </el-form-item>
          <el-form-item label="网关">
            <el-input v-model="configForm.gateway" placeholder="例如: 192.168.1.1" />
          </el-form-item>
          <el-form-item label="DNS服务器">
            <el-input v-model="configForm.dns_servers" placeholder="例如: 8.8.8.8,8.8.4.4" />
          </el-form-item>
        </template>
      </el-form>
      <template #footer>
        <el-button @click="configDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveConfig" :loading="configLoading">
          保存配置
        </el-button>
        <el-button type="success" @click="applyConfig" :loading="applyLoading">
          应用配置
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  getNetworkInterfaces,
  getNetworkConfigs,
  configureNetwork,
  applyNetworkConfig,
  pingTest
} from '../api/network'

const interfaces = ref([])
const configDialogVisible = ref(false)
const configLoading = ref(false)
const applyLoading = ref(false)
const pingLoading = ref(false)
const pingResult = ref('')

const configForm = reactive({
  interface_name: '',
  mode: 'dhcp',
  ip_address: '',
  netmask: '',
  gateway: '',
  dns_servers: ''
})

const pingForm = reactive({
  host: 'www.baidu.com'
})

const fetchInterfaces = async () => {
  try {
    interfaces.value = await getNetworkInterfaces()
  } catch (error) {
    console.error('Failed to fetch interfaces:', error)
  }
}

const openConfigDialog = (iface) => {
  configForm.interface_name = iface.name
  configForm.mode = 'dhcp'
  configForm.ip_address = ''
  configForm.netmask = ''
  configForm.gateway = ''
  configForm.dns_servers = ''
  configDialogVisible.value = true
}

const resetConfigForm = () => {
  configForm.interface_name = ''
  configForm.mode = 'dhcp'
  configForm.ip_address = ''
  configForm.netmask = ''
  configForm.gateway = ''
  configForm.dns_servers = ''
}

const saveConfig = async () => {
  try {
    configLoading.value = true
    await configureNetwork(configForm)
    ElMessage.success('配置保存成功')
    configDialogVisible.value = false
  } catch (error) {
    console.error('Failed to save config:', error)
  } finally {
    configLoading.value = false
  }
}

const applyConfig = async () => {
  try {
    applyLoading.value = true
    await applyNetworkConfig({ interface_name: configForm.interface_name })
    ElMessage.success('配置已应用')
    configDialogVisible.value = false
    fetchInterfaces()
  } catch (error) {
    console.error('Failed to apply config:', error)
  } finally {
    applyLoading.value = false
  }
}

const handlePing = async () => {
  if (!pingForm.host) {
    ElMessage.warning('请输入目标地址')
    return
  }
  try {
    pingLoading.value = true
    pingResult.value = '正在测试...'
    const result = await pingTest({ host: pingForm.host })
    pingResult.value = result.output || JSON.stringify(result, null, 2)
  } catch (error) {
    pingResult.value = error.response?.data?.output || error.message
  } finally {
    pingLoading.value = false
  }
}

onMounted(() => {
  fetchInterfaces()
})
</script>

<style scoped>
.network {
  padding: 0;
}

.mb-20 {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
