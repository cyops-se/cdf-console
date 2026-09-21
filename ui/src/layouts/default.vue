<template>
  <v-layout>
    <AppHeader @togglemenu="showmenu = !showmenu" />

    <v-navigation-drawer v-model="showmenu" permanent>
      <v-list>
        <div v-for="(link) in menuitems" :key="link.text">
          <v-list-item v-if="!link.items" :prepend-icon="link.icon" :href="link.href" :title="link.text"></v-list-item>

          <v-list-group v-if="link.items" :key="link.text" :prepend-icon="link.icon" :value="false" no-action>
            <template v-slot:activator="{ props }">
              <v-list-item v-bind="props" :key="link.text">
                <v-list-item-title v-text="link.text" />
              </v-list-item>
            </template>

            <v-list-item v-for="sublink in link.items" :key="sublink.text" :href="sublink.href" density="compact">
              <v-list-item-title>{{ sublink.text }}</v-list-item-title>
            </v-list-item>
          </v-list-group>
        </div>
      </v-list>

      <template v-slot:append>
        <div class="pa-3 text-center version-info">
          <div v-if="sysinfo.gitversion" class="text-caption text-grey">
            Version: {{ sysinfo.gitversion }}
          </div>
          <div v-if="sysinfo.gitcommit" class="text-caption text-grey">
            Commit: {{ sysinfo.gitcommit.substring(0, 7) }}
          </div>
        </div>
      </template>
    </v-navigation-drawer>

    <v-main class="main-content">
      <div class="content-wrapper">
        <v-alert v-model="app.alert" class="mb-2 mx-2 mt-2" :text="app.alertmsg" type="error" closable
          style="max-height: 60px"></v-alert>
        <router-view />
      </div>

      <v-snackbar v-model="app.snackbar" elevation="20" :color="app.snackcolor" :timeout="timeout">
        {{ app.snacktext }}
      </v-snackbar>
    </v-main>
  </v-layout>
</template>
<script setup>
import { useAppStore } from '@/stores/app'
import axios from 'axios'

const app = useAppStore()

const disconnected = ref(false)
const timeout = ref(2000)
const showmenu = ref(true)
const sysinfo = ref({
  gitversion: '',
  gitcommit: ''
})

const menuitems = [
  { text: 'Dashboard', icon: 'mdi-home', href: '/ui' },
  { text: 'Systems', icon: 'mdi-sitemap', href: '/ui/systems/list' },
  { text: 'Devices', icon: 'mdi-switch', href: '/ui/devices' },
  { text: 'Onboarding', icon: 'mdi-rocket', href: '/ui/onboarding' },
  {
    text: 'Monitoring', icon: 'mdi-security-network', href: '/ui/monitoring', items: [
      { text: 'Grey', icon: 'mdi-security-network', href: '/ui/monitoring/grey' },
      { text: 'White', icon: 'mdi-security-network', href: '/ui/monitoring/white' },
      { text: 'Black', icon: 'mdi-security-network', href: '/ui/monitoring/black' }]
  },
]

onMounted(async () => {
  try {
    const response = await axios.get('/api/system/sysinfo')
    if (response.data.success && response.data.sysinfo) {
      sysinfo.value = response.data.sysinfo
    }
  } catch (error) {
    console.error('Failed to fetch system info:', error)
  }
})

</script>

<style scoped lang="scss">
.subitem {
  font-size: .8rem;
  margin: 5px;
}

.version-info {
  border-top: 1px solid rgba(0, 0, 0, 0.12);
  line-height: 1.4;
}

.version-info .text-caption {
  font-size: 0.7rem;
}

.main-content {
  height: 100vh;
  overflow-y: auto;
  overflow-x: hidden;
}

.content-wrapper {
  padding: 0;
  width: 100%;
}
</style>
