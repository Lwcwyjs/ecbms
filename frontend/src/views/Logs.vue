<template>
  <div class="logs">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>系统日志</span>
          <el-button type="primary" size="small" @click="fetchLogs">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>
      <el-table :data="logs" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getTypeTagType(row.type)" size="small">
              {{ getTypeLabel(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="level" label="级别" width="100">
          <template #default="{ row }">
            <el-tag :type="getLevelTagType(row.level)" size="small">
              {{ getLevelLabel(row.level) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="message" label="消息" />
        <el-table-column prop="created_at" label="时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getSystemLogs } from '../api/system'

const logs = ref([])

const fetchLogs = async () => {
  try {
    logs.value = await getSystemLogs()
  } catch (error) {
    console.error('Failed to fetch logs:', error)
  }
}

const getTypeLabel = (type) => {
  const labels = {
    auth: '认证',
    network: '网络',
    system: '系统',
    user: '用户'
  }
  return labels[type] || type
}

const getTypeTagType = (type) => {
  const types = {
    auth: 'primary',
    network: 'success',
    system: 'warning',
    user: 'info'
  }
  return types[type] || ''
}

const getLevelLabel = (level) => {
  const labels = {
    info: '信息',
    warning: '警告',
    error: '错误'
  }
  return labels[level] || level
}

const getLevelTagType = (level) => {
  const types = {
    info: 'info',
    warning: 'warning',
    error: 'danger'
  }
  return types[level] || ''
}

const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString()
}

onMounted(() => {
  fetchLogs()
})
</script>

<style scoped>
.logs {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
