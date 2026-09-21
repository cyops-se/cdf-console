<template>
  <v-card style="width: 100%">
    <v-card title="Access settings">
      <v-card-text>
        <v-row>
          <v-col cols="12" sm="2"><v-text-field v-model="ip"
              label="Access IP address to onboard"></v-text-field></v-col>
          <v-col cols="12" sm="2"><v-text-field v-model="pwd" label="Access password"></v-text-field></v-col>
          <v-col cols="12" sm="2"><v-alert :color="access == 'ok' ? 'success' : 'error'" class="text-uppercase" label>
              <div class="text-center">{{ access }}</div>
            </v-alert></v-col>
          <v-col cols="12" sm="2">
            <v-btn @click="testAccess" :disabled="disableBtn" width="100%" height="56">Test access</v-btn>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
    <v-card title="Device info" class="mt-4" v-if="deviceInfo.model">
      <v-card-text>
        <v-row>
          <v-col cols="12" sm="2"><v-text-field v-model="deviceInfo.model" label="Model"
              readonly></v-text-field></v-col>
          <v-col cols="12" sm="2"><v-text-field v-model="deviceInfo.firmware" label="Current Firmware"
              readonly></v-text-field></v-col>
          <v-col cols="12" sm="2"><v-text-field v-model="deviceInfo.serialno" label="Serial Number"
              readonly></v-text-field></v-col>
          <v-col cols="12" sm="2"><v-select v-model="selectedFirmware" label="Target Firmware"
              :items="availableFirmware" item-title="version" item-value="version"
              :disabled="availableFirmware.length === 0"></v-select></v-col>
        </v-row>
      </v-card-text>
    </v-card>
    <v-card title="Device settings" class="mt-4">
      <v-card-text>
        <v-row>
          <v-col cols="12" sm="2"><v-text-field v-model="hostname" label="Hostname"></v-text-field></v-col>
          <v-col cols="12" sm="2"><v-text-field v-model="newpwd" label="New password"></v-text-field></v-col>
          <v-col cols="12" sm="2"><v-select v-model="role" label="Role" :items="['Endpoint', 'Hub']"></v-select></v-col>
          <v-col cols="12" sm="2"><!-- v-text-field v-model="postfix" label="Postfix for new address"
              @change="postfixChanged"></v-text-field --></v-col>
        </v-row>
        <v-row>
          <v-col cols="12" sm="2"><v-text-field v-model="lancidr" label="Lan CIDR"
              :error="!isValidCidr(lancidr)"></v-text-field></v-col>
          <v-col cols="12" sm="2"><v-text-field v-model="wancidr" label="Wan CIDR"
              :error="!isValidCidr(wancidr)"></v-text-field></v-col>
          <v-col cols="12" sm="2"><v-text-field v-model="apn4gcidr" label="Apn4G CIDR"
              disabled></v-text-field></v-col>
          <v-col cols="12" sm="2">
            <v-btn @click="startOnboarding" :disabled="!areCidrsValid" width="100%" height="56" color="primary">Start
              onboarding</v-btn>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
    <v-card title="Onboarding progress" class="mt-4">
      <v-card-text>
        <v-expansion-panels>
          <v-expansion-panel v-for="(item, index) in flow" :key="index" :disabled="item.disabled">
            <v-expansion-panel-title>
              {{ item.name }} {{ item.state }}
              <template v-slot:actions>
                <v-icon :color="item.color" :icon="item.icon"></v-icon>
              </template>
            </v-expansion-panel-title>
            <v-expansion-panel-text>
              {{ item.text }}
            </v-expansion-panel-text>
          </v-expansion-panel>
        </v-expansion-panels>
      </v-card-text>
    </v-card>
  </v-card>
</template>

<script setup>
  import axios from "axios"
  import { onMounted, computed } from "vue"
  import { Netmask } from 'netmask'
  import { parseCidr } from 'cidr-tools'
  import { device } from '@/plugins/device'

  const access = ref('unknown')
  const ip = ref('192.168.1.1')
  const pwd = ref('admin01')
  const role = ref('Endpoint')
  const hostname = ref('')
  const newpwd = ref('')
  const lancidr = ref('192.168.1.1/24')
  const wancidr = ref('10.49.1.1/24')
  const apn4gcidr = ref('')
  const disableBtn = ref(false)
  const deviceInfo = ref({ model: '', firmware: '', serialno: '', hostname: '' })
  const allFirmwareVersions = ref([])
  const availableFirmware = ref([])
  const selectedFirmware = ref('')

  BigInt.prototype.toJSON = function () {
    return { $bigint: this.toString() };
  };

  // Function to validate CIDR notation
  function isValidCidr(cidrString) {
    if (!cidrString || cidrString.trim() === '') {
      return false
    }
    try {
      const parsed = parseCidr(cidrString.trim())
      // parseCidr returns an object with ip, prefix, and other properties
      // Check that we got valid values
      const valid = parsed && parsed.ip && parsed.prefixPresent
      return valid
    } catch (error) {
      return false
    }
  }

  // Computed property to check if both CIDRs are valid
  const areCidrsValid = computed(() => {
    return isValidCidr(lancidr.value) && isValidCidr(wancidr.value)
  })

  const flow = ref([
    { name: 'Check availability', action: checkAccessibility },
    { name: 'Check SSH keys', action: checkSshKeys, success: 1 },
    { name: 'Install SSH keys', action: initSsh },
    { name: 'Check firmware', action: checkFirmware, success: 2 },
    { name: 'Initiate firmware upgrade', action: installFirmware },
    { name: 'Wait for firmware upgrade to complete', action: waitForUpgrade },
    // { name: 'Install SSH keys (required after some firmware upgrades)', action: initSsh },
    { name: 'Check VXLAN', action: checkVxlan, success: 1 },
    { name: 'Initiate VXLAN install', action: installVxlan },
    { name: 'Minimize device', action: minimizeDevice },
    { name: 'Setting LAN and WAN addresses', action: setLanWan },
    { name: 'Wait for new LAN and WAN addresses to become available', action: waitForAddresses },
    { name: 'Changing password', action: changePassword },
    { name: 'Initialize wireguard', action: initWireguard },
    { name: 'Setup syslog', action: initSyslog },
    { name: 'Adopted!', action: adopted },
  ])

  onMounted(async () => {
    reset()
    // Load available firmware versions
    try {
      const response = await axios.get('/api/device/firmware/available')
      if (response.data.success) {
        allFirmwareVersions.value = response.data.versions
      }
    } catch (error) {
      console.error('Failed to load firmware versions:', error)
    }
  })

  async function reset() {
    flow.value.forEach(element => {
      element.text = ''
      element.state = ''
      element.icon = ''
      element.color = ''
      element.disabled = true
    });
  }

  async function startOnboarding() {
    reset()
    console.log('adoption flow started')
    await runFlowStep(0)
    console.log('adoption flow ended')
  }

  async function runFlowStep(step) {
    if (step < 0 || step >= flow.value.length) return false
    console.log('running flow step: ' + step + ', flow.length: ' + flow.value.length)
    console.log('running flow step name: ' + flow.value[step].name)
    flow.value[step].disabled = false
    var result = await flow.value[step].action(flow.value[step])
    if (result) {
      flow.value[step].icon = 'mdi-check-circle'
      flow.value[step].color = 'success'
      if (flow.value[step].success) step += flow.value[step].success
      await runFlowStep(step + 1)
    } else {
      flow.value[step].icon = 'mdi-alert-circle'
      flow.value[step].color = 'error'
      if (flow.value[step].success) {
        flow.value[step].color = 'warning'
        await runFlowStep(step + 1)
      }
    }
  }

  async function checkAccessibility(item) {
    console.log('checkAccessibility called')
    const accessible = await device.present(ip.value)
    item.text = accessible ? 'Device accessible' : 'Device not accessible, aborting!'
    return accessible
  }

  async function checkSshKeys(item) {
    console.log('checkSshKeys called')
    const response = await axios.get("/api/device/serial/" + ip.value)
    item.text = JSON.stringify(response.data)
    console.log('response: ' + JSON.stringify(response.data))
    item.state = '- done!'
    return response.data.success
  }

  async function initSsh(item) {
    const body = { username: 'admin', password: pwd.value }
    const response = await axios.post("/api/device/initssh/" + ip.value, body)
    item.text = JSON.stringify(response.data)
    console.log('response: ' + JSON.stringify(response.data))
    item.state = '- done!'
    return response.data.success
  }

  async function checkFirmware(item) {
    console.log('checkFirmware called')
    const response = await axios.get("/api/device/firmware/version/" + ip.value)
    item.text = JSON.stringify(response.data)
    console.log('response: ' + JSON.stringify(response.data))
    const output = response.data.commands[0].output
    const firmware = output.split('\'')[1]
    console.log('firmware: [' + firmware + ']')
    item.state = '- done!'

    // Check if current firmware matches selected firmware
    const targetFirmware = selectedFirmware.value || deviceInfo.value.firmware
    console.log('target firmware: [' + targetFirmware + ']')
    return response.data.success && firmware.includes(targetFirmware)
  }

  async function installFirmware(item) {
    console.log('installFirmware called')
    item.state = '- please wait'

    // Use selected firmware version
    const body = {
      ip: ip.value,
      model: deviceInfo.value.model,
      version: selectedFirmware.value
    }

    const response = await axios.post("/api/device/firmware/upgrade", body)
    item.text = JSON.stringify(response.data)
    console.log('response: ' + JSON.stringify(response.data))
    item.state = '- done!'
    return response.data.success
  }

  async function waitForUpgrade(item) {
    // const wait = (msec) => new Promise((resolve, reject) => {
    //   setTimeout(resolve, msec)
    // })

    // for (var i = 180; i > 0; i--) {
    //   await wait(1000)
    //   item.state = '- in progress (' + i + ' secs left)'
    // }

    item.state = '- in progress '
    var stillwaiting = true
    while (stillwaiting) {
      item.state += '.'
      const body = { username: 'admin', password: pwd.value }
      const response = await axios.post("/api/device/initssh/" + ip.value, body)
      stillwaiting = !response.data.success
    }

    item.state += '- done!'
    return true
  }

  async function checkVxlan(item) {
    console.log('checkVxlan called')
    const response = await axios.get("/api/device/vxlan/check/" + ip.value)
    item.text = JSON.stringify(response.data)
    console.log('response: ' + JSON.stringify(response.data))
    item.state = '- done!'
    return response.data.success
  }

  async function installVxlan(item) {
    console.log('installVxlan called')
    item.state = '- please wait'

    // Use selected firmware version
    const body = {
      ip: ip.value,
      model: deviceInfo.value.model,
      version: selectedFirmware.value
    }

    const response = await axios.post("/api/device/vxlan/install", body)
    item.text = JSON.stringify(response.data)
    console.log('response: ' + JSON.stringify(response.data))
    item.state = '- done!'
    return response.data.success
  }

  async function minimizeDevice(item) {
    console.log('minimizeDevice called')
    item.state = '- please wait'
    const response = await device.runscript(ip.value, 'minimize.sh')
    item.text = JSON.stringify(response)
    console.log('response: ' + JSON.stringify(response))
    item.state = '- done!'
    return response
  }

  async function setLanWan(item) {
    console.log('setLanWan called')
    item.state = '- please wait'

    const cidrlan = parseCidr(lancidr.value)
    const blocklan = new Netmask(lancidr.value)
    const cidrwan = parseCidr(wancidr.value)
    const blockwan = new Netmask(wancidr.value)

    console.log('cidrwan.ip == ip.value: ', cidrwan.ip, ' == ', ip.value)

    // Set hostname and devicename on the device if hostname field is not empty
    if (hostname.value != '') {
      const hostnameCmd = `uci set system.system.hostname='${hostname.value}' && uci set system.system.devicename='${hostname.value}' && uci commit && /etc/init.d/system reload`
      const hostnameResponse = await device.shell(ip.value, hostnameCmd)
      console.log('hostname set response: ' + JSON.stringify(hostnameResponse))
    }

    if (cidrwan.ip != ip.value) {
      const response = await device.runscript(ip.value, 'set_lan_wan.sh', cidrlan.ip + ' ' + cidrwan.ip + ' ' + blocklan.mask)
      item.text = JSON.stringify(response)
      console.log('response: ' + JSON.stringify(response))
      ip.value = cidrlan.ip
    } else {
      item.text = 'Wan IP same as access IP, skipping this step'
    }

    item.state = '- done!'
    return true
  }

  async function waitForAddresses(item) {
    const wait = (msec) => new Promise((resolve, reject) => {
      setTimeout(resolve, msec)
    })

    var notavailable = true
    while (notavailable) {
      for (var i = 10; i > 0; i--) {
        await wait(1000)
        item.state = '- in progress (' + i + ' secs left)'
      }

      const serial = await device.serial(ip.value)
      notavailable = !serial.success
    }

    item.state = '- done!'
    return true
  }

  async function changePassword(item) {
    console.log('changePassword called')
    item.state = '- please wait'

    if (newpwd.value && newpwd.value.length >= 8) {
      const response = await device.runscript(ip.value, 'change_password.sh', newpwd.value)
      item.text = JSON.stringify(response)
      console.log('response: ' + JSON.stringify(response))
      item.state = '- done!'
      return response.success
    }

    item.text = 'New password must be at least 8 characters, skipping this step'
    return true
  }

  async function initWireguard(item) {
    console.log('initWireguard called')
    item.state = '- please wait'
    const cidrlan = parseCidr(lancidr.value)
    const cidrwan = parseCidr(wancidr.value)

    const parts = cidrwan.ip.split('.')
    const wgcidr = '10.48.' + parts[2] + '.' + parts[3] + '/16'

    var response = await device.initwg(ip.value, wgcidr)
    item.text = JSON.stringify(response)
    console.log('response: ' + JSON.stringify(response))
    item.state = '- done!'
    return response.success
  }

  async function initSyslog(item) {
    // TODO: Move to config page
    const response = await device.runscript(ip.value, 'init_syslog.sh', '10.49.88.254 8514 udp')
    item.text = JSON.stringify(response)
    console.log('response: ' + JSON.stringify(response))
    item.state = '- done!'
    return response.success
  }

  async function adopted(item) {
    console.log('adopted called')
    const existing = await device.findBySerial(deviceInfo.value.serialno)
    item.text = JSON.stringify(existing)
    console.log('existing: ' + JSON.stringify(existing))

    var newdevice = {}
    if (existing && existing.id > 0) {
      newdevice = existing
    }

    const wgresponse = await device.wg(ip.value)
    console.log('wg: ' + JSON.stringify(wgresponse))

    if (wgresponse != null) {
      // Set legacy pubkey field for backward compatibility
      newdevice.pubkey = wgresponse['network.wg1.public_key']

      // Create wireguard interface for new structure
      if (!newdevice.wireguardinterfaces) {
        newdevice.wireguardinterfaces = []
      }
      newdevice.wireguardinterfaces.push({
        name: 'wg1',
        publickey: wgresponse['network.wg1.public_key'],
        status: 'unknown'
      })
    }

    const system = await device.uci(ip.value, 'uci show system.system')
    console.log('system: ' + JSON.stringify(system))

    // Set both hostname and device name to the same value from hostname field
    if (hostname.value != '') {
      newdevice.hostname = hostname.value
      newdevice.name = hostname.value
    } else {
      newdevice.hostname = system['system.system.hostname']
      newdevice.name = system['system.system.devicename']
    }
    newdevice.model = system['system.system.device_code']
    newdevice.firmware = system['system.system.device_fw_version']

    const cidrlan = parseCidr(lancidr.value)
    const blocklan = new Netmask(lancidr.value)
    const cidrwan = parseCidr(wancidr.value)
    const blockwan = new Netmask(wancidr.value)

    newdevice.serialno = deviceInfo.value.serialno
    newdevice.ip = cidrlan.ip
    newdevice.networkid = 0
    newdevice.defaultpwd = pwd.value // Ri0k3C6S
    newdevice.password = newpwd.value
    newdevice.hubid = null
    newdevice.state = 'adopted'
    newdevice.role = role.value.toLowerCase()

    // Use helper methods to set interfaces from CIDR values
    device.setInterfaceFromCidr(newdevice, 'lan', lancidr.value)
    device.setInterfaceFromCidr(newdevice, 'wan', wancidr.value)

    // Parse mobile/4G CIDR if provided
    if (apn4gcidr.value && apn4gcidr.value !== '') {
      device.setInterfaceFromCidr(newdevice, 'apn4g', apn4gcidr.value)
    }

    console.log('saving onboarded device: ' + JSON.stringify(newdevice))
    const response = await device.save(newdevice)
    item.text = JSON.stringify(response.data)
    console.log('save response: ' + JSON.stringify(response))
    return response
  }

  async function testAccess() {
    access.value = "Testing, please wait!"
    disableBtn.value = true
    const body = { username: 'admin', password: pwd.value }
    const response = await axios.post("/api/device/initssh/" + ip.value, body)
    console.log('response commands: ' + JSON.stringify(response.data.commands))

    if (!response.data.success) {
      if (response.data.commands[0].error.includes('timeout')) access.value = 'NOT REACHABLE'
      if (response.data.commands[0].output.includes('password was not')) access.value = 'WRONG PASSWORD'
      disableBtn.value = false
      return false
    }

    // Access successful, now get device info
    access.value = "Getting device info..."

    // Get device info (model, firmware, serial, hostname)
    var deviceInfoResponse
    try {
      deviceInfoResponse = await axios.get('/api/device/info/' + ip.value)
      deviceInfo.value = deviceInfoResponse.data
      console.log('device info: ' + JSON.stringify(deviceInfo.value))

      // Filter firmware versions by device model
      if (deviceInfo.value.model) {
        availableFirmware.value = allFirmwareVersions.value.filter(
          fw => fw.model === deviceInfo.value.model
        )

        // Auto-select current firmware or latest version
        if (availableFirmware.value.length > 0) {
          const currentMatch = availableFirmware.value.find(
            fw => deviceInfo.value.firmware.includes(fw.version)
          )
          selectedFirmware.value = currentMatch ? currentMatch.version : availableFirmware.value[0].version
        }
      }
    } catch (error) {
      console.error('Failed to get device info:', error)
      return
    }

    // Check if device exists in database
    access.value = "Checking database..."
    const existingDevice = await device.findBySerial(deviceInfo.value.serialno)
    console.log('existing device: ' + JSON.stringify(existingDevice))

    if (existingDevice && existingDevice.id) {
      // Device found - prefill all fields
      access.value = 'DEVICE FOUND'
      hostname.value = existingDevice.hostname || ''
      newpwd.value = existingDevice.password || ''
      role.value = existingDevice.role ? existingDevice.role.charAt(0).toUpperCase() + existingDevice.role.slice(1) : 'Endpoint'

      // Parse and set CIDR values using embedded interface fields
      if (existingDevice.lan && existingDevice.lanmask) {
        const bits = cidrToBits(existingDevice.lanmask)
        lancidr.value = existingDevice.lan + '/' + bits
      }

      if (existingDevice.wan && existingDevice.wanmask) {
        const bits = cidrToBits(existingDevice.wanmask)
        wancidr.value = existingDevice.wan + '/' + bits
      }

      // if (existingDevice.apn4g && existingDevice.apn4gmask) {
      //   const bits = cidrToBits(existingDevice.apn4gmask)
      //   apn4gcidr.value = existingDevice.apn4g + '/' + bits
      // }

      apn4gcidr.value = deviceInfo.value.apn4gcidr

      console.log('Prefilled fields from existing device plus apn: "' + apn4gcidr.value + '" ');
    } else {
      access.value = 'NEW DEVICE'
    }

    disableBtn.value = false
    return true
  }

  // Helper function to convert netmask to CIDR bits
  function cidrToBits(netmask) {
    const parts = netmask.split('.')
    let bits = 0
    for (const part of parts) {
      const num = parseInt(part)
      bits += (num.toString(2).match(/1/g) || []).length
    }
    return bits
  }

</script>
