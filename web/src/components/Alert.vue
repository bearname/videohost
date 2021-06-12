<template>
  <div id="alert" class="alert-popup hide"></div>
</template>

<script>
import {EventBus} from "@/events/event-bus";

export default {
  name: "Alert"
}

const clickHandler = function (events) {
  let alertElement = document.getElementById("alert");

  alertElement.innerText = events.message;
  alertElement.classList.remove('hide');
  setTimeout(() => {
    alertElement.classList.add('hide');
  }, 2500);
}

EventBus.$on('events', clickHandler);
</script>

<style scoped>
.alert-popup {
  outline: none;
  position: fixed;
  box-sizing: border-box;
  left: 10px;
  top: 10px;
  z-index: 100;
  max-width: 288px;
  max-height: 48px;
  width: 288px;
  height: 48px;
  color: #E5E5E5;

  background-color: #323232;
  box-shadow: 0 2px 5px 0 rgb(0 0 0);
  border-radius: 2px;

  transition: 0.2s;
  line-height: 37px;
  font-size: 17px;
  padding: 6px 12px;
  transform: translateY(0);
}

.hide {
  opacity: 0;
  /*display: none;*/
  transform: translateY(-100px);
}

.success {
  color: #0eff70;
}

.failed {

}
</style>