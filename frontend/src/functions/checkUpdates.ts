import { ElMessageBox, ElNotification } from 'element-plus';
import { h } from 'vue';
import { CheckUpdates, Update as U } from '../../wailsjs/go/main/App.js';
import Updater from '../components/Updater.vue';

export function checkUpdates(type: 'load' | 'menu') {
    CheckUpdates()
        .then((res: any) => {
            if (res.Error) {
                ElNotification({
                    title: 'Error while checking for updates!',
                    message: res.Error,
                    type: 'error',
                    position: 'bottom-right',
                });
                return;
            }
            if (type == 'menu' && res.IsLatest) {
                ElMessageBox({
                    message: 'You are using latest version',
                }).catch(() => {});
            } else if (!res.IsLatest) {
                ElMessageBox({
                    title: 'A new version is avaible!\nDo you want to update?',
                    message: h(Updater, {
                        version: res.Version,
                        description: res.Description,
                    }),
                    confirmButtonText: 'Yes',
                    showCancelButton: true,
                    cancelButtonText: 'No',
                    showClose: false,
                    closeOnClickModal: false,
                    closeOnHashChange: false,
                    closeOnPressEscape: false,
                    center: true,
                })
                    .then(() => {
                        ElMessageBox({
                            title: 'Warning!',
                            message:
                                'The app will now exit and will be re-opened automatically',
                            type: 'warning',
                            showClose: false,
                            closeOnClickModal: false,
                            closeOnHashChange: false,
                            closeOnPressEscape: false,
                            center: true,
                        }).then(() => {
                            U();
                        });
                    })
                    .catch(() => {});
            }
        })
        .catch(() => {
            ElNotification({
                title: 'Error while checking for updates!',
                type: 'error',
                position: 'bottom-right',
            });
        });
}
