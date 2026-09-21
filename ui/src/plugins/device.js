import axios from "axios"
import { Netmask } from 'netmask'
import { parseCidr } from 'cidr-tools'
import { useAppStore } from '@/stores/app'

const app = useAppStore()

const device = {

    uci: async (ip, cmd) => {
        var result = {}
        var response = await axios.post('/api/device/shell/' + ip, cmd)

        if (response.data.success) {
            var lines = response.data.commands[0].output.split('\n')
            lines.forEach(element => {
                if (element != '') {
                    element = element.replace('=', ':')
                    var parts = element.split(':')
                    result[parts[0]] = parts[1].replaceAll('\'', '')
                }
            });
        }

        return result.uci ? null : result
    },

    uciarray: async (ip, cmd) => {
        var result = []
        var response = await axios.post('/api/device/shell/' + ip, cmd)

        if (response.data.success) {
            var lines = response.data.commands[0].output.split('\n')
            lines.forEach(element => {
                if (element != '') {
                    element = element.replace('=', ':')
                    var parts = element.split(':')
                    result.push({ key: parts[0], value: parts[1].replaceAll('\'', '') })
                }
            });
        }

        return result
    },

    shell: async (ip, cmd) => {
        var response = await axios.post('/api/device/shell/' + ip, cmd)
        return response
    },

    runscript: async (ip, script, ...args) => {
        const request = { script: script, args: ''.concat(' ', ...args) }
        var response = await axios.post('/api/device/script/' + ip, request)
        return response.data
    },

    runscript2: async (ip, script, args) => {
        const request = { script: script, args: args }
        var response = await axios.post('/api/device/script2/' + ip, request)
        return response.data
    },

    serial: async (ip) => {
        var result = { model: '', serial: '' }
        var response = await axios.get('/api/device/serial/' + ip)
        // console.log("get serial response: " + JSON.stringify(response.data))
        result.success = response.data.success
        if (response.data.success) {
            var parts = response.data.commands[0].output.split('\n')
            result.serial = parts[0]
            // result.model = parts[0]
        }

        return result
    },

    wgOld: async (ip) => {
        var result = {}
        var response = await axios.post('/api/device/shell/' + ip, 'uci show network.wg1')
        // console.log('wg response: ', JSON.stringify(response.data))

        if (response.data.success) {
            var lines = response.data.commands[0].output.split('\n')
            lines.forEach(element => {
                if (element != '') {
                    element = element.replace('=', ':')
                    var parts = element.split(':')
                    result[parts[0]] = parts[1].replaceAll('\'', '')
                }
            });
        }

        return result.uci ? null : result
    },

    wg: async (ip) => {
        const proto = await device.uciarray(ip, 'uci show network | grep wireguard | grep proto')
        if (proto.length == 0) return {}
        // console.log('proto: ' + JSON.stringify(proto))
        const wgiface = proto[0].key.split('.')[1]
        var response = await device.uci(ip, 'uci show network.' + wgiface)
        return response
    },

    wgpubkey: async (ip) => {
        const proto = await device.uciarray(ip, 'uci show network | grep wireguard | grep proto')
        const wgiface = proto[0].key.split('.')[1]
        var response = await device.uci(ip, 'uci show network.' + wgiface)
        const pubkey = response['network.' + wgiface + '.public_key']
        return pubkey
    },

    initwg: async (ip, ...args) => {
        const request = { script: 'init_wg.sh', args: ''.concat(...args) }
        var response = await axios.post('/api/device/script/' + ip, request)

        return response.data
    },

    checkSshKeys: async (ip) => {
        const response = await axios.get("/api/device/serial/" + ip)
        return response.data.success
    },

    initSsh: async (ip, pwd, askCb) => {
        const body = { username: 'admin', password: pwd }
        var response = await axios.post("/api/device/initssh/" + ip, body)
        return response.data.success
    },

    backup: async (id) => {
        var response = await axios.get("/api/device/backup/" + id)
        return response.data.success
    },

    // sysinfo: async (ip) => {
    //     const system = await device.uci(importip, 'uci show system.system')
    //     console.log('system: ' + JSON.stringify(system))

    //     var d = {
    //         hostname: system['system.system.hostname'], name: system['system.system.devicename'],
    //         model: system['system.system.device_code'], firmware: system['system.system.device_fw_version']
    //     }
    //     return d
    // },

    present: async (ip) => {
        if (!ip) return false
        var response = await axios.get('/api/device/present/' + ip)
        return response.data.success
    },

    checklist: async (items) => {
        for (var dev of items) {
            dev.useip = dev.wanpresent ? dev.wan : (dev.lanpresent ? dev.lan : dev.apn4g);
            dev.usemask = dev.wanpresent ? dev.wanmask : (dev.lanpresent ? dev.lanmask : dev.apn4gmask);
            dev.lastcheck = dev.lastcheck.substring(0, 19).replace('T', ' ')
            if (!dev.wgstatus) dev.wgstatus = 'unknown'
            if (!dev.sshstatus) dev.sshstatus = 'unknown'
        }
    },

    importdevice: async (importip, importpwd, progressCb) => {
        const present = await device.present(importip)
        if (!present) {
            app.showsnack('Import failed, IP address not reachable: ' + importip, 'error')
            return
        }

        const keysok = await device.checkSshKeys(importip)
        if (!keysok) {
            const pwdok = await device.initSsh(importip, importpwd)
            if (!pwdok) {
                app.showsnack('Import failed, wrong password: ' + importpwd, 'error')
                return
            }
        }

        var dev = { defaultpwd: importpwd, password: importpwd }
        const serial = await device.serial(importip)
        const existing = await device.findBySerial(serial.serial)
        if (existing && !existing.length) {
            dev = existing
        }

        dev.serialno = serial.serial
        // console.log('import device serial: ' + dev.serialno)

        // Get system info
        const system = await device.uci(importip, 'uci show system.system')

        dev.hostname = system['system.system.hostname']
        dev.name = system['system.system.devicename']
        dev.model = system['system.system.device_code']
        dev.firmware = system['system.system.device_fw_version']

        // Get regular network info
        const network = await device.uci(importip, 'uci show network')

        // Get modem info
        const modemcmd = await device.shell(importip, 'gsmctl -j')
        // console.log('modemconnected: ' + JSON.stringify(modemcmd.data))
        if (modemcmd.data.success) {
            if (modemcmd.data.commands[0].output.startsWith('Connected')) {
                // console.log('There is a MODEM connected')
                const modemcidrcmd = await device.shell(importip, `ip -f inet addr show qmimux0 | awk '/inet / {print $2}'`)
                // console.log('MODEM CIDR response: ' + JSON.stringify(modemcidrcmd.data))
                if (modemcidrcmd.data.success) {
                    const modemcidr = parseCidr(modemcidrcmd.data.commands[0].output.replace('\n', ''))
                    // console.log('MODEM CIDR ip: ', modemcidr.ip)
                    dev.apn4g = modemcidr.ip
                    dev.apn4gmask = '255.255.255.255'
                }
            }
        }

        // Get wireguard public key
        const pubkey = await device.wgpubkey(importip)
        dev.pubkey = pubkey
        return dev
    },

    importendpoint: async (importip, importpwd, progressCb) => {
        var endpoint = await device.importdevice(importip, importpwd, progressCb)
        if (!endpoint) return

        endpoint.role = 'endpoint'
        endpoint.state = 'adopted'

        return endpoint
    },

    importhub: async (importip, importpwd, progressCb) => {
        var hub = await device.importdevice(importip, importpwd, progressCb)
        if (!hub) return

        if (!hub.endpoints) hub.endpoints = []
        hub.role = 'hub'
        hub.state = 'adopted'

        var peers = await device.uciarray(importip, 'uci show network | grep endpoint')
        // console.log('wireguard peers: ' + JSON.stringify(peers))
        for (var kv of peers) {
            var endpoint = await device.importendpoint(kv.value, importpwd)
            if (!endpoint) {
                endpoint = {
                    state: 'pending',
                    status: 'not reachable',
                    role: 'endpoint'
                }
                // console.log('not able to reach endpoint, initializing a pending placeholder')
            }

            if (endpoint) {
                hub.endpoints.push(endpoint)
                // console.log('endpoint added to hub: ', JSON.stringify(endpoint))
            }
        }

        // console.log('imported hub: ' + JSON.stringify(hub))
        return hub
    },

    create: async (item) => {
        var response = await axios.post('/api/data/devices', item)
        return response.data
    },

    save: async (item) => {
        delete item.lastcheck
        var response = await axios.put('/api/data/devices', item)
        return response.data
    },

    // Each endpoint has a footprint in the hub that includes:
    // - ebtables rules
    // - vxlan[postfix] interface (uci delete network.vxlan[postfix], uci del_list network.br_lan.ports='vxlan[postfix]' )
    // - wireguard peer (uci delete network.[postfix])
    deleteFromHub: async (endpoint) => {
        try {
            if (!endpoint.hubid || endpoint.hubid == 0) return;
            var response = await device.findById(endpoint.hubid);
            // console.log('deleteFromHub response: ' + JSON.stringify(response));
            if (!response) return;

            var hub = response.data;
        } catch (e) {
            return { error: e.message }
        }
    },

    delete: async (item) => {
        try {
            if (item.role == 'endpoint') await device.deleteFromHub(item);
            var response = await axios.delete('/api/data/devices/' + item.id);
            return response.data
        } catch (e) {
            return { error: e.message }
        }
    },

    // api.Get("/data/:type/id/:id", GetDataByID)
    findById: async (id) => {
        try {
            var response = await axios.get('/api/data/devices/id/' + id);
            return response.data;
        } catch (e) {
            return { error: e.message }
        }
    },

    findBySerial: async (serial) => {
        var response = await axios.get('/api/data/devices/field/serial_no/' + serial)
        return response.data && response.data.length === 0 ? undefined : response.data[0]
    },

    hubendpoints: async (hub) => {
        var response = await axios.get('/api/data/devices/field/role/endpoint')
        const endpoints = response.data.filter((e) => e.hubid === hub.id)
        return endpoints
    },

    allendpoints: async () => {
        var response = await axios.get('/api/data/devices/field/role/endpoint')
        const available = response.data.filter((e) => e.hubid == null)
        return available
    },

    alldevices: async () => {
        var response = await axios.get('/api/data/devices')
        return response.data
    },

    resetstats: async () => {
        var response = await axios.post('/api/device/resetstats')
        return response.data
    },

    getWireguardInterface: (device, name) => {
        if (device.wireguardinterfaces && device.wireguardinterfaces.length > 0) {
            return device.wireguardinterfaces.find(iface => iface.name === name)
        }
        // Fallback to legacy field
        if (name === 'wg1' && device.pubkey) {
            return { name: 'wg1', publickey: device.pubkey, status: device.wgstatus }
        }
        return null
    },

    // Get interface IP and mask as CIDR
    getInterfaceCidr: (device, name) => {
        // Get the interface directly
        let ip = '';
        let mask = '';
        if (name === 'wan') {
            ip = device.wan;
            mask = device.wanmask;
        } else if (name === 'lan') {
            ip = device.lan;
            mask = device.lanmask;
        } else if (name === 'apn' || name === 'apn4g') {
            ip = device.apn4g;
            mask = device.apn4gmask;
        }

        try {
            const block = new Netmask(ip, mask)
            return ip + '/' + block.bitmask
        } catch (e) {
            return ''
        }
    },

    // Set interface from CIDR
    setInterfaceFromCidr: (device, name, cidr) => {
        if (!cidr) return

        try {
            const parsed = parseCidr(cidr)
            const block = new Netmask(cidr)

            // Set interface directly
            const emptyInterface = { ipaddress: '', netmask: '', macaddress: '', present: false, stats: {} }
            let iface = null
            if (name === 'wan') {
                device.wan = parsed.ip;
                device.wanmask = block.mask;
            } else if (name === 'lan') {
                device.lan = parsed.ip;
                device.lanmask = block.mask;
            } else if (name === 'apn' || name === 'apn4g') {
                device.apn4g = parsed.ip;
                device.apn4gmask = block.mask;
            }
        } catch (e) {
            console.error('Invalid CIDR:', cidr, e)
        }
    },

    preferredip: (item) => {
        var ip = item.preferredip;
        if (item.wanpresent) ip = item.wan;
        else if (item.lanpresent) ip = item.lan;
        else if (item.apn4gpresent) ip = item.apn4g;
        return ip;
    }
}

export { device }