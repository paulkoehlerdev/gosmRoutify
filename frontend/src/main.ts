import './assets/main.css'
import 'bootstrap/scss/bootstrap.scss'
import "bootstrap-icons/font/bootstrap-icons.css";

import { createApp } from 'vue'

import App from './App.vue'
import mitt from 'mitt'

const app = createApp(App)

const eventBus = mitt();
app.provide('eventBus', eventBus);

app.mount('#app')
