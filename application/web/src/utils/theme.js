// 角色主题配置
export const roleThemes = {
  admin: {
    primary: '#409EFF',
    menuBg: '#2c3e50',
    menuHover: '#34495e',
    subMenuBg: '#1f2d3d',
    subMenuHover: '#263445',
    gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    lightBg: '#ecf5ff'
  },
  doctor: {
    primary: '#67C23A',
    menuBg: '#2d5f3f',
    menuHover: '#3a7a52',
    subMenuBg: '#1e4029',
    subMenuHover: '#2a5438',
    gradient: 'linear-gradient(135deg, #11998e 0%, #38ef7d 100%)',
    lightBg: '#f0f9ff'
  },
  patient: {
    primary: '#5DADE2',
    menuBg: '#34495e',
    menuHover: '#415b76',
    subMenuBg: '#2c3e50',
    subMenuHover: '#34495e',
    gradient: 'linear-gradient(135deg, #89CFF0 0%, #5DADE2 100%)',
    lightBg: '#EBF5FB'
  },
  drugstore: {
    primary: '#9C27B0',
    menuBg: '#4a2d5f',
    menuHover: '#5f3a7a',
    subMenuBg: '#311e40',
    subMenuHover: '#422954',
    gradient: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
    lightBg: '#f5f0ff'
  }
}

// 应用主题
export function applyTheme(role) {
  const theme = roleThemes[role] || roleThemes.admin
  const root = document.documentElement
  
  // 设置 CSS 变量
  root.style.setProperty('--theme-primary', theme.primary)
  root.style.setProperty('--theme-menu-bg', theme.menuBg)
  root.style.setProperty('--theme-menu-hover', theme.menuHover)
  root.style.setProperty('--theme-submenu-bg', theme.subMenuBg)
  root.style.setProperty('--theme-submenu-hover', theme.subMenuHover)
  root.style.setProperty('--theme-gradient', theme.gradient)
  root.style.setProperty('--theme-light-bg', theme.lightBg)
  
  // 保存当前主题
  localStorage.setItem('currentTheme', role)
}

// 获取当前主题
export function getCurrentTheme() {
  return localStorage.getItem('currentTheme') || 'admin'
}

// 获取主题配置
export function getThemeConfig(role) {
  return roleThemes[role] || roleThemes.admin
}
