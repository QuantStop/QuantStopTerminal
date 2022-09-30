import { createApp } from 'vue'
import App from './App.vue'
import * as appRouter from './router'
import "bootstrap";
import "bootstrap/dist/css/bootstrap.min.css";
import { FontAwesomeIcon } from './plugins/font-awesome'
import $bus from "./websocket/events.js"

const app = createApp(App)
app.config.globalProperties.$bus = $bus;
app.use(appRouter.routeConfig);
//app.use(store);
app.component("font-awesome-icon", FontAwesomeIcon)
app.mount('body')
