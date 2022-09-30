<template>
  <div
    class="toggle"
    :class=[this.state_class]
    @click.self="onClick"
    :style="styleBtn"
  >
    <div
      class="draggable"
      @mousedown.prevent="dragStart"
      :style="style"
    >
    </div>
  </div>
</template>

<script>
export default {
  name: "ToggleSwitch",
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
  data() {
    return {
      state: false,
      pressed: 0,
      position: 0
    }
  },
  mounted() {
    this.toggle(this.value)
  },
  computed: {
    styleBtn() {
      return {
        width: `${this.width}px`,
        height: `${this.height}px`
      }
    },
    style() {
      return {
        transform: `translateX(${this.pos_percentage})`,
        width: `${this.height -1}px`,
        height: `${this.height -2}px`
      }
    },
    pos_percentage() {
      return `${this.position / this.width * this.width}%`
    },
    state_class() {
      if (this.state) {
        return 'active'
      }
    }
  },
  watch: {
    position() {
      this.state = this.position >= 50
    }
  },
  methods: {
    onClick() {
      this.toggle(!this.state)
      this.emit()
    },
    toggle(state) {
      this.state = state
      this.position = !state
          ? 0
          : 100
    },
    dragging(e) {
      const pos = e.clientX - this.$el.offsetLeft
      const percent = pos / this.width * this.width
      this.position = percent <= 0
          ? 0
          : percent >= 100
              ? 100
              : percent
    },
    dragStart(e) {
      this.startTimer()
      window.addEventListener('mousemove', this.dragging)
      window.addEventListener('mouseup', this.dragStop)
    },
    dragStop() {
      window.removeEventListener('mousemove', this.dragging)
      window.removeEventListener('mouseup', this.dragStop)
      this.resolvePosition()
      clearInterval(this.$options.interval)
      if (this.pressed < 30) {
        this.toggle(!this.state)
      }
      this.pressed = 0
      this.emit()
    },
    startTimer() {
      this.$options.interval = setInterval(() => {
        this.pressed++
      }, 1)
    },
    resolvePosition() {
      this.position = this.state ? 100 : 0
    },
    emit() {
      this.$emit('input', this.state)
    }
  }
}
</script>

<style scoped>

.toggle {
  cursor: pointer;
  background: var(--background-color-secondary);
  border: 1px solid var(--border-color);
  border-radius: 200px;
  padding: 0;
  transition: background 0.6s;
}
.draggable {
  position: absolute;

  background: var(--background-color-primary);
  border-radius: 100%;
  /*box-shadow: 0px 3px 10px rgba(0,0,0, 0.6);*/
  transform: translateX(0%);
  transition: transform 0.05s ease-in-out;
}

.active {
   /*background: var(--background-color-primary);*/
   transition: background 0.6s;
}



</style>