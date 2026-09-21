<template>
  <v-container fluid>
    <v-row>
      <v-col cols="12">
        <v-toolbar color="transparent">
          <v-toolbar-title class="text-h6">
            Device Dashboard
          </v-toolbar-title>
          <v-spacer></v-spacer>
          <v-btn
            icon="mdi-help-circle-outline"
            size="small"
            variant="text"
            color="primary"
            title="View dashboard help"
            @click="openHelp"
            class="mr-2"
          ></v-btn>
          <v-btn icon="mdi-counter" @click="resetStats" title="Reset all traffic statistics"></v-btn>
          <v-btn icon="mdi-refresh" @click="refresh"></v-btn>
        </v-toolbar>
      </v-col>
    </v-row>

    <v-row>
      <v-col v-for="device in devices" :key="device.id" cols="12" sm="6" md="4" lg="3">
        <v-card :color="device.status === 'reachable' ? 'green-lighten-4' : 'red-lighten-4'" elevation="2">
          <v-card-title>
            <div class="d-flex align-center justify-space-between">
              <span class="text-h6">{{ device.name || device.hostname }}</span>
              <v-chip :color="device.status === 'reachable' ? 'success' : 'error'" size="small" label>
                {{ device.status }}
              </v-chip>
            </div>
          </v-card-title>

          <v-card-subtitle>
            {{ device.serialno }}
          </v-card-subtitle>

          <v-card-text>
            <!-- Interface Statistics -->
            <div v-if="getInterfacesWithStats(device).length > 0">
              <v-divider class="mb-3"></v-divider>
              <div v-for="iface in getInterfacesWithStats(device)" :key="iface.name" class="mb-3">
                <div class="d-flex align-center justify-space-between mb-2">
                  <div>
                    <div class="text-subtitle-2">
                      {{ iface.name }}
                      <span class="text-caption font-weight-light text-grey ml-2">{{
                        }}</span>
                    </div>
                  </div>
                  <v-chip size="x-small" :color="iface.stats.operstate === 'up' ? 'success' : 'warning'">
                    {{ iface.stats.operstate }}
                  </v-chip>
                </div>

                <!-- RX/TX Stats -->
                <div class="mb-2">
                  <div class="d-flex align-center justify-space-between mb-1">
                    <span class="text-caption text-grey">RX:</span>
                    <span class="text-body-2 font-weight-medium text-blue">{{ formatBytes(iface.stats.rxbytes) }}</span>
                  </div>

                  <div class="d-flex align-center justify-space-between mb-1">
                    <span class="text-caption text-grey">TX:</span>
                    <span class="text-body-2 font-weight-medium text-orange">{{ formatBytes(iface.stats.txbytes)
                      }}</span>
                  </div>
                </div>

                <!-- Packet counters -->
                <div class="d-flex justify-space-between text-caption text-grey">
                  <span>{{ formatNumber(iface.stats.rxpackets) }} pkts ↓</span>
                  <span>{{ formatNumber(iface.stats.txpackets) }} pkts ↑</span>
                </div>
              </div>
            </div>

            <!-- No interface stats -->
            <div v-else class="text-caption text-grey">
              No interface statistics available
            </div>

            <!-- Last check timestamp -->
            <v-divider class="my-2"></v-divider>
            <div class="text-caption text-grey">
              Last check: {{ formatTimestamp(device.lastcheck) }}
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
  import { ref, onMounted } from 'vue'
  import { useRouter } from 'vue-router'
  import { device } from '@/plugins/device'

  const router = useRouter()
  const devices = ref([])

  onMounted(async () => {
    await refresh()

    // Auto-refresh every 30 seconds
    setInterval(refresh, 30000)
  })

  async function refresh() {
    const allDevices = await device.alldevices()
    devices.value = allDevices

    console.log('Dashboard refresh')
    console.log('Number of devices:', devices.value ? devices.value.length : 0)
    await device.checklist(devices.value)
  }

  async function resetStats() {
    const result = await device.resetstats()
    if (result.success) {
      console.log('Traffic stats reset successfully')
      await refresh()
    } else {
      console.error('Failed to reset traffic stats:', result.errors)
    }
  }

  function openHelp() {
    // Open help in a new window
    const helpUrl = router.resolve({
      path: '/help',
      query: { doc: 'Sidor/dashboard.md' }
    }).href
    window.open(helpUrl, '_blank', 'width=1200,height=800')
  }

  function getInterfacesWithStats(device) {
    const interfaces = []

    // Add WAN interface if it has stats
    if (device.wanstats && device.wanstats.rxbytes > 0) {
      interfaces.push({
        name: 'WAN',
        ipaddress: device.wan,
        stats: device.wanstats
      })
    }

    // Add LAN interface if it has stats
    if (device.lanstats && device.lanstats.rxbytes > 0) {
      interfaces.push({
        name: 'LAN',
        ipaddress: device.lan,
        stats: device.lanstats
      })
    }

    // Add APN interface if it has stats
    if (device.apn4gstats && device.apn4gstats.rxbytes > 0) {
      interfaces.push({
        name: 'APN',
        ipaddress: device.apn4g,
        stats: device.apn4gstats
      })
    }

    return interfaces
  }

  function formatBytes(bytes) {
    if (!bytes || bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  function formatNumber(num) {
    if (!num || num === 0) return '0'
    if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M'
    if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
    return num.toString()
  }

  function formatTimestamp(timestamp) {
    if (!timestamp) return 'Never'
    const date = new Date(timestamp)
    return date.toLocaleString()
  }
</script>
