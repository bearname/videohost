<template>
  <div>
    <div class="margin-bottom-20">
      <form @submit.prevent>
        <div>
          <label for="username"><span>Username</span> <input id="username" type="text" v-model="username"
                                                             placeholder="Username" required/></label>
        </div>
        <div v-if="!loginSelected">
          <label for="username"><span>Email</span> <input id="email" type="email" v-model="email" placeholder="Email"
                                                          required/></label>
        </div>
        <div v-if="!loginSelected">
          <label for="username"><input id="isSubscribed" type="checkbox" v-model="isSubscribed" placeholder="Subscribe"
                                       required/><span>Subscribe to video processing result</span></label>
        </div>
        <div>
          <label for="password"><span>Password</span> <input id="password" type="password" v-model="password"
                                                             placeholder="Password" required/></label>
        </div>
        <v-btn @click="onLoginButtonClick">{{ buttonString }}</v-btn>
      </form>
    </div>
    <p v-if="error">{{ error }}</p>
    <v-btn @click.prevent="toggleLogin">{{ textLoginString }}</v-btn>
  </div>
</template>

<script>
import {mapActions} from "vuex";

export default {
  name: "LoginPage",
  data() {
    return {
      loginSelected: true,
      username: "",
      password: "",
      email: "",
      isSubscribed: false,
      error: ""
    }
  },
  methods: {
    ...mapActions({
      login: "auth/login",
      signup: "auth/signup"
    }),
    toggleLogin() {
      this.loginSelected = !this.loginSelected;
    },
    async onLoginButtonClick() {
      try {
        if (this.loginSelected) {
          await this.login({username: this.username, password: this.password});
        } else {
          await this.signup({
            username: this.username,
            password: this.password,
            email: this.email,
            isSubscribed: this.isSubscribed
          });
        }

        await this.$router.push({name: "home"})
      } catch (error) {
        this.error = error;
      }
    }
  },
  computed: {
    buttonString() {
      if (this.loginSelected) {
        return "LOGIN"
      } else {
        return "SIGNUP"
      }
    },
    textLoginString() {
      if (this.loginSelected) {
        return "Signup instead";
      } else {
        return "Login instead";
      }
    }
  }
}
</script>

<style scoped>
.margin-bottom-20 {
  margin-bottom: 20px;
}
</style>