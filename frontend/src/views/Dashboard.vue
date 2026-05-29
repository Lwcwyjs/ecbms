<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <div class="stat-card">
          <div class="stat-value">{{ systemInfo.hostname || '-' }}</div>
          <div class="stat-label">主机名</div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card blue">
          <div class="stat-value">{{ systemStats.cpu_usage?.toFixed(1) || '0' }}%</div>
          <div class="stat-label">CPU使用率</div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card green">
          <div class="stat-value">{{ systemStats.memory_usage?.toFixed(1) || '0' }}%</div>
          <div class="stat-label">内存使用率</div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card orange">
          <div class="stat-value">{{ systemStats.disk_usage?.toFixed(1) || '0' }}%</div>
          <div class="stat-label">磁盘使用率</div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-20">
      <el-col :span="12">
        <div class="card">
          <h3 class="card-title">系统信息</h3>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="操作系统">{{ systemInfo.os }} / {{ systemInfo.platform }}</el-descriptions-item>
            <el-descriptions-item label="系统版本">{{ systemInfo.platform_version }}</el-descriptions-item>
            <el-descriptions-item label="内核版本">{{ systemInfo.kernel_version }}</el-descriptions-item>
            <el-descriptions-item label="CPU型号">{{ systemInfo.cpu_model }}</el-descriptions-item>
            <el-descriptions-item label="CPU核心数">{{ systemInfo.cpu_cores }}</el-descriptions-item>
            <el-descriptions-item label="总内存">{{ formatBytes(systemInfo.total_memory) }}</el-descriptions-item>
            <el-descriptions-item label="运行时间">{{ formatUptime(systemInfo.uptime) }}</el-descriptions-item>
          </el-descriptions>
        </div>
      </el-col>
      <el-col :span="12">
        <div class="card">
          <h3 class="card-title">快捷操作</h3>
          <el-space direction="vertical" style="width: 100%;">
            <el-button type="primary" style="width: 100%;" @click="refreshData">
              <el-icon><Refresh /></el-icon>
              刷新数据
            </el-button>
            <el-button type="warning" style="width: 100%;" @click="handleReboot">
              <el-icon><RefreshRight /></el-icon>
              重启系统
            </el-button>
            <el-button type="danger" style="width: 100%;" @click="handleShutdown">
              <el-icon><SwitchButton /></el-icon>
              关闭系统
            </el-button>
          </el-space>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import { getSystemInfo, getSystemStats, rebootSystem, shutdownSystem } from '../api/system'

const systemInfo = ref({})
const systemStats = ref({})

const fetchSystemInfo = async () => {
  try {
    systemInfo.value = await getSystemInfo()
  } catch (error) {
    console.error('Failed to fetch system info:', error)
  }
}

const fetchSystemStats = async () => {
  try {
    systemStats.value = await getSystemStats()
  } catch (error) {
    console.error('Failed to fetch system stats:', error)
  }
}

const refreshData = () => {
  fetchSystemInfo()
  fetchSystemStats()
}

const handleReboot = () => {
  ElMessageBox.confirm('确定要重启系统吗？', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    await rebootSystem()
    ElMessage.success('系统重启命令已发送')
  }).catch(() => {})
}

const handleShutdown = () => {
  ElMessageBox.confirm('确定要关闭系统吗？', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    await shutdownSystem()
    ElMessage.success('系统关机命令已发送')
  }).catch(() => {})
}

const formatBytes = (bytes) => {
  if (!bytes) return '-'
  const gb = bytes / (1024 * 1024 * 1024)
  return `${gb.toFixed(2)} GB`
}

const formatUptime = (seconds) => {
  if (!seconds) return '-'
  const days = Math.floor(seconds / 86400)
  const hours = Math.floor((seconds % 86400) / 3600)
  const mins = Math.floor((seconds % 3600) / 60)
  return `${days}天 ${hours}小时 ${mins}分钟`
}

let statsInterval = null

onMounted(() => {
  fetchSystemInfo()
  fetchSystemStats()
  statsInterval = setInterval(fetchSystemStats, 5000)
})

onUnmounted(() => {
  if (statsInterval) {
    clearInterval(statsInterval)
  }
})
</script>

<style scoped>
.dashboard {
  padding: 0;
}

.card {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.card-title {
  margin: 0 0 20px 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.mt-20 {
  margin-top: 20px;
}
</style>
