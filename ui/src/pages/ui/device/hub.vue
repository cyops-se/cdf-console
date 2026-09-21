<template>
    <v-card title="Hub provisioning" width="100%" height="100%">
        <template v-slot:append>
            <!-- <v-btn class="ml-2" @click="changeRoleToEndpoint" color="primary">Change to endpoint</v-btn>
            <v-btn class="ml-2" @click="configure" color="primary">Configure all</v-btn> -->
            <v-btn class="ml-4" color="primary" icon="mdi-refresh" @click="refresh"></v-btn>
            <v-btn class="ml-2" @click="backup" color="secondary">Backup all</v-btn>
            <v-btn class="ml-2" icon="mdi-content-save" @click="save"></v-btn>
        </template>
        <v-card-text>
            <v-card>
                <v-row>
                    <v-col cols="12" sm="3"><v-text-field v-model="app.device.hostname"
                            label="Hostname"></v-text-field></v-col>
                    <v-col cols="12" sm="3"><v-text-field v-model="lancidr" label="Lan CIDR"
                            :base-color="app.device.lan && app.device.lanpresent ? 'success' : 'error'"></v-text-field></v-col>
                    <v-col cols="12" sm="3"><v-text-field v-model="wancidr" label="Wan CIDR"
                            :base-color="app.device.wan && app.device.wanpresent ? 'success' : 'error'"></v-text-field></v-col>
                    <v-col cols="12" sm="3"><v-text-field v-model="apncidr" label="Apn 4G CIDR"
                            :base-color="app.device.apn4g && app.device.apn4gpresent ? 'success' : 'error'"></v-text-field></v-col>
                </v-row>
                <v-row>
                    <v-col cols="12" sm="3"><v-text-field v-model="syslog" label="Syslog IP"></v-text-field></v-col>
                    <v-col cols="12" sm="4"><v-text-field label="Public key" :model-value="app.device.pubkey" readonly
                            disabled></v-text-field></v-col>
                </v-row>
            </v-card>

            <AllowedList type='Hub' />

            <v-card>
                <v-card-text>
                    <v-data-table :headers="headers" :items="app.device.endpoints" style="white-space: nowrap"
                        :items-per-page="-1" :footer-props="{ itemsPerPageOptions: -1 }">
                        <template v-slot:item.lan="{ item }">
                            <div class="text-end">
                                <v-chip :color="item.lanpresent ? 'success' : 'error'" :text="item.lan || ''"
                                    class="text-uppercase" size="small" label></v-chip>
                            </div>
                        </template>
                        <template v-slot:item.wan="{ item }">
                            <div class="text-end">
                                <v-chip :color="item.wanpresent ? 'success' : 'error'" :text="item.wan || ''"
                                    class="text-uppercase" size="small" label></v-chip>
                            </div>
                        </template>
                        <template v-slot:item.apn4g="{ item }">
                            <div class="text-end">
                                <v-chip :color="item.apn4gpresent ? 'success' : 'error'" :text="item.wan || ''"
                                    class="text-uppercase" size="small" label></v-chip>
                            </div>
                        </template>
                        <template v-slot:item.status="{ item }">
                            <div class="text-end">
                                <v-chip :color="item.status == 'reachable' ? 'success' : 'error'" :text="item.status"
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
                            <v-icon size="small" @click.stop="configureEndpoint(item)">mdi-wrench</v-icon>
                            <v-icon size="small" @click.stop="deleteEndpoint(item)" color="error">mdi-delete</v-icon>
                        </template>
                        <template v-slot:top>
                            <v-dialog v-model="dialog">
                                <template v-slot:activator="{ props }">
                                    <v-toolbar>
                                        <v-toolbar-title>Endpoints</v-toolbar-title>
                                        <v-spacer />
                                        <!-- <v-btn color="primary" text="Adopt endpoints" @click="adoptEndpoints"></v-btn> -->
                                        <v-btn color="primary" @click="dialog = !dialog">Add endpoints</v-btn>
                                    </v-toolbar>
                                </template>
                                <v-card title="Add endpoints">
                                    <v-card-text>
                                        <v-data-table v-model="selected" :headers="headers" :items="endpoints"
                                            show-select></v-data-table>
                                    </v-card-text>
                                    <v-card-actions>
                                        <v-btn text="Cancel" @click="dialog = false"></v-btn>
                                        <v-spacer></v-spacer>
                                        <v-btn color="success" text="Add endpoint" @click="addEndpoint"></v-btn>
                                    </v-card-actions>
                                </v-card>
                            </v-dialog>
                        </template>
                    </v-data-table>
                </v-card-text>
            </v-card>
        </v-card-text>
    </v-card>
    <v-dialog v-model="progressdialog" persistent max-width="800">
        <v-card title="Device Configuration Progress">
            <v-card-text>
                <v-list>
                    <v-list-item v-for="(step, index) in configSteps" :key="index">
                        <template v-slot:prepend>
                            <div class="d-flex align-center mr-3" style="width: 24px; height: 24px;">
                                <v-icon v-if="step.status === 'pending'" color="grey">mdi-circle-outline</v-icon>
                                <v-progress-circular v-else-if="step.status === 'running'" indeterminate size="24"
                                    width="3" color="primary"></v-progress-circular>
                                <v-icon v-else-if="step.status === 'success'" color="success">mdi-check-circle</v-icon>
                                <v-icon v-else-if="step.status === 'error'" color="error">mdi-alert-circle</v-icon>
                            </div>
                        </template>
                        <v-list-item-title>{{ step.name }}</v-list-item-title>
                        <v-list-item-subtitle v-if="step.error" class="text-error">{{ step.error
                        }}</v-list-item-subtitle>
                    </v-list-item>
                </v-list>
                <v-alert v-if="configComplete" :type="configSuccess ? 'success' : 'warning'" class="mt-4">
                    <strong v-if="configSuccess">Configuration completed successfully!</strong>
                    <strong v-else>Configuration completed with errors.</strong>
                    <div class="mt-2">Please review the steps above for details.</div>
                </v-alert>
            </v-card-text>
            <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn :disabled="!configComplete" color="primary" @click="closeProgressDialog">Close</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup>
    import axios from "axios"
    import { onMounted } from "vue"
    import { useRouter } from 'vue-router'
    import { useAppStore } from '@/stores/app'
    import { device } from '@/plugins/device'

    const router = useRouter()
    const app = useAppStore()
    const dialog = ref(false)
    const progressdialog = ref(false)
    const progress = ref(0)
    const present = ref(false)
    const previp = ref('')
    const prevmask = ref('')
    const endpoints = ref([])
    const selected = ref([])
    const ip = ref('')
    const lancidr = ref('')
    const wancidr = ref('')
    const apncidr = ref('')
    const wgcidr = ref('')
    const newdevice = ref('')
    const syslog = ref('10.49.88.254 8514 udp')
    const configSteps = ref([])
    const configComplete = ref(false)
    const configSuccess = ref(false)

    // Computed properties for interfaces
    const wgInterface = computed(() => device.getWireguardInterface(app.device, 'wg1'))

    // Computed property for device IP to use for operations
    const deviceIp = computed(() => app.device.preferredip || app.device.wan || app.device.lan || app.device.apn4g)
    const headers = ref([{ title: 'Hostname', key: 'hostname', align: 'start' },
    { title: 'Serial', key: 'serialno', align: 'end' },
    { title: 'WAN', key: 'wan', align: 'end' },
    { title: 'LAN', key: 'lan', align: 'end' },
    { title: 'APN4G', key: 'apn4g', align: 'end' },
    // { title: 'Reachable', key: 'status', align: 'end' },
    { title: 'Wireguard', key: 'wgstatus', align: 'end' },
    { title: 'SSH', key: 'sshstatus', align: 'end' },
    { title: 'Last check', key: 'lastcheck', align: 'end' },
    { title: 'Public key', key: 'pubkey', align: 'end' },
    { title: 'Actions', key: 'actions', align: 'end' }])

    onMounted(async () => {
        app.clearalert()
        refresh()
    })

    async function refresh() {
        // console.log('hub: ', JSON.stringify(app.device))

        // Initialize CIDR values from interface fields
        lancidr.value = device.getInterfaceCidr(app.device, 'lan') || ''
        wancidr.value = device.getInterfaceCidr(app.device, 'wan') || ''
        apncidr.value = device.getInterfaceCidr(app.device, 'apn4g') || ''

        app.device.endpoints = await device.hubendpoints(app.device)
        await device.checklist(app.device.endpoints)

        var response = await device.allendpoints()
        endpoints.value = response
    }

    async function addEndpoint() {
        // move selected items from endpoints.value to app.device.endpoints to avoid them being added more than once
        for (var i = 0; i < selected.value.length; i++) {
            for (var j = endpoints.value.length - 1; j >= 0; j--) {
                if (endpoints.value[j].id == selected.value[i]) {
                    var item = endpoints.value.splice(j, 1)
                    app.device.endpoints.push(item[0])
                    break
                }
            }
        }

        selected.value = []
        dialog.value = false
    }

    async function deleteEndpoint(item) {
        // move selected items from app.device.endpoints to endpoints.value to avoid duplicates
        for (var i = 0; i < app.device.endpoints.length; i++) {
            if (app.device.endpoints[i].id == item.id) {
                await device.delete(app.device.endpoints[i]);
                break
            }
        }

        dialog.value = false
        this.refresh();
    }

    // Helper function to update step status
    function updateStepStatus(stepIndex, status, error = null) {
        if (stepIndex < configSteps.value.length) {
            configSteps.value[stepIndex].status = status
            if (error) {
                configSteps.value[stepIndex].error = error
            }
        }
    }

    // Step 1: Save configuration and setup syslog
    async function stepSaveConfiguration(endpoint) {
        const stepIndex = 0
        updateStepStatus(stepIndex, 'running')
        try {
            await save()
            await setntp(endpoint)
            updateStepStatus(stepIndex, 'success')
            return true
        } catch (error) {
            updateStepStatus(stepIndex, 'error', `Failed to save configuration: ${error.message}`)
            return false
        }
    }

    // Step 2: Add WireGuard peer on hub
    async function stepAddWireguardPeerOnHub(endpoint, config) {
        const stepIndex = 1
        updateStepStatus(stepIndex, 'running')
        try {
            const { peerid, peerpubkey, peerip, allowedip } = config
            // console.log('HUB add_wg_peer.sh: ', peerid, peerpubkey, peerip, allowedip)
            const response = await device.runscript(
                deviceIp.value,
                'add_wg_peer.sh',
                `${peerid} ${peerpubkey} ${peerip} ${allowedip}`
            )
            if (!response.success) {
                const errorMsg = response.error || 'Failed to add WireGuard peer on hub'
                // console.log('response (after adding peer to hub): ' + JSON.stringify(response))
                endpoint.status = 'hub wg failed'
                updateStepStatus(stepIndex, 'error', errorMsg)
                return false
            }
            updateStepStatus(stepIndex, 'success')
            return true
        } catch (error) {
            updateStepStatus(stepIndex, 'error', `Exception: ${error.message}`)
            return false
        }
    }

    // Step 3: Add VXLAN peer on hub
    async function stepAddVxlanPeerOnHub(endpoint, config) {
        const stepIndex = 2
        updateStepStatus(stepIndex, 'running')
        try {
            const { vxlanip, peerid } = config
            // console.log('HUB add_vxlan_peer.sh: ', vxlanip, peerid)
            const response = await device.runscript(
                deviceIp.value,
                'add_vxlan_peer.sh',
                `${vxlanip} ${peerid}`
            )
            if (!response.success) {
                const errorMsg = response.error || 'Failed to add VXLAN peer on hub'
                // console.log('response (after adding vxlan to peer): ' + JSON.stringify(response))
                endpoint.status = 'hub vxlan failed'
                updateStepStatus(stepIndex, 'error', errorMsg)
                return false
            }
            updateStepStatus(stepIndex, 'success')
            return true
        } catch (error) {
            updateStepStatus(stepIndex, 'error', `Exception: ${error.message}`)
            return false
        }
    }

    // Step 4: Initialize WireGuard firewall on hub
    async function stepInitWireguardFirewallOnHub(endpoint) {
        const stepIndex = 3
        updateStepStatus(stepIndex, 'running')
        try {
            // console.log('HUB init_wg_firewall.sh')
            const response = await device.runscript(deviceIp.value, 'init_wg_firewall.sh', '')
            if (!response.success) {
                const errorMsg = response.error || 'Failed to initialize WireGuard firewall on hub'
                // console.log('response (after setting firewall rules): ' + JSON.stringify(response))
                endpoint.status = 'hub firewall failed'
                updateStepStatus(stepIndex, 'error', errorMsg)
                return false
            }
            updateStepStatus(stepIndex, 'success')
            return true
        } catch (error) {
            updateStepStatus(stepIndex, 'error', `Exception: ${error.message}`)
            return false
        }
    }

    // Step 5: Add WireGuard peer on endpoint
    async function stepAddWireguardPeerOnEndpoint(endpoint, config) {
        const stepIndex = 4
        updateStepStatus(stepIndex, 'running')
        try {
            const { hubid, hubpubkey, hubip } = config
            // console.log('PEER add_wg_peer.sh: ', hubid, hubpubkey, hubip, '0.0.0.0/0')
            const response = await device.runscript(
                endpoint.wan,
                'add_wg_peer.sh',
                `${hubid} ${hubpubkey} ${hubip} 0.0.0.0/0`
            )
            if (!response.success) {
                const errorMsg = response.error || 'Failed to add WireGuard peer on endpoint'
                // console.log('response (after adding hub to peer): ' + JSON.stringify(response))
                endpoint.status = 'peer wg failed'
                updateStepStatus(stepIndex, 'error', errorMsg)
                return false
            }
            updateStepStatus(stepIndex, 'success')
            return true
        } catch (error) {
            updateStepStatus(stepIndex, 'error', `Exception: ${error.message}`)
            return false
        }
    }

    // Step 6: Initialize VXLAN on endpoint
    async function stepInitVxlanOnEndpoint(endpoint, config) {
        const stepIndex = 5
        updateStepStatus(stepIndex, 'running')
        try {
            const { vxlanhubip, vxlanid } = config
            // console.log('PEER init_vxlan.sh: ', vxlanhubip, vxlanid)
            const response = await device.runscript(
                endpoint.wan,
                'init_vxlan.sh',
                `${vxlanhubip} ${vxlanid}`
            )
            if (!response.success) {
                const errorMsg = response.error || 'Failed to initialize VXLAN on endpoint'
                // console.log('response (after adding vxlan to peer): ' + JSON.stringify(response))
                endpoint.status = 'peer vxlan failed'
                updateStepStatus(stepIndex, 'error', errorMsg)
                return false
            }
            updateStepStatus(stepIndex, 'success')
            return true
        } catch (error) {
            updateStepStatus(stepIndex, 'error', `Exception: ${error.message}`)
            return false
        }
    }

    // Step 7: Initialize WireGuard firewall on endpoint
    async function stepInitWireguardFirewallOnEndpoint(endpoint, config) {
        const stepIndex = 6
        updateStepStatus(stepIndex, 'running')
        try {
            const { vxlanip, vxlanid } = config
            // console.log('PEER init_wg_firewall.sh: ', vxlanip, vxlanid)
            const response = await device.runscript(
                endpoint.wan,
                'init_wg_firewall.sh',
                `${vxlanip} ${vxlanid}`
            )
            if (!response.success) {
                const errorMsg = response.error || 'Failed to initialize WireGuard firewall on endpoint'
                // console.log('response (after setting firewall rules for peer): ' + JSON.stringify(response))
                endpoint.status = 'peer firewall failed'
                updateStepStatus(stepIndex, 'error', errorMsg)
                return false
            }
            updateStepStatus(stepIndex, 'success')
            return true
        } catch (error) {
            updateStepStatus(stepIndex, 'error', `Exception: ${error.message}`)
            return false
        }
    }

    // Step 8: Refresh device data
    async function stepRefreshDeviceData() {
        const stepIndex = 7
        updateStepStatus(stepIndex, 'running')
        try {
            await refresh()
            updateStepStatus(stepIndex, 'success')
            return true
        } catch (error) {
            updateStepStatus(stepIndex, 'error', `Failed to refresh: ${error.message}`)
            return false
        }
    }

    // Main configuration function
    async function configureEndpoint(item) {
        // console.log('configuring device: ' + JSON.stringify(item))

        // Initialize progress dialog with steps
        configSteps.value = [
            { name: 'Save configuration and setup NTP/Syslog', status: 'pending', error: null },
            { name: 'Add WireGuard peer on hub', status: 'pending', error: null },
            { name: 'Add VXLAN peer on hub', status: 'pending', error: null },
            { name: 'Initialize WireGuard firewall on hub', status: 'pending', error: null },
            { name: 'Add WireGuard peer on endpoint', status: 'pending', error: null },
            { name: 'Initialize VXLAN on endpoint', status: 'pending', error: null },
            { name: 'Initialize WireGuard firewall on endpoint', status: 'pending', error: null },
            { name: 'Refresh device data', status: 'pending', error: null }
        ]

        configComplete.value = false
        configSuccess.value = false
        progressdialog.value = true

        const endpoint = item
        delete endpoint.status

        // Calculate configuration parameters
        const parts = endpoint.lan?.split('.')
        const peerid = parts[2] + parts[3]
        const peerip = endpoint.wan
        const peerpubkey = endpoint.pubkey
        const vxlanip = '10.48.' + parts[2] + '.' + parts[3]
        const allowedip = vxlanip + '/32'

        const hubid = 'wghub1'
        const hubip = app.device.wan
        const wgIface = device.getWireguardInterface(app.device, 'wg1')
        const hubpubkey = wgIface ? wgIface.publickey : app.device.pubkey

        const peerparts = endpoint.lan?.split('.')
        const vxlanid = peerparts[2] + peerparts[3] // same as peerid?
        const hubipparts = hubip.split('.');
        // const vxlanhubip = hubip.replace('10.49.', '10.48.')
        const vxlanhubip = '10.48.' + hubipparts[2] + '.' + hubipparts[3]

        const hubConfig = { peerid, peerpubkey, peerip, allowedip, vxlanip }
        const endpointConfig = { hubid, hubpubkey, hubip, vxlanhubip, vxlanid, vxlanip }

        // Execute configuration steps - continue even if errors occur
        let success = true

        // Step 1: Save configuration
        const step1Success = await stepSaveConfiguration(endpoint)
        success = success && step1Success

        // Step 2: Add WireGuard peer on hub
        const step2Success = await stepAddWireguardPeerOnHub(endpoint, hubConfig)
        success = success && step2Success

        // Step 3: Add VXLAN peer on hub
        const step3Success = await stepAddVxlanPeerOnHub(endpoint, hubConfig)
        success = success && step3Success

        // Step 4: Initialize WireGuard firewall on hub
        const step4Success = await stepInitWireguardFirewallOnHub(endpoint)
        success = success && step4Success

        // Step 5: Add WireGuard peer on endpoint
        const step5Success = await stepAddWireguardPeerOnEndpoint(endpoint, endpointConfig)
        success = success && step5Success

        // Step 6: Initialize VXLAN on endpoint
        const step6Success = await stepInitVxlanOnEndpoint(endpoint, endpointConfig)
        success = success && step6Success

        // Step 7: Initialize WireGuard firewall on endpoint
        const step7Success = await stepInitWireguardFirewallOnEndpoint(endpoint, endpointConfig)
        success = success && step7Success

        // Step 8: Refresh device data
        await stepRefreshDeviceData()

        // Mark configuration as complete
        configComplete.value = true
        configSuccess.value = success
    }

    function closeProgressDialog() {
        progressdialog.value = false
        configSteps.value = []
        configComplete.value = false
        configSuccess.value = false
    }

    async function setntp(endpoint) {
        // const args = [{key: 'hubip', value: ''}, {key: 'hubwgip', value: ''}, {key: 'hubvxlanip', value: ''}, {key: 'vxlanid', value: ''}, {key: 'peerip', value: ''}, {key: 'hubpubkey', value: ''}]
        const wanIp = app.device.wan
        const args = [{ key: 'hubip', value: wanIp }]

        var response = await device.runscript2(wanIp, 'init_ntpserver.sh', [])
        if (!response.success) {
            // console.log('test 1 response: ' + JSON.stringify(response))
        }

        response = await device.runscript2(endpoint.preferredip, 'init_ntpclient.sh', args)
        // await app.device.endpoints.forEach(async (endpoint) => {
        //     response = await device.runscript2(endpoint.wan, 'init_ntpclient.sh', args)
        //     if (!response.success) {
        //         // console.log('test 2 response: ' + JSON.stringify(response))
        //     }
        // })
    }

    async function save() {
        // Use helper methods to set interfaces from CIDR values
        if (lancidr.value) {
            device.setInterfaceFromCidr(app.device, 'lan', lancidr.value)
        }
        if (wancidr.value) {
            device.setInterfaceFromCidr(app.device, 'wan', wancidr.value)
        }
        if (apncidr.value) {
            device.setInterfaceFromCidr(app.device, 'apn4g', apncidr.value)
        }

        delete app.device.lastcheck
        for (var i = 0; i < app.device.endpoints.length; i++) {
            delete app.device.endpoints[i].lastcheck
        }

        // Use WAN interface IP for syslog script
        const wanIp = app.device.waniface
        await device.runscript(wanIp, 'init_syslog.sh', syslog.value);

        // console.log('saving item: ' + JSON.stringify(app.device))
        var response = await axios.put("/api/data/devices", app.device)
        // console.log('save response: ' + JSON.stringify(response.data))
    }

    async function backup() {
        const response = await device.backup(app.device.id)
        // console.log('backup complete: ', response)
    }
</script>
