<template>
    <v-card width="100%" height="100%">
        <v-toolbar color="transparent">
            <v-toolbar-title class="text-h6" text="Devices"></v-toolbar-title>

            <v-text-field v-model="search" label="Search" prepend-inner-icon="mdi-magnify" variant="outlined"
                hide-details single-line></v-text-field>
            <template v-slot:append>
                <v-btn class="ml-4" color="primary" icon="mdi-refresh" @click="refresh"></v-btn>
                <v-btn class="ml-2" @click="scandialog = true">Scan</v-btn>
                <v-btn class="ml-2" color="success" icon="mdi-plus" @click="dialog = true"></v-btn>
            </template>
        </v-toolbar>
        <v-card-text>
            <v-data-table density="compact" :search="search" :headers="headers" :items="items" @click:row="onClick"
                :items-per-page="-1" :footer-props="{ itemsPerPageOptions: -1 }">
                <template v-slot:item.lan="{ item }">
                    <div class="text-end">
                        <v-chip v-if="item.lan" :color="item.lanpresent ? 'success' : 'error'" :text="item.lan"
                            class="text-uppercase" size="small" label></v-chip>
                        <span v-else class="text-grey">-</span>
                    </div>
                </template>
                <template v-slot:item.wan="{ item }">
                    <div class="text-end">
                        <v-chip v-if="item.wan" :color="item.wanpresent ? 'success' : 'error'" :text="item.wan"
                            class="text-uppercase" size="small" label></v-chip>
                        <span v-else class="text-grey">-</span>
                    </div>
                </template>
                <template v-slot:item.apn4g="{ item }">
                    <div class="text-end">
                        <v-chip v-if="item.apn4g" :color="item.apn4gpresent ? 'success' : 'error'" :text="item.apn4g"
                            class="text-uppercase" size="small" label></v-chip>
                        <span v-else class="text-grey">-</span>
                    </div>
                </template>
                <template v-slot:item.status="{ item }">
                    <div class="text-end">
                        <v-chip :color="item.status == 'reachable' ? 'success' : 'error'" :text="item.status || 'unknown'"
                            class="text-uppercase" size="small" label></v-chip>
                    </div>
                </template>
                <template v-slot:item.wgstatus="{ item }">
                    <div class="text-end">
                        <v-chip :color="item.wgstatus == 'active' ? 'success' : 'error'" :text="item.wgstatus"
                            class="text-uppercase" size="small" label></v-chip>
                    </div>
                </template>
                <template v-slot:item.sshstatus="{ item }">
                    <div class="text-end">
                        <v-chip :color="item.sshstatus == 'valid' ? 'success' : 'error'" :text="item.sshstatus"
                            class="text-uppercase" size="small" label></v-chip>
                    </div>
                </template>
                <template v-slot:item.actions="{ item }">
                    <!-- <v-icon class="me-2" size="small" @click="editItem(item)">mdi-pencil</v-icon> -->
                    <v-icon class="me-2" size="small" @click="configureItem(item)">mdi-pencil</v-icon>
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
                    <v-text-field v-model="pwd" label="Password" required></v-text-field>
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
    </v-dialog>
</template>

<script setup>
import axios from "axios"
import { onMounted } from "vue"
import { useRouter } from 'vue-router'
import { expandCidr, parseCidr } from "cidr-tools"
import { useAppStore } from '@/stores/app'
import { device } from '@/plugins/device'

const props = defineProps(['items'])

const router = useRouter()
const app = useAppStore()
const items = ref([])
const footer = ref({ 'items-per-page-options': -1 })
const dialog = ref(false)
const scandialog = ref(false)
// const headers = ref([{ title: 'Hostname', key: 'hostname', align: 'start' }, { title: 'Role', key: 'role', align: 'end' }, { title: 'Hub', key: 'hubid', align: 'end' }, { title: 'Serial', key: 'serialno', align: 'end' }, { title: 'WAN', key: 'wan', align: 'end' }, { title: 'APN 4G', key: 'apn4g', align: 'end' }, { title: 'LAN', key: 'lan', align: 'end' }, { title: 'Public key', key: 'pubkey', align: 'end' }, { title: 'Status', key: 'status', align: 'end' }, { title: 'Actions', key: 'actions', align: 'end' }])
const headers = ref([{ title: 'Hostname', key: 'hostname', align: 'start' },
{ title: 'Role', key: 'role', align: 'end' },
{ title: 'Hub', key: 'hubid', align: 'end' },
{ title: 'Serial', key: 'serialno', align: 'end' },
{ title: 'WAN', key: 'wan', align: 'end' },
{ title: 'APN', key: 'apn4g', align: 'end' },
{ title: 'LAN', key: 'lan', align: 'end' },
// { title: 'Reachable', key: 'status', align: 'end' },
{ title: 'Wireguard', key: 'wgstatus', align: 'end' },
{ title: 'SSH', key: 'sshstatus', align: 'end' },
{ title: 'Actions', key: 'actions', align: 'end' }])
const search = ref('')
const scancidr = ref('192.168.1.0/24')
const pwd = ref('admin01')
const scanning = ref('')
const progress = ref(0)

onMounted(async () => {
    refresh()
})

async function refresh() {
    var response = await device.alldevices()
    items.value = response
    await device.checklist(items.value)
    // console.log('devicelist: ' + JSON.stringify(items.value));
}

async function save() {
    dialog.value = false
    delete app.device.id
    const response = await device.save(app.device) // axios.put("/api/data/devices", device.value)
    // console.log('save device response: ' + JSON.stringify(response))
    refresh()
}

async function deleteItem(item) {
    const response = await axios.delete("/api/data/devices/" + item.id)
    // console.log('delete device response: ' + JSON.stringify(response))
    refresh()
}

async function onClick(e, item) {
    app.device = item.item
    if (!app.device.role) app.device.role = 'endpoint'
    if (app.device.endpoints == null) app.device.endpoints = []
    if (app.device.allowedaddresses == null) app.device.allowedaddresses = []
    // console.log('device selected: ' + JSON.stringify(app.device))
    if (app.device.role === 'endpoint') {
        // console.log('navigating to endpoint page')
        router.push('/ui/device/endpoint')
    } else if (app.device.role === 'hub') {
        // console.log('navigating to hub page')
        router.push('/ui/device/hub')
    } else {
        // console.log('navigating to general provisioning page')
        router.push('/ui/provision')
    }
}

async function scan() {
    // scandialog.value = false
    const ips = Array.from(expandCidr(scancidr.value))
    for (var i = 0; i < ips.length; i++) {
        var ip = ips[i]
        var skip = false
        for (var item in items.value) {
            if (item.ip === ip) {
                skip = true
                break
            }
        }

        if (!skip) {
            scanning.value = ip
            progress.value = (i * 100) / ips.length
            const reachable = await device.present(ip)
            if (reachable) {
                var newdevice = await device.importdevice(ip, '')
                newdevice.role = 'endpoint'
                const response = await device.save(newdevice)
                refresh()
            }
        }
    }

    progress.value = 100
    scandialog.value = false
}
</script>