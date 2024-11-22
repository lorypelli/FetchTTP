import { useStorage } from '@vueuse/core';

export const httpTabItem = useStorage('httpTab', []);

export const wsTabItem = useStorage('wsTab', []);