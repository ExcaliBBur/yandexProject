import { createRouter, createWebHistory } from 'vue-router'
import CalculatorMain from '../components/CalculatorMain.vue'
import CalculatorSettings from '../components/CalculatorSettings.vue'
import ComputingResourses from '../components/ComputingResourses.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/calculator',
      component: CalculatorMain
    },
    {
      path: '/calculator-settings',
      component: CalculatorSettings
    },
    {
      path: '/computing-resourses',
      component: ComputingResourses
    },
    {
      path: '/:catchAll(.*)',
      component: CalculatorMain
    }
  ]
})

export default router
