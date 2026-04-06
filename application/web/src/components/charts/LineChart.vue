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
  name: 'LineChart',
  props: {
    title: {
      type: String,
      default: '数据趋势'
    },
    xData: {
      type: Array,
      default: () => ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
    },
    yData: {
      type: Array,
      default: () => [120, 200, 150, 80, 70, 110, 130]
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
            type: 'cross'
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
          boundaryGap: false,
          data: this.xData,
          axisLine: {
            lineStyle: {
              color: '#E4E7ED'
            }
          },
          axisLabel: {
            color: '#606266'
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
            type: 'line',
            smooth: true,
            data: this.yData,
            lineStyle: {
              width: 3,
              color: {
                type: 'linear',
                x: 0,
                y: 0,
                x2: 1,
                y2: 0,
                colorStops: [
                  { offset: 0, color: this.color[0] },
                  { offset: 1, color: this.color[1] }
                ]
              }
            },
            areaStyle: {
              color: {
                type: 'linear',
                x: 0,
                y: 0,
                x2: 0,
                y2: 1,
                colorStops: [
                  { offset: 0, color: this.color[0] + '4D' },
                  { offset: 1, color: this.color[1] + '1A' }
                ]
              }
            },
            itemStyle: {
              color: this.color[0],
              borderColor: '#fff',
              borderWidth: 2
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
