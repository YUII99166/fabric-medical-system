import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'             the icon show in the sidebar
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
export const constantRoutes = [{
  path: '/login',
  component: () => import('@/views/login/index'),
  hidden: true
},

{
  path: '/404',
  component: () => import('@/views/404'),
  hidden: true
},

{
  path: '/',
  component: Layout,
  redirect: '/main',
  children: [{
    path: 'main',
    name: 'main',
    component: () => import('@/views/mainPage/list/index'),
    meta: {
      title: '主页',
      icon: 'el-icon-s-home'
    }
  }]
}
]

/**
 * asyncRoutes 异步路由
 * the routes that need to be dynamically loaded based on user roles   异步路由是根据用户角色动态加载的路由。
 * 每个路由对象包括路径（path）、组件（component）、重定向（redirect）、路由名称（name）、子路由（children）等元素
 */
export const asyncRoutes = [
  {
    path: '/account',
    component: Layout,
    redirect: '/account/architecture',
    name: 'Account',
    alwaysShow: true,
    meta: {
      roles: ['admin'],
      title: '监管中心',
      icon: '/image/系统管理员.png'
    },
    children: [
      {
        path: 'architecture',
        name: 'Architecture',
        component: () => import('@/views/account/architecture/index'),
        meta: {
          roles: ['admin'],
          title: '架构概览',
          icon: '/image/3-联盟链部署.png'
        }
      },
      {
        path: 'statistics',
        name: 'Statistics',
        component: () => import('@/views/account/statistics/index'),
        meta: {
          roles: ['admin'],
          title: '数据统计',
          icon: 'el-icon-data-line'
        }
      },
      {
        path: 'activity-monitor',
        name: 'ActivityMonitor',
        component: () => import('@/views/account/activity-monitor/index'),
        meta: {
          roles: ['admin'],
          title: '日志监控',
          icon: 'el-icon-document-copy'
        }
      },
      {
        path: 'all',
        name: 'AccountAll',
        component: () => import('@/views/account/list/index'),
        meta: {
          roles: ['admin'],
          title: '用户管理',
          icon: 'el-icon-user'
        }
      },
      {
        path: 'delete',
        name: 'AccountDelete',
        component: () => import('@/views/account/delete/index'),
        meta: {
          roles: ['admin'],
          title: '删除用户',
          icon: 'el-icon-delete'
        }
      },

      {
        path: 'add',
        name: 'Add',
        component: () => import('@/views/account/add/index'),
        meta: {
          roles: ['admin'],
          title: '新增账户',
          icon: 'el-icon-plus'
        }
      },

      {
        path: 'detail',
        name: 'AccountDetail',
        component: () => import('@/views/account/detail/index'),
        hidden: true,
        meta: {
          roles: ['admin'],
          title: '用户详情',
          icon: 'el-icon-document'
        }
      },

      {
        path: 'edit',
        name: 'AccountEdit',
        component: () => import('@/views/account/edit/index'),
        hidden: true,
        meta: {
          roles: ['admin'],
          title: '编辑用户',
          icon: 'el-icon-edit'
        }
      }
    ]
  },

  {
    path: '/prescription',
    component: Layout,
    redirect: '/prescription/all',
    name: 'Prescription',
    alwaysShow: true,
    meta: {
      roles: ['admin','doctor','patient'],
      title: '病历',
      icon: '/image/男医生.png'
    },
    children: [
      {
        path: 'all',
        name: 'PrescriptionAll',
        component: () => import('@/views/prescription/list/index'),
        meta: {
          roles: ['admin','doctor'],
          title: '所有病历',
          icon: 'el-icon-folder-opened'
        }
      },
      {
        path: 'mine',
        name: 'PrescriptionOfMine',
        component: () => import('@/views/prescription/mine/index'),
        meta: {
          roles: ['admin','patient'],
          title: '我的病历',
          icon: '/image/病人-男.png'
        }
      },
      {
        path: 'health-profile',
        name: 'HealthProfile',
        component: () => import('@/views/patient/health-profile/index'),
        meta: {
          roles: ['patient'],
          title: '健康档案',
          icon: 'el-icon-data-line'
        }
      },
      {
        path: 'access-trace',
        name: 'AccessTrace',
        component: () => import('@/views/patient/access-trace/index'),
        meta: {
          roles: ['patient'],
          title: '隐私溯源',
          icon: 'el-icon-view'
        }
      },
      {
        path: 'add',
        name: 'Add',
        component: () => import('@/views/prescription/add/index'),
        meta: {
          roles: ['admin','doctor'],
          title: '新增病历',
          icon: 'el-icon-edit-outline'
        }
      },
      {
        path: 'patient-search',
        name: 'PatientSearch',
        component: () => import('@/views/prescription/patient-search/index'),
        meta: {
          roles: ['admin','doctor'],
          title: '患者病历查询',
          icon: 'el-icon-search'
        }
      },
      {
        path: 'detail',
        name: 'PrescriptionDetail',
        component: () => import('@/views/prescription/detail/index'),
        hidden: true,
        meta: {
          roles: ['admin','doctor','patient'],
          title: '病历详情',
          icon: 'el-icon-document'
        }
      }
    ]
  },

  {
    path: '/authorization',
    component: Layout,
    children: [{
      path: '',
      name: 'Authorization',
      component: () => import('@/views/authorization/index'),
      meta: {
        roles: ['admin','patient','doctor'],
        title: '授权管理',
        icon: 'el-icon-key'
      }
    }]
  },

  {
    path: '/drug',
    component: Layout,
    redirect: '/drug/all',
    name: 'Drug',
    alwaysShow: true,
    meta: {
      roles: ['admin','drugstore','patient'],
      title: '药品订单',
      icon: '/image/药店-copy.png'
    },
    children: [{
      path: 'all',
      name: 'DrugAll',
      component: () => import('@/views/drugOrder/list/index'),
      meta: {
        roles: ['admin','drugstore'],
        title: '所有订单',
        icon: 'el-icon-tickets'
      }
    },
      {
        path: 'mine',
        name: 'DrugMine',
        component: () => import('@/views/drugOrder/mine/index'),
        meta: {
          roles: ['admin', 'patient'],
          title: '我的订单',
          icon: 'el-icon-shopping-bag-2'
        }
      },
      {
        path: '/addDrug',
        name: 'AddDrug',
        component: () => import('@/views/drugOrder/add/index'),
        meta: {
          roles: ['admin','drugstore'],
          title: '新增订单',
          icon: 'el-icon-circle-plus-outline'
        }
      }]
  },

  {
    path: '/about',
    component: Layout,
    children: [{
      path: '',
      name: 'About',
      component: () => import('@/views/about/index'),
      meta: {
        roles: ['admin', 'doctor', 'patient', 'drugstore'],
        title: '关于',
        icon: 'el-icon-info'
      }
    }]
  },

  // 404 page must be placed at the end !!!
  {
    path: '*',
    redirect: '/404',
    hidden: true
  }
]

const createRouter = () => new Router({
  base: '/',
  // mode: 'history', // require service support
  scrollBehavior: () => ({
    y: 0
  }),
  routes: constantRoutes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}

export default router
