#!/usr/bin/env python3
# -*- coding: utf-8 -*-

with open('survey-results-visualization.html', 'r', encoding='utf-8', errors='ignore') as f:
    content = f.read()

# 扩大所有字体3个字号
# Chart.js 默认字体: 16 -> 19
content = content.replace('Chart.defaults.font.size = 16;', 'Chart.defaults.font.size = 19;')

# 图例字体: 15 -> 18
content = content.replace('font: { size: 15 }', 'font: { size: 18 }')

# X轴字体: 16 -> 19
content = content.replace('font: { size: 16 }', 'font: { size: 19 }')

# Y轴字体: 14 -> 17
content = content.replace('font: { size: 14 }', 'font: { size: 17 }')

with open('survey-results-visualization.html', 'w', encoding='utf-8') as f:
    f.write(content)

print("所有字体已扩大3个字号")
print("- Chart默认: 16→19")
print("- 图例: 15→18")
print("- X轴: 16→19")
print("- Y轴: 14→17")
