/* eslint-disable */
import Vue from 'vue'
import Router from 'vue-router'
import echartview from '@/components/echartview'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'echartview',
      component: echartview
    }
  ]
})
