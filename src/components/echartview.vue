<template>
  <div class='hello'>
    <v-chart v-for="(item, index) in chartDataList" :key="index" :options="item"/>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  data () {
    return {
      options: {
        title: {
          x: 'center',
          text: 'chart'
        },
        xAxis: {
          type: 'category',
          axisLabel: {
            interval: 0,
            rotate: 40
          }
        },
        yAxis: {
          type: 'value'
        },
        series: [{
          data: [],
          type: 'line',
          smooth: true
        }]
      },
      chartDataList: []
    }
  },
  mounted () {
    this.getData()
  },
  methods: {
    getData () {
      let param = {}
      param.mobile = ''
      let url = '/api/TestChart'
      axios.post(url, param).then((responseData) => {
        responseData.data.forEach(item => {
          let temp = JSON.parse(JSON.stringify(this.options))
          temp.xAxis.data = item.X
          temp.series[0].data = item.Data
          temp.title.text = item.Name
          this.chartDataList.push(temp)
        })
      }).catch((error) => {
        console.log(error)
      })
    }
  }

}
</script>
