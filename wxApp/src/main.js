import {
	createSSRApp
} from "vue";
import App from "./App.vue";
import pinia from "./store";
import uviewPlus from 'uview-plus';

export function createApp() {
	const app = createSSRApp(App);
	app.use(pinia);
	app.use(uviewPlus);
	return {
		app,
	};
}
