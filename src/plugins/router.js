import Vue from 'vue';
import VueRouter from 'vue-router'
import Auth from '../components/Auth'
import List from '../components/List'

Vue.use(VueRouter)

export default new VueRouter({
    mode: 'history',
    routes: [
        { path: '/auth', component: Auth },
        { path: '/list', component: List }
    ]
})