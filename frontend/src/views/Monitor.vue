<template>
  <div class="monitor">
    <el-row :gutter="20">
      <el-col :span="6">
        <div class="stat-card">
          <div class="stat-icon">
            <el-icon><Cpu /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ currentStats.cpu_usage?.toFixed(1) || '0' }}%</div>
            <div class="stat-label">CPU使用率</div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card blue">
          <div class="stat-icon">
            <el-icon><Coin /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ currentStats.memory_usage?.toFixed(1) || '0' }}%</div>
            <div class="stat-label">内存使用率</div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card green">
          <div class="stat-icon">
            <el-icon><Files /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ currentStats.disk_usage?.toFixed(1) || '0' }}%</div>
            <div class="stat-label">磁盘使用率</div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card orange">
          <div class="stat-icon">
            <el-icon><Connection /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ formatNetwork(currentStats.network_in || 0) }}</div>
            <div class="stat-label">网络流入</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-20">
      <el-col :span="12">
        <div class="chart-card">
          <h3 class="chart-title">CPU使用率趋势</h3>
          <div ref="cpuChartRef" class="chart"></div>
        </div>
      </el-col>
      <el-col :span="12">
        <div class="chart-card">
          <h3 class="chart-title">内存使用率趋势</h3>
          <div ref="memoryChartRef" class="chart"></div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-20">
      <el-col :span="12">
        <div class="chart-card">
          <h3 class="chart-title">磁盘使用率</h3>
          <div ref="diskChartRef" class="chart"></div>
        </div>
      </el-col>
      <el-col :span="12">
        <div class="chart-card">
          <h3 class="chart-title">网络流量趋势</h3>
          <div ref="networkChartRef" class="chart"></div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import * as echarts from 'echarts'
import { getSystemStats } from '../api/system'

const currentStats = ref({})
const cpuChartRef = ref(null)
const memoryChartRef = ref(null)
const diskChartRef = ref(null)
const networkChartRef = ref(null)

let cpuChart = null
let memoryChart = null
let diskChart = null
let networkChart = null

const maxDataPoints = 30
const cpuData = ref([])
const memoryData = ref([])
const networkInData = ref([])
const networkOutData = ref([])
const timeLabels = ref([])

const formatNetwork = (bytes) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB'
  if (bytes < 1024 * 1024 * 1024) return (bytes / 1024 / 1024).toFixed(2) + ' MB'
  return (bytes / 1024 / 1024 / 1024).toFixed(2) + ' GB'
}

const initCharts = () => {
  cpuChart = echarts.init(cpuChartRef.value)
  memoryChart = echarts.init(memoryChartRef.value)
  diskChart = echarts.init(diskChartRef.value)
  networkChart = echarts.init(networkChartRef.value)

  const lineOption = {
    tooltip: {
      trigger: 'axis'
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: []
    },
    yAxis: {
      type: 'value',
      max: 100,
      axisLabel: {
        formatter: '{value}%'
      }
    },
    series: [{
      type: 'line',
      smooth: true,
      data: [],
      areaStyle: {
        opacity: 0.3
      }
    }]
  }

  cpuChart.setOption({ ...lineOption, color: ['#667eea'] })
  memoryChart.setOption({ ...lineOption, color: ['#4facfe'] })

  diskChart.setOption({
    tooltip: {
      trigger: 'item'
    },
    series: [{
      type: 'gauge',
      startAngle: 180,
      endAngle: 0,
      min: 0,
      max: 100,
      splitNumber: 10,
      axisLine: {
        lineStyle: {
          width: 20,
          color: [
            [0.3, '#43e97b'],
            [0.7, '#4facfe'],
            [1, '#fa709a']
          ]
        }
      },
      pointer: {
        itemStyle: {
          color: 'auto'
        }
      },
      axisTick: {
        distance: -20,
        length: 8,
        lineStyle: {
          color: '#fff',
          width: 2
        }
      },
      splitLine: {
        distance: -25,
        length: 15,
        lineStyle: {
          color: '#fff',
          width: 3
        }
      },
      axisLabel: {
        color: '#999',
        distance: 30,
        fontSize: 12
      },
      detail: {
        valueAnimation: true,
        formatter: '{value}%',
        color: 'auto',
        fontSize: 24,
        offsetCenter: [0, '60%']
      },
      data: [{ value: 0 }]
    }]
  })

  networkChart.setOption({
    tooltip: {
      trigger: 'axis'
    },
    legend: {
      data: ['流入', '流出']
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: []
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        formatter: '{value} KB/s'
      }
    },
    series: [
      {
        name: '流入',
        type: 'line',
        smooth: true,
        data: [],
        areaStyle: { opacity: 0.3 }
      },
      {
        name: '流出',
        type: 'line',
        smooth: true,
        data: [],
        areaStyle: { opacity: 0.3 }
      }
    ]
  })
}

const fetchStats = async () => {
  try {
    const stats = await getSystemStats()
    currentStats.value = stats

    const now = new Date()
    const timeStr = now.toLocaleTimeString()

    timeLabels.value.push(timeStr)
    cpuData.value.push(stats.cpu_usage)
    memoryData.value.push(stats.memory_usage)
    networkInData.value.push((stats.network_in || 0) / 1024)
    networkOutData.value.push((stats.network_out || 0) / 1024)

    if (timeLabels.value.length > maxDataPoints) {
      timeLabels.value.shift()
      cpuData.value.shift()
      memoryData.value.shift()
      networkInData.value.shift()
      networkOutData.value.shift()
    }

    updateCharts()
  } catch (error) {
    console.error('Failed to fetch stats:', error)
  }
}

const updateCharts = () => {
  if (cpuChart) {
    cpuChart.setOption({
      xAxis: { data: timeLabels.value },
      series: [{ data: cpuData.value }]
    })
  }

  if (memoryChart) {
    memoryChart.setOption({
      xAxis: { data: timeLabels.value },
      series: [{ data: memoryData.value }]
    })
  }

  if (diskChart) {
    diskChart.setOption({
      series: [{ data: [{ value: currentStats.value.disk_usage || 0 }] }]
    })
  }

  if (networkChart) {
    networkChart.setOption({
      xAxis: { data: timeLabels.value },
      series: [
        { data: networkInData.value },
        { data: networkOutData.value }
      ]
    })
  }
}

const handleResize = () => {
  cpuChart?.resize()
  memoryChart?.resize()
  diskChart?.resize()
  networkChart?.resize()
}

let statsInterval = null

onMounted(async () => {
  await nextTick()
  initCharts()
  await fetchStats()
  statsInterval = setInterval(fetchStats, 3000)
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  if (statsInterval) {
    clearInterval(statsInterval)
  }
  window.removeEventListener('resize', handleResize)
  cpuChart?.dispose()
  memoryChart?.dispose()
  diskChart?.dispose()
  networkChart?.dispose()
})
</script>

<style scoped>
.monitor {
  padding: 0;
}

.stat-card {
  display: flex;
  align-items: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  border-radius: 8px;
  padding: 20px;
}

.stat-card.blue {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-card.green {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-card.orange {
  background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
}

.stat-icon {
  font-size: 48px;
  margin-right: 20px;
  opacity: 0.8;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  opacity: 0.9;
}

.mt-20 {
  margin-top: 20px;
}

.chart-card {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.chart-title {
  margin: 0 0 20px 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.chart {
  height: 280px;
}
</style>
