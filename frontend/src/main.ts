import './assets/main.css'
import 'bootstrap/scss/bootstrap.scss'
import "bootstrap-icons/font/bootstrap-icons.css";

import { createApp } from 'vue'

import App from './App.vue'
import mitt from 'mitt'
import L from 'leaflet'

// eslint-disable-next-line
delete L.Icon.Default.prototype._getIconUrl
// eslint-disable-next-line
L.Icon.Default.mergeOptions({
  iconRetinaUrl: import('leaflet/dist/images/marker-icon-2x.png'),
  iconUrl: import('leaflet/dist/images/marker-icon.png'),
  shadowUrl: import('leaflet/dist/images/marker-shadow.png')
})

const app = createApp(App)

const eventBus = mitt();
app.provide('eventBus', eventBus);

app.mount('#app')
