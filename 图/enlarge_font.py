#!/usr/bin/env python3
# -*- coding: utf-8 -*-

with open('graphviz.svg', 'rb') as f:
    content = f.read()

content = content.replace(b'font-size="10.00"', b'font-size="14.00"')
content = content.replace(b'font-size="11.00"', b'font-size="15.00"')
content = content.replace(b'font-size="12.00"', b'font-size="16.00"')

with open('graphviz.svg', 'wb') as f:
    f.write(content)

print("字体已放大：10→14, 11→15, 12→16")
