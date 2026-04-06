#!/usr/bin/env python3
# -*- coding: utf-8 -*-

with open('survey-results-visualization.html', 'rb') as f:
    content = f.read()

# 替换所有剩余的 font size 13 为 15
content = content.replace(b'font: { size: 13 }', b'font: { size: 15 }')
content = content.replace(b'font: { size: 12 }', b'font: { size: 14 }')

with open('survey-results-visualization.html', 'wb') as f:
    f.write(content)

print("问卷调查文件字体已优化完成")
