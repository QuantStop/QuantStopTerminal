<template>
  <Layout title="TechnicalIndicatorKLineChart">
    <div id="technical-indicator-k-line" class="k-line-chart"/>
    <div
      class="k-line-chart-menu-container">
      <span style="padding-right: 10px">Main Chart Indicators</span>
      <button
        v-for="type in mainTechnicalIndicatorTypes"
        :key="type"
        v-on:click="setCandleTechnicalIndicator(type)">
        {{type}}
      </button>
      <button
        v-on:click="setCandleTechnicalIndicator('EMOJI')">
        Custom
      </button>
      <span style="padding-right: 10px;padding-left: 12px">Submap Indicator</span>
      <button
        v-for="type in subTechnicalIndicatorTypes"
        :key="type"
        v-on:click="setSubTechnicalIndicator(type)">
        {{type}}
      </button>
      <button
        v-on:click="setSubTechnicalIndicator('EMOJI')">
        Custom
      </button>
    </div>
  </Layout>
</template>

<script>
  import {dispose, init} from 'klinecharts'
  import generatedKLineDataList from './generatedKLineDataList'
  import Layout from "./Layout"

  const fruits = [
    '🍏', '🍎', '🍐', '🍊', '🍋', '🍌',
    '🍉', '🍇', '🍓', '🍈', '🍒', '🍑',
    '🍍', '🥥', '🥝', '🥭', '🥑', '🍏'
  ]

  // custom indicator
  const emojiTechnicalIndicator = {
    name: 'EMOJI',
    plots: [
      { key: 'emoji' }
    ],
    calcTechnicalIndicator: (kLineDataList) => {
      const result = []
      kLineDataList.forEach(kLineData => {
        result.push({ emoji: kLineData.close, text: fruits[Math.floor(Math.random() * 17)] })
      })
      return result
    },
    render: ({ ctx, dataSource, viewport, xAxis, yAxis }) => {
      ctx.font = `${viewport.barSpace}px Helvetica Neue`
      ctx.textAlign = 'center'
      for (let i = dataSource.from; i < dataSource.to; i++) {
        const data = dataSource.technicalIndicatorDataList[i]
        const x = xAxis.convertToPixel(i)
        const y = yAxis.convertToPixel(data.emoji)
        ctx.fillText(data.text, x, y)
      }
    }
  }

  export default {
    name: 'TechnicalIndicatorKLineChart',
    components: {Layout},
    data: function () {
      return {
        mainTechnicalIndicatorTypes: ['MA', 'EMA', 'SAR'],
        subTechnicalIndicatorTypes: ['VOL', 'MACD', 'KDJ']
      }
    },
    mounted: function () {
      this.kLineChart = init('technical-indicator-k-line')
      this.kLineChart.addTechnicalIndicatorTemplate(emojiTechnicalIndicator)
      this.paneId = this.kLineChart.createTechnicalIndicator('VOL', false)
      this.kLineChart.applyNewData(generatedKLineDataList())
    },
    methods: {
      setCandleTechnicalIndicator: function (type) {
        this.kLineChart.createTechnicalIndicator(type, false, { id: 'candle_pane' })
      },
      setSubTechnicalIndicator: function (type) {
        this.kLineChart.createTechnicalIndicator(type, false, { id: this.paneId })
      },
    },
    destroyed: function () {
      dispose('technical-indicator-k-line')
    }
  }
</script>
