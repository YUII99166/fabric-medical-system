<script>
export default {
  name: 'MenuItem',
  functional: true,
  props: {
    icon: {
      type: String,
      default: ''
    },
    title: {
      type: String,
      default: ''
    }
  },
  render(h, context) {
    const { icon, title } = context.props
    const vnodes = []

    if (icon) {
      // 判断是否为图片路径（以 /image/ 开头）
      if (icon.startsWith('/image/')) {
        vnodes.push(
          <img 
            src={icon} 
            style="width: 20px; height: 20px; margin-right: 10px; vertical-align: middle; object-fit: contain;"
            alt="icon"
          />
        )
      }
      // 判断是否为 Element UI 图标（以 el-icon- 开头）
      else if (icon.startsWith('el-icon-')) {
        vnodes.push(<i class={icon} style="margin-right: 10px;"/>)
      } 
      // 使用 SVG 图标
      else {
        vnodes.push(<svg-icon icon-class={icon}/>)
      }
    }

    if (title) {
      vnodes.push(<span slot='title'>{(title)}</span>)
    }
    return vnodes
  }
}
</script>
