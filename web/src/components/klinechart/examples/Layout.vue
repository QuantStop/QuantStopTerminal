
<template>
  <div class="k-line-chart-container w-100">
    <h3 class="k-line-chart-title">{{title}}</h3>
    <slot></slot>
  </div>
</template>

<script>
export default {
  name: 'Layout',
  props: {
    title: String
  },
}
</script>
