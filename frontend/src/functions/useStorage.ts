import { useStorage } from '@vueuse/core';
import { HTTPItem, WSItem } from '../types';

export const httpTabItem = useStorage<HTTPItem[]>('httpTab', []);

export const wsTabItem = useStorage<WSItem[]>('wsTab', []);
