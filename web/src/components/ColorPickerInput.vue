<template>
  <div class="input-group input-group-sm">
    <input type="text" class="form-control form-control-sm" aria-describedby="button-addon2" v-model="color" @input="updateFromInput" />
    <span class="input-group-addon color-picker-container">
      <input id="colorPicker" type="color" :value="color" @input="updateFromPicker" />
    </span>
  </div>
</template>
<script>
export default {
  name: "ColorPickerInput",
  props: {
    color: ''
  },
  data() {
    return {
      colors: {
        hex: '#000000',
      },
      colorValue: '',
      displayPicker: false,
    }
  },
  mounted() {
    this.setColor(this.color || '#000000');
  },
  methods: {
    setColor(color) {
      this.updateColors(color);
      this.colorValue = color;
    },
    updateColors(color) {
      this.colors = {
        hex: color
      }
      /*if(color.slice(0, 1) == '#') {
        this.colors = {
          hex: color
        };
      }
      else if(color.slice(0, 4) == 'rgba') {
        var rgba = color.replace(/^rgba?\(|\s+|\)$/g,'').split(','),
            hex = '#' + ((1 << 24) + (parseInt(rgba[0]) << 16) + (parseInt(rgba[1]) << 8) + parseInt(rgba[2])).toString(16).slice(1);
        this.colors = {
          hex: hex,
          a: rgba[3],
        }
      }*/
    },
    updateFromInput() {
      this.updateColors(this.colorValue);
    },
    updateFromPicker(color) {
      //console.log(color.target.value)
      this.setColor(color.target.value)
      /*if(color.rgba.a == 1) {
        this.colorValue = color.hex;
      }
      else {
        this.colorValue = 'rgba(' + color.rgba.r + ', ' + color.rgba.g + ', ' + color.rgba.b + ', ' + color.rgba.a + ')';
      }*/
    },

  },
  watch: {
    colorValue(val) {
      if(val) {
        this.updateColors(val);
        this.$emit('input', val);
      }
    }
  },
}
</script>

<style scoped>
.current-color {
  display: inline-block;
  width: 16px;
  height: 16px;
  background-color: #000;
  cursor: pointer;
}
.vc-chrome {
  position: absolute;
  top: 35px;
  right: 0;
  z-index: 9;
}
.form-control {
  border: 1px solid var(--border-color) !important;
  background-color: var(--background-color-primary) !important;
}
.form-control:focus {
  border: 1px solid var(--border-color) !important;
}
button:focus, select:focus, input:focus {
  border: none !important;
  outline: none !important;
  box-shadow: none !important;
}
#hexInput {
  height: 25px !important;
}
#colorPicker {
  height: 100%;
  padding: 0 !important;
  margin: 0 !important;
  border: none !important;
  background: none;
}
#colorPicker:hover {
  cursor: pointer !important;
}
input {
  background-color: transparent !important;
}
</style>