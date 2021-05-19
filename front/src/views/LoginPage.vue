<template>
  <div>
<!--    <h2></h2>-->
    <div class="margin-bottom-20">
      <form @submit.prevent>
        <div>
          <label for="username"><span>Username</span> <input id="username" type="text" v-model="username" placeholder="Username"/></label>
        </div>
        <div>
          <label for="password"><span>Password</span> <input id="password" type="password" v-model="password" placeholder="Password"/></label>
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
    onLoginButtonClick() {
      let promise;
      if (this.loginSelected) {
        promise = this.login({username: this.username, password: this.password});
      } else {
        promise = this.signup({username: this.username, password: this.password});
      }

      promise.then(() => [
        this.$router.push({name: "home"})
      ]).catch(error => {
        this.error = error;
      });
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
.margin-bottom-20{
  margin-bottom: 20px;
}
</style>