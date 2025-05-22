import {POSITION} from 'vue-toastification';
import "vue-toastification/dist/index.css";
import { defineNuxtPlugin } from '#app';

import Toast from "vue-toastification";
import type { PluginOptions } from "vue-toastification";
import "vue-toastification/dist/index.css";

export default defineNuxtPlugin(nuxtApp => {
    const options: PluginOptions = {
        position: POSITION.BOTTOM_RIGHT,
    };
    nuxtApp.vueApp.use(Toast, options);
})