const path = require("path");
const PptxGenJS = require("pptxgenjs");

const pptx = new PptxGenJS();
pptx.layout = "LAYOUT_WIDE";
pptx.author = "OpenAI Codex";
pptx.company = "OpenAI";
pptx.subject = "Academic defense presentation";
pptx.title = "基于区块链的社区医疗档案管理信息系统";
pptx.lang = "zh-CN";
pptx.theme = {
  headFontFace: "Microsoft YaHei",
  bodyFontFace: "Microsoft YaHei",
  lang: "zh-CN",
};

const C = {
  navy: "17324D",
  blue: "2F6B9A",
  blue2: "4E8DBA",
  blue3: "7AA7C7",
  pale: "EEF4F8",
  pale2: "F7FAFC",
  text: "243746",
  muted: "5F7285",
  line: "D9E4EC",
  white: "FFFFFF",
  green: "4E8F7A",
  red: "C65D57",
  gold: "D8A94B",
};

const root = path.resolve(__dirname, "..");
const assets = {
  logo: path.join(root, "application", "web", "public", "image", "logo.png"),
  admin: path.join(root, "图", "管理员端功能结构图-黑白版 (1).png"),
};

function addSlideBase(slide, title, index, total) {
  slide.background = { color: C.pale2 };
  slide.addShape(pptx.ShapeType.rect, {
    x: 0,
    y: 0,
    w: 10,
    h: 0.38,
    fill: { color: C.navy },
    line: { color: C.navy, transparency: 100 },
  });
  slide.addText(title, {
    x: 0.55,
    y: 0.52,
    w: 6.8,
    h: 0.45,
    fontFace: "Microsoft YaHei",
    fontSize: 24,
    bold: true,
    color: C.navy,
    margin: 0,
  });
  slide.addShape(pptx.ShapeType.line, {
    x: 0.55,
    y: 1.06,
    w: 8.9,
    h: 0,
    line: { color: C.line, width: 1.1 },
  });
  slide.addText(`${index}/${total}`, {
    x: 9.1,
    y: 0.54,
    w: 0.45,
    h: 0.28,
    fontFace: "Microsoft YaHei",
    fontSize: 11,
    color: C.muted,
    align: "right",
    margin: 0,
  });
}

function addFooter(slide) {
  slide.addText("基于论文与项目文档自动整理", {
    x: 0.55,
    y: 5.18,
    w: 2.5,
    h: 0.2,
    fontFace: "Microsoft YaHei",
    fontSize: 10,
    color: "7F8C97",
    margin: 0,
  });
}

function addBulletList(slide, items, opts = {}) {
  const runs = [];
  items.forEach((item, idx) => {
    runs.push({
      text: item,
      options: { bullet: true, breakLine: idx !== items.length - 1 },
    });
  });
  slide.addText(runs, {
    x: opts.x,
    y: opts.y,
    w: opts.w,
    h: opts.h,
    fontFace: "Microsoft YaHei",
    fontSize: opts.fontSize || 20,
    color: opts.color || C.text,
    breakLine: false,
    paraSpaceAfterPt: opts.paraSpaceAfterPt || 8,
    valign: "top",
    margin: 0,
  });
}

function addInfoCard(slide, cfg) {
  slide.addShape(pptx.ShapeType.roundRect, {
    x: cfg.x,
    y: cfg.y,
    w: cfg.w,
    h: cfg.h,
    rectRadius: 0.08,
    fill: { color: cfg.fill || C.white },
    line: { color: cfg.line || C.line, width: 1 },
  });
  if (cfg.bandColor) {
    slide.addShape(pptx.ShapeType.rect, {
      x: cfg.x,
      y: cfg.y,
      w: 0.08,
      h: cfg.h,
      fill: { color: cfg.bandColor },
      line: { color: cfg.bandColor, transparency: 100 },
    });
  }
  slide.addText(cfg.title, {
    x: cfg.x + 0.16,
    y: cfg.y + 0.12,
    w: cfg.w - 0.24,
    h: 0.22,
    fontFace: "Microsoft YaHei",
    fontSize: cfg.titleSize || 15,
    bold: true,
    color: cfg.titleColor || C.navy,
    margin: 0,
  });
  if (cfg.body) {
    slide.addText(cfg.body, {
      x: cfg.x + 0.16,
      y: cfg.y + 0.42,
      w: cfg.w - 0.26,
      h: cfg.h - 0.5,
      fontFace: "Microsoft YaHei",
      fontSize: cfg.bodySize || 12.5,
      color: cfg.bodyColor || C.text,
      valign: "top",
      margin: 0,
    });
  }
}

function addMetricCard(slide, cfg) {
  slide.addShape(pptx.ShapeType.roundRect, {
    x: cfg.x,
    y: cfg.y,
    w: cfg.w,
    h: cfg.h,
    rectRadius: 0.08,
    fill: { color: cfg.fill || C.white },
    line: { color: cfg.line || C.line, width: 1 },
  });
  slide.addText(cfg.value, {
    x: cfg.x + 0.1,
    y: cfg.y + 0.16,
    w: cfg.w - 0.2,
    h: 0.36,
    fontFace: "Microsoft YaHei",
    fontSize: cfg.valueSize || 24,
    bold: true,
    color: cfg.valueColor || C.navy,
    align: "center",
    margin: 0,
  });
  slide.addText(cfg.label, {
    x: cfg.x + 0.1,
    y: cfg.y + cfg.h - 0.34,
    w: cfg.w - 0.2,
    h: 0.18,
    fontFace: "Microsoft YaHei",
    fontSize: cfg.labelSize || 11.5,
    color: cfg.labelColor || C.muted,
    align: "center",
    margin: 0,
  });
}

function addStep(slide, cfg) {
  slide.addShape(pptx.ShapeType.roundRect, {
    x: cfg.x,
    y: cfg.y,
    w: cfg.w,
    h: cfg.h,
    rectRadius: 0.08,
    fill: { color: cfg.fill || C.white },
    line: { color: cfg.line || C.line, width: 1 },
  });
  slide.addShape(pptx.ShapeType.ellipse, {
    x: cfg.x + 0.12,
    y: cfg.y + 0.12,
    w: 0.38,
    h: 0.38,
    fill: { color: cfg.dotColor || C.blue },
    line: { color: cfg.dotColor || C.blue, transparency: 100 },
  });
  slide.addText(String(cfg.num), {
    x: cfg.x + 0.12,
    y: cfg.y + 0.17,
    w: 0.38,
    h: 0.15,
    fontFace: "Microsoft YaHei",
    fontSize: 12,
    bold: true,
    color: C.white,
    align: "center",
    margin: 0,
  });
  slide.addText(cfg.title, {
    x: cfg.x + 0.6,
    y: cfg.y + 0.11,
    w: cfg.w - 0.72,
    h: 0.2,
    fontFace: "Microsoft YaHei",
    fontSize: 13.5,
    bold: true,
    color: C.navy,
    margin: 0,
  });
  slide.addText(cfg.body, {
    x: cfg.x + 0.6,
    y: cfg.y + 0.38,
    w: cfg.w - 0.72,
    h: cfg.h - 0.45,
    fontFace: "Microsoft YaHei",
    fontSize: 11.5,
    color: C.text,
    margin: 0,
  });
}

function addArrow(slide, x, y, w) {
  slide.addShape(pptx.ShapeType.chevron, {
    x,
    y,
    w,
    h: 0.28,
    fill: { color: "B9CAD7" },
    line: { color: "B9CAD7", transparency: 100 },
  });
}

const totalSlides = 15;

// 1 Cover
{
  const slide = pptx.addSlide();
  slide.background = { color: C.pale2 };
  slide.addShape(pptx.ShapeType.rect, {
    x: 0,
    y: 0,
    w: 10,
    h: 5.625,
    fill: { color: C.pale2 },
    line: { color: C.pale2, transparency: 100 },
  });
  slide.addShape(pptx.ShapeType.rect, {
    x: 0,
    y: 0,
    w: 3.2,
    h: 5.625,
    fill: { color: C.navy },
    line: { color: C.navy, transparency: 100 },
  });
  slide.addText("学术汇报", {
    x: 0.48,
    y: 0.82,
    w: 1.5,
    h: 0.36,
    fontFace: "Microsoft YaHei",
    fontSize: 22,
    bold: true,
    color: C.white,
    margin: 0,
  });
  slide.addText("毕业答辩 / 组会 / 会议", {
    x: 0.48,
    y: 1.22,
    w: 2.1,
    h: 0.24,
    fontFace: "Microsoft YaHei",
    fontSize: 12.5,
    color: "D4E1EA",
    margin: 0,
  });
  slide.addText("基于区块链的社区医疗档案管理信息系统", {
    x: 3.62,
    y: 1.05,
    w: 5.55,
    h: 0.9,
    fontFace: "Microsoft YaHei",
    fontSize: 26,
    bold: true,
    color: C.navy,
    margin: 0,
  });
  slide.addText("基于 Hyperledger Fabric 的医疗数据共享与隐私保护方案", {
    x: 3.62,
    y: 2.04,
    w: 5.35,
    h: 0.3,
    fontFace: "Microsoft YaHei",
    fontSize: 15,
    color: C.blue,
    italic: true,
    margin: 0,
  });
  slide.addShape(pptx.ShapeType.line, {
    x: 3.62,
    y: 2.42,
    w: 5.2,
    h: 0,
    line: { color: C.line, width: 1.1 },
  });
  addInfoCard(slide, {
    x: 3.62,
    y: 2.78,
    w: 1.85,
    h: 0.9,
    title: "研究主题",
    body: "区块链医疗数据管理\n跨机构共享与授权访问",
    bandColor: C.blue,
  });
  addInfoCard(slide, {
    x: 5.63,
    y: 2.78,
    w: 1.85,
    h: 0.9,
    title: "技术路线",
    body: "Fabric 联盟链\nGo + Gin + Vue.js",
    bandColor: C.green,
  });
  addInfoCard(slide, {
    x: 7.64,
    y: 2.78,
    w: 1.65,
    h: 0.9,
    title: "汇报重点",
    body: "创新点\n实验验证\n结论展望",
    bandColor: C.gold,
  });
  slide.addText("答辩人：XXX\n指导教师：XXX\n专业：计算机科学与技术", {
    x: 3.68,
    y: 4.08,
    w: 3.4,
    h: 0.78,
    fontFace: "Microsoft YaHei",
    fontSize: 14,
    color: C.text,
    margin: 0,
  });
  slide.addText("请基于实际信息替换封面占位内容", {
    x: 3.68,
    y: 4.96,
    w: 3.3,
    h: 0.18,
    fontFace: "Microsoft YaHei",
    fontSize: 10.5,
    color: C.muted,
    italic: true,
    margin: 0,
  });
  slide.addImage({
    path: assets.logo,
    x: 7.62,
    y: 4.08,
    w: 1.2,
    h: 1.2,
    sizing: { type: "contain", w: 1.2, h: 1.2 },
  });
}

// 2 Agenda
{
  const slide = pptx.addSlide();
  addSlideBase(slide, "目录", 2, totalSlides);
  const items = [
    ["01", "研究背景", "行业痛点与研究动机"],
    ["02", "问题与方法", "目标、路线与系统设计"],
    ["03", "创新实现", "联盟链架构与关键机制"],
    ["04", "实验验证", "问卷、性能与安全结果"],
    ["05", "结论展望", "成果总结与后续方向"],
  ];
  items.forEach((item, i) => {
    const y = 1.28 + i * 0.74;
    slide.addShape(pptx.ShapeType.roundRect, {
      x: 0.88,
      y,
      w: 8.2,
      h: 0.56,
      rectRadius: 0.07,
      fill: { color: i % 2 === 0 ? C.white : C.pale },
      line: { color: C.line, width: 1 },
    });
    slide.addShape(pptx.ShapeType.roundRect, {
      x: 1.06,
      y: y + 0.11,
      w: 0.7,
      h: 0.34,
      rectRadius: 0.05,
      fill: { color: C.navy },
      line: { color: C.navy, transparency: 100 },
    });
    slide.addText(item[0], {
      x: 1.06,
      y: y + 0.17,
      w: 0.7,
      h: 0.12,
      fontFace: "Microsoft YaHei",
      fontSize: 11,
      bold: true,
      color: C.white,
      align: "center",
      margin: 0,
    });
    slide.addText(item[1], {
      x: 1.98,
      y: y + 0.13,
      w: 1.5,
      h: 0.18,
      fontFace: "Microsoft YaHei",
      fontSize: 18,
      bold: true,
      color: C.navy,
      margin: 0,
    });
    slide.addText(item[2], {
      x: 4.1,
      y: y + 0.16,
      w: 3.8,
      h: 0.16,
      fontFace: "Microsoft YaHei",
      fontSize: 12.5,
      color: C.muted,
      margin: 0,
    });
  });
  addFooter(slide);
}

// 3 Background
{
  const slide = pptx.addSlide();
  addSlideBase(slide, "研究背景：社区医疗数据共享痛点突出", 3, totalSlides);
  addBulletList(slide, [
    "社区医疗信息化加速推进",
    "跨院病历共享仍然困难",
    "中心化存储存在隐私风险",
    "病历真伪与追溯能力不足",
    "分级诊疗协同效率受限",
  ], { x: 0.72, y: 1.32, w: 3.5, h: 2.3, fontSize: 18.5 });
  slide.addText("问卷痛点证据", {
    x: 4.82,
    y: 1.28,
    w: 2.0,
    h: 0.18,
    fontFace: "Microsoft YaHei",
    fontSize: 14,
    bold: true,
    color: C.navy,
    margin: 0,
  });
  slide.addChart(pptx.ChartType.bar, [
    {
      name: "占比",
      labels: ["系统互通困难", "重复检查", "跨院查询受阻", "隐私泄露担忧"],
      values: [100, 81.6, 78.1, 68.4],
    },
  ], {
    x: 4.72,
    y: 1.58,
    w: 4.1,
    h: 2.65,
    catAxisLabelFontFace: "Microsoft YaHei",
    valAxisLabelFontFace: "Microsoft YaHei",
    showLegend: false,
    showValue: true,
    valAxisMinVal: 0,
    valAxisMaxVal: 100,
    valAxisMajorUnit: 20,
    chartColors: [C.blue],
    chartArea: { fill: { color: C.white }, border: { color: C.line, pt: 1 } },
    catAxisLabelColor: C.muted,
    valAxisLabelColor: C.muted,
    valGridLine: { color: C.line, pt: 0.6 },
    catGridLine: { color: C.white, transparency: 100 },
    dataLabelColor: C.text,
    dataLabelPosition: "outEnd",
  });
  addMetricCard(slide, { x: 0.86, y: 4.15, w: 1.55, h: 0.72, value: "92", label: "有效问卷", fill: C.white });
  addMetricCard(slide, { x: 2.55, y: 4.15, w: 1.55, h: 0.72, value: "4类", label: "核心用户群", fill: C.white });
  addInfoCard(slide, {
    x: 4.72,
    y: 4.34,
    w: 4.1,
    h: 0.62,
    title: "结论",
    body: "数据孤岛、重复检查与隐私风险构成社区医疗信息化的核心阻碍。",
    bandColor: C.red,
  });
  addFooter(slide);
}

// 4 Objective
{
  const slide = pptx.addSlide();
  addSlideBase(slide, "问题提出：构建可信共享与患者主导的病历体系", 4, totalSlides);
  addInfoCard(slide, {
    x: 0.72,
    y: 1.35,
    w: 3.05,
    h: 1.06,
    title: "研究问题",
    body: "如何在多机构协同场景下，兼顾病历共享效率、数据安全性与患者授权可控性？",
    bandColor: C.blue,
  });
  addInfoCard(slide, {
    x: 0.72,
    y: 2.62,
    w: 3.05,
    h: 1.2,
    title: "研究目标",
    body: "设计一套基于 Fabric 的社区医疗档案管理系统，实现跨院访问、审计追踪与双存储优化。",
    bandColor: C.green,
  });
  slide.addShape(pptx.ShapeType.ellipse, {
    x: 4.45,
    y: 1.78,
    w: 1.35,
    h: 1.35,
    fill: { color: C.navy },
    line: { color: C.navy, transparency: 100 },
  });
  slide.addText("核心目标", {
    x: 4.45,
    y: 2.28,
    w: 1.35,
    h: 0.16,
    fontFace: "Microsoft YaHei",
    fontSize: 14,
    bold: true,
    color: C.white,
    align: "center",
    margin: 0,
  });
  const spokes = [
    { x: 6.35, y: 1.08, text: "跨院共享\n避免重复检查", color: C.blue },
    { x: 6.75, y: 2.34, text: "患者授权\n保护隐私边界", color: C.green },
    { x: 5.96, y: 3.66, text: "链上存证\n实现可信追溯", color: C.gold },
  ];
  spokes.forEach((s) => {
    slide.addShape(pptx.ShapeType.line, {
      x: 5.12,
      y: 2.45,
      w: s.x - 5.12,
      h: s.y - 2.45,
      line: { color: C.line, width: 1.3 },
    });
    addInfoCard(slide, {
      x: s.x,
      y: s.y,
      w: 2.15,
      h: 0.92,
      title: s.text.split("\n")[0],
      body: s.text.split("\n")[1],
      bandColor: s.color,
    });
  });
  addInfoCard(slide, {
    x: 0.72,
    y: 4.18,
    w: 8.3,
    h: 0.64,
    title: "汇报导向",
    body: "论文内容重组为“背景 - 问题 - 方法 - 实验 - 结论”链条，突出创新点与验证证据。",
    bandColor: C.blue2,
  });
  addFooter(slide);
}

// 5 Methodology
{
  const slide = pptx.addSlide();
  addSlideBase(slide, "研究方法：需求调研驱动的系统设计与实现", 5, totalSlides);
  addBulletList(slide, [
    "文献调研界定区块链医疗场景",
    "问卷访谈提炼功能优先级",
    "采用 Fabric 联盟链总体方案",
    "以测试闭环验证系统有效性",
  ], { x: 0.72, y: 1.3, w: 3.2, h: 1.8, fontSize: 18.5 });
  addStep(slide, {
    x: 4.15, y: 1.28, w: 1.18, h: 0.92, num: 1,
    title: "痛点识别", body: "数据孤岛\n隐私泄露\n追溯不足", dotColor: C.blue,
  });
  addArrow(slide, 5.38, 1.6, 0.28);
  addStep(slide, {
    x: 5.72, y: 1.28, w: 1.18, h: 0.92, num: 2,
    title: "需求分析", body: "92份问卷\n四类用户\n七大模块", dotColor: C.green,
  });
  addArrow(slide, 6.95, 1.6, 0.28);
  addStep(slide, {
    x: 7.29, y: 1.28, w: 1.18, h: 0.92, num: 3,
    title: "架构设计", body: "联盟链\n双存储\nRBAC", dotColor: C.gold,
  });
  addStep(slide, {
    x: 4.15, y: 2.72, w: 1.18, h: 0.92, num: 4,
    title: "系统实现", body: "Go + Gin\nVue.js\n智能合约", dotColor: C.blue2,
  });
  addArrow(slide, 5.38, 3.04, 0.28);
  addStep(slide, {
    x: 5.72, y: 2.72, w: 1.18, h: 0.92, num: 5,
    title: "实验验证", body: "JMeter\nCaliper\n44用例", dotColor: C.red,
  });
  addArrow(slide, 6.95, 3.04, 0.28);
  addStep(slide, {
    x: 7.29, y: 2.72, w: 1.18, h: 0.92, num: 6,
    title: "结果总结", body: "100%通过\nTPS 150\n可扩展", dotColor: C.navy,
  });
  addInfoCard(slide, {
    x: 0.72,
    y: 4.18,
    w: 8.25,
    h: 0.7,
    title: "方法学特点",
    body: "采用“需求调研 + 原型实现 + 性能测试”三层证据结构，避免停留在概念验证。",
    bandColor: C.blue,
  });
  addFooter(slide);
}

// 6 Architecture
{
  const slide = pptx.addSlide();
  addSlideBase(slide, "方法阐述：系统总体架构", 6, totalSlides);
  addBulletList(slide, [
    "前后端分离，服务职责清晰",
    "Fabric 承担可信存证与授权",
    "MySQL 负责缓存与高频查询",
    "多层接口满足角色化访问",
  ], { x: 0.72, y: 1.3, w: 2.9, h: 1.8, fontSize: 18 });
  const layers = [
    { y: 1.3, title: "前端展示层", body: "Vue.js / Element UI / ECharts", fill: "E7F0F7" },
    { y: 2.0, title: "API 服务层", body: "Gin 路由 / 参数校验 / 权限控制", fill: "EDF5FA" },
    { y: 2.7, title: "区块链服务层", body: "Fabric SDK / ChannelQuery / ChannelExecute", fill: "F2F7FB" },
    { y: 3.4, title: "智能合约层", body: "15个链码函数 / 授权逻辑 / 审计记录", fill: "F7FAFC" },
    { y: 4.1, title: "数据存储层", body: "链上账本 + MySQL 双存储机制", fill: "FFFFFF" },
  ];
  layers.forEach((layer, i) => {
    slide.addShape(pptx.ShapeType.roundRect, {
      x: 4.15,
      y: layer.y,
      w: 4.55,
      h: 0.5,
      rectRadius: 0.06,
      fill: { color: layer.fill },
      line: { color: i === 0 ? C.blue : C.line, width: 1 },
    });
    slide.addText(layer.title, {
      x: 4.38,
      y: layer.y + 0.12,
      w: 1.8,
      h: 0.14,
      fontFace: "Microsoft YaHei",
      fontSize: 14,
      bold: true,
      color: C.navy,
      margin: 0,
    });
    slide.addText(layer.body, {
      x: 6.0,
      y: layer.y + 0.14,
      w: 2.3,
      h: 0.12,
      fontFace: "Microsoft YaHei",
      fontSize: 11.5,
      color: C.text,
      align: "right",
      margin: 0,
    });
  });
  slide.addShape(pptx.ShapeType.line, {
    x: 6.42, y: 1.8, w: 0, h: 2.32,
    line: { color: C.line, width: 1.2 },
  });
  addMetricCard(slide, { x: 0.86, y: 3.34, w: 1.4, h: 0.72, value: "5类", label: "系统角色", fill: C.white });
  addMetricCard(slide, { x: 2.42, y: 3.34, w: 1.4, h: 0.72, value: "7大", label: "业务模块", fill: C.white });
  addInfoCard(slide, {
    x: 0.72, y: 4.16, w: 2.95, h: 0.82,
    title: "设计原则",
    body: "分层架构\n模块化解耦\n安全与性能并重",
    bandColor: C.green,
  });
  addFooter(slide);
}

// 7 Topology
{
  const slide = pptx.addSlide();
  addSlideBase(slide, "方法阐述：联盟链拓扑与参与组织", 7, totalSlides);
  addMetricCard(slide, { x: 0.72, y: 1.35, w: 1.55, h: 0.76, value: "4", label: "Peer 组织", fill: C.white });
  addMetricCard(slide, { x: 2.45, y: 1.35, w: 1.55, h: 0.76, value: "8", label: "Peer 节点", fill: C.white });
  addMetricCard(slide, { x: 0.72, y: 2.3, w: 1.55, h: 0.76, value: "1", label: "Orderer", fill: C.white });
  addMetricCard(slide, { x: 2.45, y: 2.3, w: 1.55, h: 0.76, value: "15", label: "链码函数", fill: C.white });
  addInfoCard(slide, {
    x: 0.72, y: 3.45, w: 3.08, h: 1.24,
    title: "组织构成",
    body: "协和医院（TaobaoMSP）\n301医院（JDMSP）\n温江社区中心（WenjinMSP）\n监管中心（RegCenterMSP）",
    bandColor: C.blue,
  });
  slide.addShape(pptx.ShapeType.ellipse, {
    x: 5.02, y: 2.18, w: 1.52, h: 1.52,
    fill: { color: C.navy },
    line: { color: C.navy, transparency: 100 },
  });
  slide.addText("appchannel", {
    x: 5.02, y: 2.72, w: 1.52, h: 0.16,
    fontFace: "Microsoft YaHei",
    fontSize: 15,
    bold: true,
    color: C.white,
    align: "center",
    margin: 0,
  });
  const orgBoxes = [
    { x: 4.42, y: 1.06, t: "协和医院\nTaobaoMSP", c: C.blue },
    { x: 6.42, y: 1.06, t: "301医院\nJDMSP", c: C.green },
    { x: 4.42, y: 3.98, t: "温江社区\nWenjinMSP", c: C.gold },
    { x: 6.42, y: 3.98, t: "监管中心\nRegCenterMSP", c: C.red },
  ];
  orgBoxes.forEach((b) => {
    slide.addShape(pptx.ShapeType.roundRect, {
      x: b.x, y: b.y, w: 1.55, h: 0.78,
      rectRadius: 0.06,
      fill: { color: C.white },
      line: { color: b.c, width: 1.4 },
    });
    slide.addText(b.t, {
      x: b.x + 0.1, y: b.y + 0.2, w: 1.35, h: 0.28,
      fontFace: "Microsoft YaHei", fontSize: 12.5, bold: true,
      color: C.text, align: "center", margin: 0,
    });
    slide.addShape(pptx.ShapeType.line, {
      x: b.x + 0.77,
      y: b.y < 2.18 ? b.y + 0.78 : 3.0,
      w: 5.78 - (b.x + 0.77),
      h: b.y < 2.18 ? 1.4 : -0.3,
      line: { color: C.line, width: 1.1 },
    });
  });
  addInfoCard(slide, {
    x: 7.98, y: 2.16, w: 1.32, h: 1.06,
    title: "治理逻辑",
    body: "联盟链\n多机构共识\n审计可追溯",
    bandColor: C.navy,
  });
  addFooter(slide);
}

// 8 Innovations
{
  const slide = pptx.addSlide();
  addSlideBase(slide, "研究创新点：从可信存证走向可用系统", 8, totalSlides);
  const cards = [
    ["双存储机制", "关键数据上链\n查询数据入库\n兼顾安全性能", C.blue],
    ["患者主导授权", "跨院访问需审批\n授权记录上链\n隐私边界明确", C.green],
    ["完整病历历史链", "支持复诊随访\n补充记录只增不改\n形成连续病史", C.gold],
    ["全链路审计", "访问动作留痕\nTxID 可追踪\n责任界定更清晰", C.red],
  ];
  cards.forEach((card, i) => {
    const x = 0.82 + (i % 2) * 4.2;
    const y = 1.42 + Math.floor(i / 2) * 1.55;
    addInfoCard(slide, {
      x, y, w: 3.62, h: 1.18,
      title: card[0], body: card[1], bandColor: card[2],
      bodySize: 13,
    });
  });
  addInfoCard(slide, {
    x: 0.82, y: 4.72, w: 7.62, h: 0.42,
    title: "创新性总结",
    body: "该工作不是单点功能演示，而是将授权访问、订单流转和审计追踪整合进一套可部署系统。",
    bandColor: C.blue2,
  });
  addFooter(slide);
}

// 9 Authorization process
{
  const slide = pptx.addSlide();
  addSlideBase(slide, "关键机制：跨院授权访问流程", 9, totalSlides);
  const steps = [
    ["医生发起", "选择患者病历\n填写申请理由"],
    ["患者接收", "查看申请信息\n确认访问主体"],
    ["患者审批", "同意或拒绝\n形成链上授权"],
    ["医生查询", "获得授权后\n访问完整病历"],
    ["审计留痕", "访问日志记录\nTxID 全程追踪"],
  ];
  steps.forEach((s, i) => {
    const x = 0.64 + i * 1.82;
    addStep(slide, {
      x, y: 1.85, w: 1.45, h: 1.28, num: i + 1,
      title: s[0], body: s[1], dotColor: [C.blue, C.green, C.gold, C.blue2, C.red][i],
    });
    if (i < steps.length - 1) addArrow(slide, x + 1.5, 2.34, 0.18);
  });
  addBulletList(slide, [
    "授权决策权回归患者",
    "访问动作同步生成审计记录",
    "未授权场景自动拒绝跨院查看",
  ], { x: 0.92, y: 3.72, w: 3.4, h: 1.15, fontSize: 18 });
  addInfoCard(slide, {
    x: 4.72, y: 3.62, w: 4.08, h: 1.14,
    title: "机制价值",
    body: "将“共享可用”与“隐私可控”统一在一个流程中，是该系统最重要的业务闭环。",
    bandColor: C.navy,
  });
  addFooter(slide);
}

// 10 Experiment design
{
  const slide = pptx.addSlide();
  addSlideBase(slide, "实验设计：以需求样本与系统测试构成验证闭环", 10, totalSlides);
  addInfoCard(slide, {
    x: 0.72, y: 1.28, w: 2.65, h: 1.28,
    title: "数据来源",
    body: "92份有效问卷\n4类用户样本\n2024年1-2月调研",
    bandColor: C.blue,
  });
  addInfoCard(slide, {
    x: 0.72, y: 2.78, w: 2.65, h: 1.28,
    title: "测试设计",
    body: "44个测试用例\nJMeter 并发压测\nCaliper 链上评测",
    bandColor: C.green,
  });
  addInfoCard(slide, {
    x: 0.72, y: 4.28, w: 2.65, h: 0.68,
    title: "评测维度",
    body: "功能完整性 / 性能 / 安全性 / 兼容性",
    bandColor: C.gold,
  });
  slide.addTable([
    ["评价项", "实验设置", "目标或基线"],
    ["API响应", "50并发，500请求", "< 500ms"],
    ["链上写入", "20并发，100请求", "确认 < 5s"],
    ["并发能力", "100并发登录", "支持 100+ 用户"],
    ["链性能", "Caliper 压测", "TPS 100-200"],
  ], {
    x: 4.1, y: 1.42, w: 4.9,
    rowH: 0.42,
    colW: [1.2, 2.1, 1.4],
    border: { pt: 1, color: C.line },
    fill: C.white,
    fontFace: "Microsoft YaHei",
    fontSize: 12,
    color: C.text,
    bold: false,
    valign: "mid",
    margin: 0.06,
    autoFit: false,
    colorHeaders: C.white,
    fillHeader: C.navy,
    boldHeader: true,
  });
  addFooter(slide);
}

// 11 Questionnaire results
{
  const slide = pptx.addSlide();
  addSlideBase(slide, "实验验证：需求调研结果支撑系统功能优先级", 11, totalSlides);
  slide.addText("样本构成", {
    x: 0.82, y: 1.26, w: 1.2, h: 0.18,
    fontFace: "Microsoft YaHei", fontSize: 14, bold: true, color: C.navy, margin: 0,
  });
  slide.addChart(pptx.ChartType.doughnut, [
    {
      name: "样本",
      labels: ["患者", "医生/医务", "药店", "医院管理"],
      values: [41.3, 34.8, 12.0, 12.0],
    },
  ], {
    x: 0.72, y: 1.58, w: 3.55, h: 2.55,
    holeSize: 58,
    showLegend: true,
    legendPos: "b",
    chartColors: [C.blue, C.green, C.gold, C.red],
    chartArea: { fill: { color: C.white }, border: { color: C.line, pt: 1 } },
    showValue: true,
    showPercent: true,
    dataLabelColor: C.text,
    dataLabelPosition: "bestFit",
  });
  slide.addText("功能需求排序", {
    x: 4.72, y: 1.26, w: 1.6, h: 0.18,
    fontFace: "Microsoft YaHei", fontSize: 14, bold: true, color: C.navy, margin: 0,
  });
  slide.addChart(pptx.ChartType.bar, [
    {
      name: "占比",
      labels: ["电子病历创建查看", "跨院病历共享", "病历防篡改", "患者授权管理"],
      values: [87.0, 80.4, 77.2, 72.8],
    },
  ], {
    x: 4.56, y: 1.58, w: 4.22, h: 2.55,
    showLegend: false,
    showValue: true,
    valAxisMinVal: 0,
    valAxisMaxVal: 100,
    valAxisMajorUnit: 20,
    chartColors: [C.blue2],
    chartArea: { fill: { color: C.white }, border: { color: C.line, pt: 1 } },
    catAxisLabelColor: C.muted,
    valAxisLabelColor: C.muted,
    valGridLine: { color: C.line, pt: 0.6 },
    dataLabelColor: C.text,
    dataLabelPosition: "outEnd",
  });
  addInfoCard(slide, {
    x: 0.82, y: 4.36, w: 7.96, h: 0.5,
    title: "结果解释",
    body: "调研结果显示，病历共享、防篡改与患者授权是最优先能力，因此论文将其作为架构与合约设计核心。",
    bandColor: C.blue,
  });
  addFooter(slide);
}

// 12 Performance results
{
  const slide = pptx.addSlide();
  addSlideBase(slide, "实验验证：性能结果满足社区医疗场景需求", 12, totalSlides);
  slide.addChart(pptx.ChartType.bar, [
    {
      name: "平均响应",
      labels: ["用户登录", "病历列表查询", "病历创建上链"],
      values: [245, 380, 2800],
    },
  ], {
    x: 0.72, y: 1.48, w: 4.15, h: 2.78,
    showLegend: false,
    showValue: true,
    valAxisTitle: "毫秒",
    chartColors: [C.blue],
    chartArea: { fill: { color: C.white }, border: { color: C.line, pt: 1 } },
    catAxisLabelColor: C.muted,
    valAxisLabelColor: C.muted,
    valGridLine: { color: C.line, pt: 0.6 },
    dataLabelColor: C.text,
    dataLabelPosition: "outEnd",
  });
  addMetricCard(slide, { x: 5.32, y: 1.62, w: 1.32, h: 0.76, value: "150", label: "TPS", fill: C.white });
  addMetricCard(slide, { x: 6.82, y: 1.62, w: 1.32, h: 0.76, value: "2.5s", label: "平均链延迟", fill: C.white });
  addMetricCard(slide, { x: 8.32, y: 1.62, w: 0.76, h: 0.76, value: "0%", label: "错误率", fill: C.white, valueSize: 18 });
  slide.addTable([
    ["场景", "并发", "吞吐量", "结论"],
    ["用户登录", "50", "180 req/s", "响应稳定"],
    ["病历查询", "50", "120 req/s", "满足目标"],
    ["登录压测", "100", "-", "平均 580ms"],
    ["链上写入", "20", "7 req/s", "确认 < 5s"],
  ], {
    x: 5.02, y: 2.62, w: 4.0,
    rowH: 0.42,
    colW: [1.1, 0.7, 1.1, 0.9],
    border: { pt: 1, color: C.line },
    fill: C.white,
    fontFace: "Microsoft YaHei",
    fontSize: 11.5,
    color: C.text,
    margin: 0.05,
    autoFit: false,
    colorHeaders: C.white,
    fillHeader: C.navy,
    boldHeader: true,
  });
  addInfoCard(slide, {
    x: 0.72, y: 4.52, w: 8.3, h: 0.38,
    title: "评价",
    body: "查询类接口维持亚秒级响应；链上写入延迟高于数据库操作，但在医疗业务可接受区间内。",
    bandColor: C.green,
  });
  addFooter(slide);
}

// 13 Validation
{
  const slide = pptx.addSlide();
  addSlideBase(slide, "实验验证：功能、安全与兼容性全面通过", 13, totalSlides);
  slide.addChart(pptx.ChartType.bar, [
    {
      name: "通过数",
      labels: ["功能测试", "性能测试", "安全测试", "兼容性测试"],
      values: [24, 3, 6, 11],
    },
  ], {
    x: 0.72, y: 1.52, w: 4.05, h: 2.65,
    showLegend: false,
    showValue: true,
    chartColors: [C.green],
    chartArea: { fill: { color: C.white }, border: { color: C.line, pt: 1 } },
    catAxisLabelColor: C.muted,
    valAxisLabelColor: C.muted,
    valGridLine: { color: C.line, pt: 0.6 },
    dataLabelColor: C.text,
    dataLabelPosition: "outEnd",
  });
  addMetricCard(slide, { x: 5.32, y: 1.56, w: 1.52, h: 0.78, value: "44/44", label: "测试通过", fill: C.white });
  addMetricCard(slide, { x: 7.02, y: 1.56, w: 1.52, h: 0.78, value: "100%", label: "交易成功率", fill: C.white });
  addInfoCard(slide, {
    x: 5.18, y: 2.62, w: 3.74, h: 0.94,
    title: "安全验证",
    body: "密码 bcrypt 加密\n未授权跨院访问拦截\nSQL 注入与 XSS 防护通过",
    bandColor: C.red,
  });
  addInfoCard(slide, {
    x: 5.18, y: 3.74, w: 3.74, h: 0.94,
    title: "兼容性验证",
    body: "主流浏览器全通过\nWindows / Linux / macOS 可运行\n移动端响应式适配",
    bandColor: C.blue,
  });
  addFooter(slide);
}

// 14 Conclusion
{
  const slide = pptx.addSlide();
  addSlideBase(slide, "结论总结：系统实现了可信共享的核心闭环", 14, totalSlides);
  addBulletList(slide, [
    "完成社区医疗档案联盟链系统",
    "实现跨院共享与患者授权机制",
    "构建链上链下双存储体系",
    "验证性能、安全与可部署性",
    "为医疗区块链应用提供样例",
  ], { x: 0.78, y: 1.4, w: 3.6, h: 2.15, fontSize: 19 });
  addMetricCard(slide, { x: 5.02, y: 1.42, w: 1.32, h: 0.76, value: "4+1", label: "组织网络", fill: C.white });
  addMetricCard(slide, { x: 6.52, y: 1.42, w: 1.32, h: 0.76, value: "15", label: "链码函数", fill: C.white });
  addMetricCard(slide, { x: 8.02, y: 1.42, w: 1.0, h: 0.76, value: "5类", label: "角色", fill: C.white, valueSize: 18 });
  addInfoCard(slide, {
    x: 4.92, y: 2.56, w: 4.02, h: 1.02,
    title: "学术贡献",
    body: "将联盟链技术落实到社区医疗多角色业务流程，证明其在病历共享、授权控制与审计追踪中的可行性。",
    bandColor: C.navy,
  });
  addInfoCard(slide, {
    x: 0.78, y: 4.12, w: 8.16, h: 0.7,
    title: "一句话总结",
    body: "论文从实际场景出发，将“可信、可控、可追溯”的区块链特性转化为可运行的医疗信息系统能力。",
    bandColor: C.green,
  });
  addFooter(slide);
}

// 15 Outlook
{
  const slide = pptx.addSlide();
  addSlideBase(slide, "不足与展望", 15, totalSlides);
  addInfoCard(slide, {
    x: 0.82, y: 1.38, w: 3.8, h: 2.1,
    title: "当前不足",
    body: "TPS 约 150，规模扩展仍有限\n高阶隐私技术尚未引入\n跨链互操作能力有待补强\n移动端与影像数据支持不足",
    bandColor: C.red,
  });
  addInfoCard(slide, {
    x: 5.02, y: 1.38, w: 3.98, h: 2.1,
    title: "未来方向",
    body: "引入 IPFS 处理大文件\n结合国密或零知识增强隐私\n扩展更多医疗机构接入\n结合 AI 做辅助诊断分析",
    bandColor: C.blue,
  });
  slide.addShape(pptx.ShapeType.roundRect, {
    x: 1.35, y: 4.08, w: 7.2, h: 0.78,
    rectRadius: 0.08,
    fill: { color: C.navy },
    line: { color: C.navy, transparency: 100 },
  });
  slide.addText("汇报完毕，恳请各位老师批评指正", {
    x: 1.55, y: 4.34, w: 6.8, h: 0.2,
    fontFace: "Microsoft YaHei",
    fontSize: 22,
    bold: true,
    color: C.white,
    align: "center",
    margin: 0,
  });
  addFooter(slide);
}

async function main() {
  const out = path.join(root, "答辩PPT-学术汇报版-15页.pptx");
  await pptx.writeFile({ fileName: out });
  console.log(out);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
