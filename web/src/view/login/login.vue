<template>
  <div id="userLayout" class="user-layout-wrapper">
    <div class="container">
      <div class="top">
        <div class="desc">
          <img class="logo_login" src="@/assets/logo_login.png" alt="">
        </div>
        <div class="header">
          <a href="/"><span class="title">Bee Manager</span></a>
        </div>
      </div>
      <div class="main">
        <el-form
          ref="loginForm"
          :model="loginForm"
          :rules="rules"
          @keyup.enter.native="submitForm"
        >
          <el-form-item prop="username">
            <el-input v-model="loginForm.username" placeholder="请输入用户名">
              <i slot="suffix" class="el-input__icon el-icon-user" />
            </el-input>
          </el-form-item>
          <el-form-item prop="password">
            <el-input
              v-model="loginForm.password"
              :type="lock === 'lock' ? 'password' : 'text'"
              placeholder="请输入密码"
            >
              <i
                slot="suffix"
                :class="'el-input__icon el-icon-' + lock"
                @click="changeLock"
              />
            </el-input>
          </el-form-item>
          <el-form-item style="position: relative">
            <el-input
              v-model="loginForm.captcha"
              name="logVerify"
              placeholder="请输入验证码"
              style="width: 60%"
            />
            <div class="vPic">
              <img
                v-if="picPath"
                :src="picPath"
                width="100%"
                height="100%"
                alt="请输入验证码"
                @click="loginVerify()"
              >
            </div>
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              style="width: 100%"
              @click="submitForm"
            >登 录</el-button>
          </el-form-item>
        </el-form>
      </div>

      <div class="text-center text-white cursor-pointer" @click="goToRegister">注 册</div>
      <div class="footer">

      </div>
    </div>
  </div>
</template>

<script>
import { mapActions } from 'vuex'
import { captcha } from '@/api/user'
export default {
  name: 'Login',
  data() {
    const checkUsername = (rule, value, callback) => {
      if (value.length < 5) {
        return callback(new Error('请输入正确的用户名'))
      } else {
        callback()
      }
    }
    const checkPassword = (rule, value, callback) => {
      if (value.length < 6) {
        return callback(new Error('请输入正确的密码'))
      } else {
        callback()
      }
    }
    return {
      curYear: 0,
      lock: 'lock',
      loginForm: {
        username: '',
        password: '',
        captcha: '',
        captchaId: ''
      },
      rules: {
        username: [{ validator: checkUsername, trigger: 'blur' }],
        password: [{ validator: checkPassword, trigger: 'blur' }]
      },
      logVerify: '',
      picPath: ''
    }
  },
  created() {
    this.loginVerify()
    this.curYear = new Date().getFullYear()
  },
  methods: {
    ...mapActions('user', ['LoginIn']),
    async login() {
      return await this.LoginIn(this.loginForm)
    },
    goToRegister() {
      this.$router.push({
        name: 'Register'
      })
    },
    async submitForm() {
      this.$refs.loginForm.validate(async(v) => {
        if (v) {
          const flag = await this.login()
          if (!flag) {
            this.loginVerify()
          }
        } else {
          this.$message({
            type: 'error',
            message: '请正确填写登录信息',
            showClose: true
          })
          this.loginVerify()
          return false
        }
      })
    },
    goRegister() {
      this.$router.push({ name: 'Register' })
    },
    changeLock() {
      this.lock = this.lock === 'lock' ? 'unlock' : 'lock'
    },
    loginVerify() {
      captcha({}).then((ele) => {
        this.picPath = ele.data.picPath
        this.loginForm.captchaId = ele.data.captchaId
      })
    }
  }
}
</script>

<style scoped lang="scss">
@import "@/style/login.scss";
</style>
