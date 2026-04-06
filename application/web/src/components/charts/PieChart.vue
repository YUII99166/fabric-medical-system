<template>
  <div class="chart-container">
    <v-chart 
      :option="chartOption" 
      :autoresize="true"
      style="width: 100%; height: 100%;"
    />
  </div>
</template>

<script>
export default {
  name: 'PieChart',
  props: {
    title: {
      type: String,
      default: '数据分布'
    },
    data: {
      type: Array,
      default: () => [
        { value: 335, name: '医生' },
        { value: 310, name: '病人' },
        { value: 234, name: '药店' }
      ]
    }
  },
  computed: {
    chartOption() {
      return {
        title: {
          text: this.title,
          left: 'center',
          textStyle: {
            fontSize: 16,
            fontWeight: 600,
            color: '#303133'
          }
        },
        tooltip: {
          trigger: 'item',
          formatter: '{a} <br/>{b}: {c} ({d}%)'
        },
        legend: {
          orient: 'vertical',
          left: 'left',
          top: 'middle',
          textStyle: {
            color: '#606266'
          }
        },
        series: [
          {
            name: '角色分布',
            type: 'pie',
            radius: ['40%', '70%'],
            center: ['60%', '50%'],
            avoidLabelOverlap: false,
            itemStyle: {
              borderRadius: 10,
              borderColor: '#fff',
              borderWidth: 2
            },
            label: {
              show: false,
              position: 'center'
            },
            emphasis: {
              label: {
                show: true,
                fontSize: 20,
                fontWeight: 'bold'
              }
            },
            labelLine: {
              show: false
            },
            data: this.data,
            color: [
              { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [
                { offset: 0, color: '#667eea' },
                { offset: 1, color: '#764ba2' }
              ]},
              { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [
                { offset: 0, color: '#4facfe' },
                { offset: 1, color: '#00f2fe' }
              ]},
              { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [
                { offset: 0, color: '#fa709a' },
                { offset: 1, color: '#fee140' }
              ]},
              { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [
                { offset: 0, color: '#43e97b' },
                { offset: 1, color: '#38f9d7' }
              ]}
            ]
          }
        ]
      }
    }
  }
}
</script>

<style scoped>
.chart-container {
  width: 100%;
  height: 100%;
  min-height: 300px;
}
</style>
