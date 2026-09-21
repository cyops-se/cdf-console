import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', {
    state: () => ({
        alert: false,
        alertmsg: '',
        snackbar: false,
        snacktext: '',
        snackcolor: 'success',
        system: {},
        device: {},
    }),
    actions: {
        clearalert() {
            this.alert = false
            this.alertmsg = ''
        },
        setalert(msg) {
            this.alert = true
            this.alertmsg = msg
        },
        showsnack(msg, color) {
            this.snacktext = msg
            this.snackbar = true
            this.snackcolor = color ? color : 'success'
            console.log('snack: ', msg)
        },
    },
    persist: true,
})
