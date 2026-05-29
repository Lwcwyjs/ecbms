<template>
  <div class="settings">
    <el-card>
      <template #header>
        <span>个人设置</span>
      </template>
      <el-row :gutter="40">
        <el-col :span="12">
          <h3 class="section-title">修改密码</h3>
          <el-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef" label-width="100px">
            <el-form-item label="当前密码" prop="old_password">
              <el-input v-model="passwordForm.old_password" type="password" show-password />
            </el-form-item>
            <el-form-item label="新密码" prop="new_password">
              <el-input v-model="passwordForm.new_password" type="password" show-password />
            </el-form-item>
            <el-form-item label="确认密码" prop="confirm_password">
              <el-input v-model="passwordForm.confirm_password" type="password" show-password />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleChangePassword" :loading="passwordLoading">
                确认修改
              </el-button>
            </el-form-item>
          </el-form>
        </el-col>
        <el-col :span="12">
          <h3 class="section-title">账户信息</h3>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="用户名">{{ userStore.user?.username }}</el-descriptions-item>
            <el-descriptions-item label="角色">
              <el-tag :type="userStore.user?.role === 'admin' ? 'danger' : 'primary'">
                {{ userStore.user?.role === 'admin' ? '管理员' : '普通用户' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="用户ID">{{ userStore.user?.id }}</el-descriptions-item>
          </el-descriptions>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '../store'
import { changePassword } from '../api/auth'

const userStore = useUserStore()
const passwordLoading = ref(false)
const passwordFormRef = ref(null)

const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== passwordForm.new_password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const passwordRules = {
  old_password: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const handleChangePassword = async () => {
  try {
    await passwordFormRef.value.validate()
    passwordLoading.value = true
    await changePassword({
      old_password: passwordForm.old_password,
      new_password: passwordForm.new_password
    })
    ElMessage.success('密码修改成功')
    passwordForm.old_password = ''
    passwordForm.new_password = ''
    passwordForm.confirm_password = ''
  } catch (error) {
    console.error('Failed to change password:', error)
  } finally {
    passwordLoading.value = false
  }
}
</script>

<style scoped>
.settings {
  padding: 0;
}

.section-title {
  margin: 0 0 20px 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}
</style>
