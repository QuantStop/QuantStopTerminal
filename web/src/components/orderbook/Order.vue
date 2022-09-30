<template>
  <li class="ask" :key="ask[0]">
    <div class="d-flex justify-content-between">
      <span class="amount">{{ ask[0] }}</span>
      <span class="size">{{ ask[1] }}</span>
    </div>
  </li>
</template>

<script>
export default {
  name: "Order"
}
</script>

<style scoped>

</style>