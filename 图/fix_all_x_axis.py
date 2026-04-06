#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re

with open('survey-results-visualization.html', 'r', encoding='utf-8', errors='ignore') as f:
    content = f.read()

# 为所有没有 x 轴配置的 scales 添加 x 轴字体配置
# 查找所有 scales: { 后面直接跟 y: 的情况，在前面插入 x 轴配置

# 模式1: scales: { y: { ... } }  -> scales: { x: { ticks: { font: { size: 16 } } }, y: { ... } }
pattern1 = r'(scales:\s*\{\s*)(y:\s*\{)'
replacement1 = r'\1x: { ticks: { font: { size: 16 } } }, \2'
content = re.sub(pattern1, replacement1, content)

# 模式2: scales: { x: { ... } } 但没有 ticks 配置，添加 ticks
pattern2 = r'(x:\s*\{\s*)(beginAtZero)'
replacement2 = r'\1ticks: { font: { size: 16 } }, \2'
content = re.sub(pattern2, replacement2, content)

with open('survey-results-visualization.html', 'w', encoding='utf-8') as f:
    f.write(content)

print("所有图表的X轴字体已统一设置为16")
