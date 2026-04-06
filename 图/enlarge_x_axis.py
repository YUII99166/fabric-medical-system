#!/usr/bin/env python3
# -*- coding: utf-8 -*-

with open('survey-results-visualization.html', 'r', encoding='utf-8', errors='ignore') as f:
    content = f.read()

# 查找并替换 x 轴的字体大小配置
# 将 x 轴的 ticks font size 从 14 改为 16
import re

# 匹配 x 轴配置中的 ticks font size
# 模式：x: { ... ticks: { font: { size: 14 }
pattern = r'(x:\s*\{[^}]*ticks:\s*\{\s*font:\s*\{\s*size:\s*)14(\s*\})'
content = re.sub(pattern, r'\g<1>16\g<2>', content)

with open('survey-results-visualization.html', 'w', encoding='utf-8') as f:
    f.write(content)

print("X轴字体已从14扩大到16")
