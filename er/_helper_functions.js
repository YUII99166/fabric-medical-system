// ER图绘制辅助函数
// 这个文件包含所有ER图通用的绘制函数

// 计算矩形边缘点
function getRectEdgePoint(rectX, rectY, rectW, rectH, targetX, targetY) {
    const centerX = rectX + rectW / 2;
    const centerY = rectY + rectH / 2;
    const dx = targetX - centerX;
    const dy = targetY - centerY;
    
    if (dx === 0 && dy === 0) return { x: centerX, y: centerY };
    
    const angle = Math.atan2(dy, dx);
    const absAngle = Math.abs(angle);
    
    // 判断交点在哪条边
    if (absAngle < Math.atan2(rectH, rectW)) {
        // 右边
        return { x: rectX + rectW, y: centerY + (rectW / 2) * Math.tan(angle) };
    } else if (absAngle > Math.PI - Math.atan2(rectH, rectW)) {
        // 左边
        return { x: rectX, y: centerY - (rectW / 2) * Math.tan(angle) };
    } else if (angle > 0) {
        // 下边
        return { x: centerX + (rectH / 2) / Math.tan(angle), y: rectY + rectH };
    } else {
        // 上边
        return { x: centerX - (rectH / 2) / Math.tan(angle), y: rectY };
    }
}

// 计算椭圆边缘点
function getEllipseEdgePoint(cx, cy, rx, ry, targetX, targetY) {
    const dx = targetX - cx;
    const dy = targetY - cy;
    
    if (dx === 0 && dy === 0) return { x: cx, y: cy };
    
    const angle = Math.atan2(dy, dx);
    return {
        x: cx + rx * Math.cos(angle),
        y: cy + ry * Math.sin(angle)
    };
}

// 绘制实体到属性的连接线
function drawEntityAttributeLines(g, entityX, entityY, entityW, entityH, attributes) {
    attributes.forEach(attr => {
        const rectEdge = getRectEdgePoint(entityX, entityY, entityW, entityH, attr.cx, attr.cy);
        const ellipseEdge = getEllipseEdgePoint(attr.cx, attr.cy, attr.rx, attr.ry, rectEdge.x, rectEdge.y);
        
        g.strokeStyle = '#000';
        g.lineWidth = 2;
        g.beginPath();
        g.moveTo(rectEdge.x, rectEdge.y);
        g.lineTo(ellipseEdge.x, ellipseEdge.y);
        g.stroke();
    });
}
