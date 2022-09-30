<template>
  <div class="dropdown" v-if="options">

    <!-- Dropdown Input -->
    <input class="dropdown-input"
           :name="name"
           @focus="showOptions()"
           @blur="exit()"
           @keyup="keyMonitor"
           v-model="searchFilter"
           :disabled="disabled"
           :placeholder="placeholder" />

    <!-- Dropdown Menu -->
    <div class="dropdown-content"
         v-show="optionsShown">
      <div
          class="dropdown-item"
          @mousedown="selectOption(option)"
          v-for="(option, index) in filteredOptions"
          :key="index">
        {{ option.id }}
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Dropdown',
  template: 'Dropdown',
  props: {
    name: {
      type: String,
      required: false,
      default: 'dropdown',
      note: 'Input name'
    },
    options: {
      type: Array,
      required: true,
      default: [],
      note: 'Options of dropdown. An array of options with id and name',
    },
    placeholder: {
      type: String,
      required: false,
      default: 'Please select an option',
      note: 'Placeholder of dropdown'
    },
    disabled: {
      type: Boolean,
      required: false,
      default: false,
      note: 'Disable the dropdown'
    },
    maxItem: {
      type: Number,
      required: false,
      default: 50,
      note: 'Max items showing'
    }
  },
  data() {
    return {
      selected: {},
      optionsShown: false,
      searchFilter: ''
    }
  },
  created() {
    //this.$emit('selected', this.selected);
  },
  computed: {
    filteredOptions() {
      const filtered = [];
      const regOption = new RegExp(this.searchFilter, 'id');
      for (const option of this.options) {
        if (this.searchFilter.length < 1 || option.id.match(regOption)){
          if (filtered.length < this.maxItem) filtered.push(option);
        }
      }
      return filtered.sort((a, b) => {
        if (a.id < b.id)
          return -1;
        if (a.id > b.id)
          return 1;
        return 0;
      });

    },

  },
  methods: {
    selectOption(option) {
      if (!this.disabled) {
        this.selected = option;
        this.optionsShown = false;
        this.searchFilter = this.selected.id;
        this.$emit('selected', this.selected);
      }
    },
    showOptions(){
      if (!this.disabled) {
        this.searchFilter = '';
        this.optionsShown = true;
      }
    },
    exit() {
      if (!this.disabled) {
        if (!this.selected.id) {
          this.selected = {};
          this.searchFilter = '';
        } else {
          this.searchFilter = this.selected.id;
        }
        this.$emit('selected', this.selected);
        this.optionsShown = false;
      }
    },
    // Selecting when pressing Enter
    keyMonitor: function(event) {
      if (event.key === "Enter" && this.filteredOptions[0] && !this.disabled) {
        this.selectOption(this.filteredOptions[0]);
      }

    }
  },
  watch: {
    searchFilter() {
      if (!this.disabled) {
        if (this.filteredOptions.length === 0) {
          this.selected = {};
        } else {
          this.selected = this.filteredOptions[0];
        }
        this.$emit('filter', this.searchFilter);
      }
    }
  }
};
</script>


<style scoped>
input {
  all: unset;
}
.dropdown {
  position: relative;
  display: block;
  margin: auto;
}
.dropdown .dropdown-input {
  background: var(--background-color-primary);
  cursor: text !important;
  border: 1px solid var(--border-color);
  border-radius: 3px;
  color: #333;
  display: block;
  padding: 6px;
  min-width: 50px;
  max-width: 250px;
}
.dropdown .dropdown-input:hover {
  background: var(--background-color-secondary) !important;
  border: 1px solid var(--border-color-hover) !important;
}
.dropdown .dropdown-input:focus {
  background: var(--background-color-secondary) !important;
  border: 1px solid var(--border-color-hover) !important;
}
.dropdown .dropdown-input:focus-visible {
  background: var(--background-color-secondary) !important;
  border: 1px solid var(--border-color-hover) !important;
}
:focus-visible {
  background: var(--background-color-secondary) !important;
  border: 1px solid var(--border-color-hover) !important;
}

.dropdown .dropdown-content {
  position: absolute;
  background-color: var(--background-color-primary);
  width: 100% !important;
  /*min-width: 48px;
  max-width: 248px;*/
  max-height: 248px;
  border: 1px solid var(--border-color);
  box-shadow: 0px -8px 34px 0px rgba(0, 0, 0, 0.05);
  overflow: auto !important;
  z-index: 100;
}
.dropdown .dropdown-content .dropdown-item {
  color: var(--text-primary-color);
  padding: 8px;
  text-decoration: none;
  display: block;
  cursor: pointer;
}
.dropdown .dropdown-content .dropdown-item:hover {
   background-color: var(--background-color-secondary);
}
.dropdown:hover .dropdown-content {
  display: block;
}

/* dropdown scrollbar width */
::-webkit-scrollbar {
  width: 1px;
}

/* dropdown scrollbar track */
::-webkit-scrollbar-track {
  background: var(--background-color-secondary);
  box-shadow: inset 0 0 5px grey;
  border-radius: 10px;
}

/* dropdown scrollbar handle */
::-webkit-scrollbar-thumb {
  border-color: var(--border-color);
  border-radius: 10px;
  background: var(--background-color-primary);
}

/* dropdown scrollbar handle on hover */
::-webkit-scrollbar-thumb:hover {
  background: #555;
}

</style>