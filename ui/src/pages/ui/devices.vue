<template>
  <!-- <v-card width="100%" height="100%">
    <v-toolbar color="transparent">
      <v-toolbar-title class="text-h6" text="Devices"></v-toolbar-title>

      <v-text-field v-model="search" label="Search" prepend-inner-icon="mdi-magnify" variant="outlined" hide-details
        single-line></v-text-field>
      <template v-slot:append>
        <v-btn class="ml-4" color="primary" icon="mdi-refresh" @click="refresh"></v-btn>
        <v-btn class="ml-2" @click="scandialog = true">Scan</v-btn>
        <v-btn class="ml-2" color="success" icon="mdi-plus" @click="dialog = true"></v-btn>
      </template>
</v-toolbar>
<v-card-text>
  <v-data-table density="compact" :search="search" :headers="headers" :items="items" @click:row="onClick">
    <template v-slot:item.lan="{ item }">
          <div class="text-end">
            <v-chip :color="item.lanpresent ? 'success' : 'error'" :text="item.lan || ''" class="text-uppercase" size="small"
              label></v-chip>
          </div>
        </template>
    <template v-slot:item.wan="{ item }">
          <div class="text-end">
            <v-chip :color="item.wanpresent ? 'success' : 'error'" :text="item.wan || ''" class="text-uppercase" size="small"
              label></v-chip>
          </div>
        </template>
    <template v-slot:item.status="{ item }">
          <div class="text-end">
            <v-chip :color="item.status == 'ok' ? 'success' : 'error'" :text="item.status" class="text-uppercase"
              size="small" label></v-chip>
          </div>
        </template>
    <template v-slot:item.actions="{ item }">
          <v-icon size="small" @click.stop="deleteItem(item)">mdi-delete</v-icon>
        </template>
  </v-data-table>
</v-card-text>
</v-card>
<v-dialog max-width="500" v-model="dialog">
  <template v-slot:default="{ isActive }">
      <v-card title="Device">
        <v-card-text>
          <v-text-field v-model="app.device.name" label="Name" required></v-text-field>
          <v-text-field v-model="app.device.hostname" label="Hostname" required></v-text-field>
          <v-text-field v-model="app.device.serialno" label="Serial" required></v-text-field>
        </v-card-text>
        <v-card-actions>
          <v-btn text="Cancel" @click="isActive.value = false"></v-btn>
          <v-spacer></v-spacer>
          <v-btn color="success" text="Save" @click="save"></v-btn>
        </v-card-actions>
      </v-card>
    </template>
</v-dialog>
<v-dialog max-width="500" v-model="scandialog">
  <template v-slot:default="{ isActive }">
      <v-card title="Device">
        <v-card-text>
          <v-text-field v-model="scancidr" label="Network to scan" required></v-text-field>
          <v-label>Scanning: {{ scanning }}</v-label>
          <v-progress-linear v-model="progress" :height="12"></v-progress-linear>
        </v-card-text>
        <v-card-actions>
          <v-btn text="Cancel" @click="isActive.value = false"></v-btn>
          <v-spacer></v-spacer>
          <v-btn color="success" text="Scan" @click="scan"></v-btn>
        </v-card-actions>
      </v-card>
    </template>
</v-dialog> -->
  <DeviceList />
</template>

<script setup>
  // import axios from "axios"
  // import { onMounted } from "vue"
  // import { useRouter } from 'vue-router'
  // import { expandCidr, parseCidr } from "cidr-tools"
  // import { useAppStore } from '@/stores/app'
  // import device from '@/plugins/device'

  // const router = useRouter()
  // const app = useAppStore()
  // const items = ref([])
  // const networks = ref([])
  // const dialog = ref(false)
  // const scandialog = ref(false)
  // // const headers = ref([{ title: 'Hostname', key: 'hostname', align: 'start' }, { title: 'Role', key: 'role', align: 'end' }, { title: 'Hub', key: 'hubid', align: 'end' }, { title: 'Serial', key: 'serialno', align: 'end' }, { title: 'WAN', key: 'wan', align: 'end' }, { title: 'APN 4G', key: 'apn4g', align: 'end' }, { title: 'LAN', key: 'lan', align: 'end' }, { title: 'Public key', key: 'pubkey', align: 'end' }, { title: 'Status', key: 'status', align: 'end' }, { title: 'Actions', key: 'actions', align: 'end' }])
  // const headers = ref([{ title: 'Hostname', key: 'hostname', align: 'start' }, { title: 'Role', key: 'role', align: 'end' }, { title: 'Hub', key: 'hubid', align: 'end' }, { title: 'Serial', key: 'serialno', align: 'end' }, { title: 'WAN', key: 'wan', align: 'end' }, { title: 'LAN', key: 'lan', align: 'end' }, { title: 'Status', key: 'status', align: 'end' }, { title: 'Actions', key: 'actions', align: 'end' }])
  // const search = ref('')
  // const scancidr = ref('192.168.1.0/24')
  // const scanning = ref('')
  // const progress = ref(0)

  onMounted(async () => {
    refresh()
  })

  async function refresh() {
    // var response = await device.alldevices()
    // items.value = response

    // for (var dev of items.value) {
    //   dev.status = 'ok'
    //   // var reachable = await device.present(element.wan)
    //   var reachable = dev.wanpresent || dev.lanpresent
    //   dev.useip = dev.wanpresent ? dev.wan : dev.lan
    //   dev.usemask = dev.wanpresent ? dev.wanmask : dev.lanmask
    //   if (reachable) {
    //     var serial = await device.serial(dev.useip)
    //     if (serial.success) {
    //       if (dev.serialno !== serial.serial) {
    //         dev.status = 'Serial mismatch'
    //         dev.serialno += ' (' + serial.serial + ')'
    //       }
    //     } else {
    //       dev.status = 'No SSH keys'
    //     }
    //   } else {
    //     dev.status = 'Not reachable'
    //   }

    //   if (!dev.lan) dev.lan = dev.ip + '/' + dev.mask
    // }
  }

  // async function save() {
  //   dialog.value = false
  //   delete app.device.id
  //   const response = await device.save(app.device) // axios.put("/api/data/devices", device.value)
  //   refresh()
  // }

  // async function deleteItem(item) {
  //   const response = await axios.delete("/api/data/devices/" + item.id)
  //   refresh()
  // }

  // async function onClick(e, item) {
  //   app.device = item.item
  //   if (app.device.endpoints == null) app.device.endpoints = []
  //   if (app.device.allowedaddresses == null) app.device.allowedaddresses = []
  //   router.push("/provision")
  // }

  // async function scan() {
  //   // scandialog.value = false
  //   const ips = expandCidr(scancidr.value).toArray()
  //   for (var i = 0; i < ips.length; i++) {
  //     var ip = ips[i]
  //     var skip = false
  //     for (var item in items.value) {
  //       if (item.ip === ip) {
  //         skip = true
  //         break
  //       }
  //     }

  //     if (!skip) {
  //       scanning.value = ip
  //       progress.value = (i * 100) / ips.length
  //       const reachable = await device.present(ip)
  //       if (reachable) {
  //         const serial = await device.serial(ip)
  //         console.log('serial: ' + JSON.stringify(serial))
  //         var existing = await device.findBySerial(serial.serial)
  //         console.log('existing: ' + JSON.stringify(existing))

  //         const ipparts = ip.split('.')
  //         var newdevice = { name: '[unknown]', hostname: '[unknown]', serialno: serial.serial }
  //         if (existing && existing.id > 0) {
  //           newdevice = existing
  //         }

  //         newdevice.ip = ip
  //         newdevice.postfix = ipparts[2] + '.' + ipparts[3]
  //         const response = await device.save(newdevice)
  //         refresh()
  //       }
  //     }
  //   }

  //   progress.value = 100
  //   scandialog.value = false
  // }
</script>
