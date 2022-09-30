<template>
<div id="orderbook">

    <div class="orderbook-header-container row p-0 m-0 flex-row justify-content-center">
      <div class="orderbook-header d-flex align-items-center">
        <h6 class="m-0">Orderbook</h6>
      </div>
    </div>

    <div class="mid-market-container">
      <span class="mid-market-price">
        {{ midpoint }}
      </span>
      <span class="mid-market-price-text">
        Mid Market Price
      </span>
    </div>

    <div class="bid-ask-header-container row p-0 m-0 flex-row">
      <div class="col-sm-6 col-sm-push-6 p-0">
        <div class="bid-header d-flex justify-content-between">
          <h6 class="l2-size-header">Size</h6>
          <h6 class="l2-bid-header">Bid</h6>
        </div>
      </div>
      <div class="col-sm-6 col-sm-pull-6 p-0">
        <div class="ask-header d-flex justify-content-between">
          <h6 class="l2-ask-header">Ask</h6>
          <h6 class="l2-size-header">Size</h6>
        </div>
      </div>
    </div>

    <div class="orderbook-row-container row p-0 m-0 flex-row">
      <div class="col-sm-6 col-sm-push-6 p-0">
        <ul id="bids" ref="bids">
          <template v-for="bid in bids" :key="bid[0]">
            <li class="bid d-flex justify-content-between">
              <span class="size">{{ bid[1] }}</span>
              <span class="amount">{{ bid[0] }}</span>
            </li>
          </template>

        </ul>
      </div>
      <div class="col-sm-6 col-sm-pull-6 p-0">
        <ul id="asks" ref="asks">
          <template v-for="ask in asks" :key="ask[0]">
            <li class="ask d-flex justify-content-between">
              <span class="amount">{{ ask[0] }}</span>
              <span class="size">{{ ask[1] }}</span>
            </li>
          </template>
        </ul>
      </div>
    </div>

</div>
</template>

<script>
export default {
  name: "Orderbook",
  props: {
    asks: {
      type: Array,
      default: [],
    },
    bids: {
      type: Array,
      default: [],
    },
    midpoint: "",
    spread: "",
    bookInitialized: false,
  },
}
</script>

<style scoped>
#orderbook {
  margin: 0;
  padding: 0;
  height: calc(80%);
  max-width: 400px;
  width: calc(100%);
  display: flex;
  flex-direction: column;
  flex: 1;
  grid-area: order-book / order-book / order-book / order-book;
  grid-row-start: order-book;
  grid-column-start: order-book;
  grid-row-end: order-book;
  grid-column-end: order-book;

  border-left: 1px solid var(--border-color) !important;
  border-bottom: 1px solid var(--border-color) !important;
}
.orderbook-header {
  color: var(--text-primary-color) !important;
  text-align: center !important;
  height: 37px;
  width: auto;
  margin: 0;
}
.orderbook-header-container {
  border-bottom: 1px solid var(--border-color) !important;
}
.orderbook-row-container {
  overflow-y: auto;
  background-color: var(--orderbook-background-color) !important;
}
#bids {
  color: green !important;
  margin: 0;
  padding: 0;
  position: relative;
  overflow: hidden;
}
#asks {
  color: red !important;
  margin: 0;
  padding: 0 2px 0 2px;
  position: relative;
  overflow: hidden;
}

.mid-market-container {
  display: flex;
  flex-direction: column;
  flex-wrap: nowrap;
  justify-content: center;
  align-items: center;
  min-height: 75px;
}
.mid-market-price {
  font-size: 18px;
  color: var(--text-primary-color);
}
.mid-market-price-text {
  font-size: 12px;
  color: gray;
}

.bid-ask-header-container {
  border-bottom: 1px solid var(--border-color) !important;
  border-top: 1px solid var(--border-color) !important;
}
.ask-header {
  padding: 0 2px 0 2px;
}
.bid-header {
  padding: 0 2px 0 2px;
}
.ask {
  color: red !important;
  font-size: 12px !important;
  background-color: var(--orderbook-background-color) !important;
}
.bid {
  color: green !important;
  font-size: 12px !important;
  background-color: var(--orderbook-background-color) !important;
}
.size {
  color: gray;
  font-size: 12px !important;
}
.ask:hover, .bid:hover {
  background-color: var(--background-color-secondary) !important;
  cursor: pointer !important;
}
.l2-size-header {
  color: gray;
  margin: 0;
}

.l2-ask-header {
  color: red;
  margin: 0;
}
.l2-bid-header {
  color: green;
  margin: 0;
}

.red {
  color: red;
}
.green {
  color: green;
}

/* width */
::-webkit-scrollbar {
  width: 0px;
}

/* Track */
::-webkit-scrollbar-track {
  background: var(--background-color-secondary);
  box-shadow: inset 0 0 5px grey;
  border-radius: 10px;
}

/* Handle */
::-webkit-scrollbar-thumb {
  border-color: var(--border-color);
  border-radius: 10px;
  background: var(--background-color-primary);
}

/* Handle on hover */
::-webkit-scrollbar-thumb:hover {
  background: #555;
}
</style>