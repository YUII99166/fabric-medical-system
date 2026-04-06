#!/usr/bin/env python3
# -*- coding: utf-8 -*-

with open('survey-results-visualization.html', 'r', encoding='utf-8', errors='ignore') as f:
    content = f.read()

# 优化坐标轴字体大小
content = content.replace("ticks: { stepSize:", "ticks: { font: { size: 14 }, stepSize:")
content = content.replace("ticks: {", "ticks: { font: { size: 14 },")

with open('survey-results-visualization.html', 'w', encoding='utf-8') as f:
    f.write(content)

print("图表坐标轴字体已优化")
