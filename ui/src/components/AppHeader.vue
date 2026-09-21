<template>
    <v-app-bar fixed>
        <v-app-bar-nav-icon @click="$emit('togglemenu')" />
        <v-app-bar-title>CDF-CONSOLE</v-app-bar-title>
        <template v-slot:append>
            <v-btn @click="openHelp" size="small" variant="text" title="Help Documentation">Help</v-btn>
            <v-btn icon="mdi-heart" :color="connected"></v-btn>
            <v-btn icon="mdi-theme-light-dark" @click="toggleTheme"></v-btn>
        </template>
    </v-app-bar>
</template>

<script setup>
import { useTheme } from 'vuetify'
import { useRouter } from 'vue-router'
import axios from 'axios'

const theme = useTheme()
const router = useRouter()
const connected = ref('warning')

onMounted(async () => {
    setInterval(async () => {
        axios.get('/api/system/sysinfo').then(() => { connected.value = 'success' }).catch(() => connected.value = 'error')
    }, 5000)
})

function toggleTheme() {
    theme.global.name.value = theme.global.current.value.dark ? 'light' : 'dark'
}

function openHelp() {
    // Open help documentation in a new window
    const helpUrl = router.resolve({
        path: '/help'
    }).href
    window.open(helpUrl, '_blank', 'width=1200,height=800')
}
</script>

<style scoped lang="sass">
</style>
