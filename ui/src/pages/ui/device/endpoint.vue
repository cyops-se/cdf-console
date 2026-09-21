<template>
    <v-card title="Endpoint provisioning" width="100%" height="100%">
        <template v-slot:append>
            <!-- <v-btn class="ml-2" @click="changeRoleToHub" color="primary">Change to hub</v-btn>
            <v-btn class="ml-2" @click="configure" color="primary">Configure</v-btn> -->
            <v-btn class="ml-2" icon="mdi-content-save" @click="save"></v-btn>
        </template>
        <v-card-text>
            <v-card>
                <v-card-text>
                    <v-row>
                        <v-col cols="12" sm="3"><v-text-field v-model="app.device.hostname"
                                label="Hostname"></v-text-field></v-col>
                        <v-col cols="12" sm="3"><v-text-field v-model="lancidr" label="Lan CIDR"
                                :base-color="app.device.lan && app.device.lanpresent ? 'success' : 'error'"></v-text-field></v-col>
                        <v-col cols="12" sm="3"><v-text-field v-model="wancidr" label="Wan CIDR"
                                :base-color="app.device.wan && app.device.wanpresent ? 'success' : 'error'"></v-text-field></v-col>
                        <v-col cols="12" sm="3"><v-text-field v-model="apn4gcidr" label="Apn 4G CIDR"
                                :base-color="app.device.apn4g && app.device.apn4gpresent ? 'success' : 'error'"></v-text-field></v-col>
                        <v-col cols="12" sm="3"><v-text-field v-model="syslog" label="Syslog IP"></v-text-field></v-col>
                        <v-col cols="12" sm="4"><v-text-field label="Public key" :model-value="app.device.pubkey"
                                readonly disabled></v-text-field></v-col>
                    </v-row>
                </v-card-text>
            </v-card>
            <AllowedList type='endpoint' />
        </v-card-text>
    </v-card>
</template>

<script setup>
    import axios from "axios"
    import { onMounted } from "vue"
    import { useRouter } from 'vue-router'
    import { useAppStore } from '@/stores/app'
    import { device } from '@/plugins/device'

    const router = useRouter()
    const app = useAppStore()
    const lancidr = ref('')
    const wancidr = ref('')
    const apn4gcidr = ref('')
    const syslog = ref('10.49.88.254 8514 udp')

    // Computed properties for interfaces
    const wgInterface = computed(() => device.getWireguardInterface(app.device, 'wg1'))

    // Computed property for device IP to use for operations
    const deviceIp = computed(() => app.device.preferredip || app.device.wan || app.device.lan || app.device.apn4g)

    onMounted(async () => {
        app.clearalert()
        refresh()
    })

    async function refresh() {
        console.log('endpoint key: ', JSON.stringify(app.device.pubkey))
        try {
            // Initialize CIDR values from interface fields
            if (app.device.lan && app.device.lanmask) {
                lancidr.value = device.getInterfaceCidr(app.device, 'lan') || ''
            }
            if (app.device.wan && app.device.wanmask) {
                wancidr.value = device.getInterfaceCidr(app.device, 'wan') || ''
            }
            if (app.device.apn4g && app.device.apn4gmask) {
                apn4gcidr.value = device.getInterfaceCidr(app.device, 'apn4g') || ''
            }

            // if (app.device.hubid > 0) app.device.hub = await axios.get('/api/data/devices/id/' + app.device.hubid).data
        } catch (e) {
            console.log('some problem with networks: ', e)
        }

        const lanPresent = app.device.lanpresent || false
        const wanPresent = app.device.wanpresent || false
        const apn4gPresent = app.device.apn4gpresent || false

        // if (lanPresent || wanPresent || apn4gPresent) {
        //     // Use preferred IP or first available interface IP
        //     const checkIp = app.device.preferredip || app.device.wan || app.device.lan
        //     if (checkIp) {
        //         var response = await device.wg(checkIp)
        //         if (response === null) {
        //             app.setalert('Wireguard not initialized on device!');
        //         } else {
        //             const wgIface = device.getWireguardInterface(app.device, 'wg1')
        //             const pubkey = wgIface ? wgIface.publickey : app.device.pubkey
        //             if (response['network.wg1.public_key'] !== pubkey) {
        //                 app.setalert('Wireguard key in database does not match device key!')
        //             }
        //         }
        //     }
        // }
    }

    async function setntp() {
        // const args = [{key: 'hubip', value: ''}, {key: 'hubwgip', value: ''}, {key: 'hubvxlanip', value: ''}, {key: 'vxlanid', value: ''}, {key: 'peerip', value: ''}, {key: 'hubpubkey', value: ''}]
        const wanIp = app.device.wan
        const args = [{ key: 'hubip', value: wanIp }]

        var response = await device.runscript2(wanIp, 'init_ntpserver.sh', [])
        if (!response.success) {
            console.log('test 1 response: ' + JSON.stringify(response))
        }

        await app.device.endpoints.forEach(async (endpoint) => {
            response = await device.runscript2(endpoint.wan, 'init_ntpclient.sh', args)
            if (!response.success) {
                console.log('test 2 response: ' + JSON.stringify(response))
            }
        })
    }

    async function save() {
        // Use helper methods to set interfaces from CIDR values
        if (lancidr.value) {
            device.setInterfaceFromCidr(app.device, 'lan', lancidr.value)
        }
        if (wancidr.value) {
            device.setInterfaceFromCidr(app.device, 'wan', wancidr.value)
        }
        if (apn4gcidr.value) {
            device.setInterfaceFromCidr(app.device, 'wan', wancidr.value)
        }

        delete app.device.lastcheck

        // Use WAN interface IP for syslog script
        const wanIp = app.device.wan
        await device.runscript(wanIp, 'init_syslog.sh', syslog.value);

        delete app.device.wanstats;
        delete app.device.lanstats;
        delete app.device.apn4gstats;
        delete app.device.hub?.wanstats;
        delete app.device.hub?.lanstats;
        delete app.device.hub?.apn4gstats;

        console.log('saving item: ' + JSON.stringify(app.device))
        var response = await axios.put("/api/data/devices", app.device)
        console.log('save response: ' + JSON.stringify(response.data))

        // Interfaces are now saved as part of the device object (embedded fields)
        // Wireguard interfaces are still separate and saved with the device if they exist
    }

    // async function configure() {
    //     if (app.device.role !== 'endpoint') return

    //     await save()
    //     await setntp()

    //     // Configure wireguard in hub and all peers
    //     // Start with removing wireguard peers in hub
    //     var response = await device.uciarray(deviceIp.value, 'uci show network | grep wireguard_wg1')
    //     console.log('finding all peers: ' + JSON.stringify(response))
    //     for (var element of response) {
    //         await device.uci(deviceIp.value, 'uci delete ' + element.key)
    //     }

    //     // need synchrounous iteration here
    //     for (var i = 0; i < app.device.endpoints.length; i++) {
    //         var endpoint = app.device.endpoints[i]
    //         // set up peer in hub
    //         const parts = endpoint.lan?.split('.')
    //         const peerid = parts[2] + parts[3] // endpoint.wan.replace('10.49.', '').replace('.', '')
    //         const peerip = endpoint.wan
    //         const peerpubkey = endpoint.pubkey
    //         const vxlanip = '10.48.' + parts[2] + '.' + parts[3] // endpoint.wan.replace('10.49.', '10.48.')
    //         const allowedip = vxlanip + '/32' // endpoint.wan.replace('10.49', '10.48') + '/32'
    //         delete endpoint.status
    //         console.log('HUB add_wg_peer.sh: ', peerid + ' ' + peerpubkey + ' ' + peerip + ' ' + allowedip)
    //         response = await device.runscript(deviceIp.value, 'add_wg_peer.sh', peerid + ' ' + peerpubkey + ' ' + peerip + ' ' + allowedip)
    //         if (!response.success) {
    //             console.log('response (after adding peer to hub): ' + JSON.stringify(response))
    //             endpoint.status = 'hub wg failed'
    //             continue
    //         }

    //         console.log('HUB add_vxlan_peer.sh: ', vxlanip + ' ' + peerid)
    //         response = await device.runscript(deviceIp.value, 'add_vxlan_peer.sh', vxlanip + ' ' + peerid)
    //         if (!response.success) {
    //             console.log('response (after adding vxlan to peer): ' + JSON.stringify(response))
    //             endpoint.status = 'hub vxlan failed'
    //             continue
    //         }
    //     }

    //     console.log('HUB init_wg_firewall.sh: ')
    //     response = await device.runscript(deviceIp.value, 'init_wg_firewall.sh', '')
    //     if (!response.success) {
    //         console.log('response (after setting firewall rules): ' + JSON.stringify(response))
    //         endpoint.status = 'hub firewall failed'
    //     }

    //     // we can run the provisioning of endpoints in parallell
    //     await app.device.endpoints.forEach(async (endpoint) => {
    //         // set up hub in peer

    //         if (endpoints.status) {
    //             // no need to configure endpoint when it couldn't be configured in the hub
    //             return
    //         }

    //         const hubid = 'wghub1'
    //         const hubip = app.device.wan
    //         const wgIface = device.getWireguardInterface(app.device, 'wg1')
    //         const hubpubkey = wgIface ? wgIface.publickey : app.device.pubkey

    //         const hubparts = hubip.split('.')
    //         const peerparts = endpoint.lan?.split('.')
    //         const vxlanid = peerparts[2] + peerparts[3] // endpoint.wan.replace('10.49.', '').replace('.', '')
    //         const vxlanip = '10.48.' + hubparts[2] + '.' + hubparts[3] // hubip.replace('10.49.', '10.48.')
    //         // const vxlanip = hubip.replace('10.49.', '10.48.')

    //         // await device.uci(deviceIp.value, 'uci -q delete network.wghub1')
    //         console.log('PEER add_wg_peer.sh: ', hubid + ' ' + hubpubkey + ' ' + hubip + ' 0.0.0.0/0')
    //         response = await device.runscript(endpoint.wan, 'add_wg_peer.sh', hubid + ' ' + hubpubkey + ' ' + hubip + ' 0.0.0.0/0')
    //         if (!response.success) {
    //             console.log('response (after adding hub to peer): ' + JSON.stringify(response))
    //             endpoint.status = 'peer wg failed'
    //             return
    //         }

    //         console.log('PEER init_vxlan.sh: ', vxlanip + ' ' + vxlanid)
    //         response = await device.runscript(endpoint.wan, 'init_vxlan.sh', vxlanip + ' ' + vxlanid)
    //         if (!response.success) {
    //             console.log('response (after adding vxlan to peer): ' + JSON.stringify(response))
    //             endpoint.status = 'peer vxlan failed'
    //             return
    //         }

    //         console.log('PEER init_wg_firewall.sh: ', vxlanip + ' ' + vxlanid)
    //         response = await device.runscript(endpoint.wan, 'init_wg_firewall.sh', vxlanip + ' ' + vxlanid)
    //         if (!response.success) {
    //             console.log('response (after setting firewall rules for peer): ' + JSON.stringify(response))
    //             endpoint.status = 'peer firewall failed'
    //         }
    //     })
    // }

    // async function changeRoleToHub() {
    //     app.device.role = 'hub'
    //     await save()
    //     router.push('/ui/devices')
    // }
</script>
