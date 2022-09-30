<template>
  <div class="led-box" :style="style">
    <div class="led-light" :style="style" :class="connected ? 'led-green' : 'led-red'"></div>
  </div>
</template>

<script>
export default {
  name: 'StatusIndicator',
  props: {
    value:{
      type: Boolean,
      required: true,
      default: false
    },
    width: {
      type: Number,
      required: true,
      default: 100
    },
    height: {
      type: Number,
      required: true,
      default: 50
    }
  },
  computed: {
    connected() {
      return this.value
    },
    style() {
      return {
        width: `${this.height}px`,
        height: `${this.height}px`
      }
    },
  }
};
</script>

<style scoped>
.led-box {
  display: inline-block!important;
}
.led-light {
  display: inline-block!important;
  border-radius: 50%;
}
.led-red {
  background-color: red;
  box-shadow: rgba(0, 0, 0, 0.2) 0 -1px 5px 1px, inset #441313 0 -1px 9px, darkred 0 2px 8px;
}
.led-green {
  background-color: green;
  box-shadow: rgba(0, 0, 0, 0.2) 0 -1px 5px 1px, inset #304701 0 -1px 9px, darkgreen 0 2px 8px;
}
</style>