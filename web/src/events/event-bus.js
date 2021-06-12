import Vue from 'vue';

export const EventBus = new Vue();

export function publishEvent(success, message) {
  EventBus.$emit("events", {type: "events", success: success, message: message});
}