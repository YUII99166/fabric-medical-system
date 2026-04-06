#!/usr/bin/env python3
# -*- coding: utf-8 -*-

with open('survey-results-visualization.html', 'r', encoding='utf-8', errors='ignore') as f:
    content = f.read()

# 再扩大5个字号
# Chart.js 默认字体: 19 -> 24
content = content.replace('Chart.defaults.font.size = 19;', 'Chart.defaults.font.size = 24;')

# 图例字体: 18 -> 23
content = content.replace('font: { size: 18 }', 'font: { size: 23 }')

# X轴字体: 19 -> 24
content = content.replace('font: { size: 19 }', 'font: { size: 24 }')

# Y轴字体: 17 -> 22
content = content.replace('font: { size: 17 }', 'font: { size: 22 }')

# 确保Chart.js使用宋体（已经设置过，但再次确认）
# Chart.defaults.font.family 已经是 SimSun

with open('survey-results-visualization.html', 'w', encoding='utf-8') as f:
    f.write(content)

print("所有字体已再次扩大5个字号，并使用宋体")
print("- Chart默认: 19→24")
print("- 图例: 18→23")
print("- X轴: 19→24")
print("- Y轴: 17→22")
print("- 字体: 全部使用宋体(SimSun)")
