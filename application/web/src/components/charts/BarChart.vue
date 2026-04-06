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
  name: 'BarChart',
  props: {
    title: {
      type: String,
      default: '数据对比'
    },
    xData: {
      type: Array,
      default: () => ['协和医院', '301医院', '温江医疗中心', '监管中心']
    },
    yData: {
      type: Array,
      default: () => [120, 200, 150, 80]
    },
    color: {
      type: Array,
      default: () => ['#4facfe', '#00f2fe']
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
          trigger: 'axis',
          axisPointer: {
            type: 'shadow'
          }
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '3%',
          top: '15%',
          containLabel: true
        },
        xAxis: {
          type: 'category',
          data: this.xData,
          axisLine: {
            lineStyle: {
              color: '#E4E7ED'
            }
          },
          axisLabel: {
            color: '#606266',
            interval: 0,
            rotate: 30
          }
        },
        yAxis: {
          type: 'value',
          axisLine: {
            lineStyle: {
              color: '#E4E7ED'
            }
          },
          axisLabel: {
            color: '#606266'
          },
          splitLine: {
            lineStyle: {
              color: '#F2F6FC'
            }
          }
        },
        series: [
          {
            name: '数量',
            type: 'bar',
            data: this.yData,
            barWidth: '50%',
            itemStyle: {
              borderRadius: [8, 8, 0, 0],
              color: {
                type: 'linear',
                x: 0,
                y: 0,
                x2: 0,
                y2: 1,
                colorStops: [
                  { offset: 0, color: this.color[0] },
                  { offset: 1, color: this.color[1] }
                ]
              }
            },
            emphasis: {
              itemStyle: {
                color: {
                  type: 'linear',
                  x: 0,
                  y: 0,
                  x2: 0,
                  y2: 1,
                  colorStops: [
                    { offset: 0, color: '#667eea' },
                    { offset: 1, color: '#764ba2' }
                  ]
                }
              }
            }
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
