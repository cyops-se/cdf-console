<template>
  <v-card title="Networks" width="100%" height="100%">
    <template v-slot:append>
      <v-btn color="success" icon="mdi-plus" @click="dialog = true"></v-btn>
    </template>
    <v-card-text>
      <v-data-table density="compact" :headers="headers" :items="items">
        <template v-slot:item.actions="{ item }">
          <v-icon class="me-2" size="small" @click="editItem(item)">
            mdi-pencil
          </v-icon>
          <v-icon size="small" @click="deleteItem(item)">
            mdi-delete
          </v-icon>
        </template>
      </v-data-table>
    </v-card-text>
  </v-card>
  <v-dialog max-width="500" v-model="dialog">
    <template v-slot:default="{ isActive }">
      <v-card title="Network">
        <v-card-text>
          <v-form v-model="valid">
            <v-text-field v-model="network.name" label="Name" required></v-text-field>
            <v-text-field v-model="network.network" label="Network" required validate-on="input"
              :rules="cidrRules"></v-text-field>
          </v-form>
        </v-card-text>

        <v-card-actions>
          <v-btn text="Cancel" @click="isActive.value = false"></v-btn>
          <v-spacer></v-spacer>
          <v-btn color="success" text="Save" @click="save"></v-btn>
        </v-card-actions>
      </v-card>
    </template>
  </v-dialog>
</template>

<script setup>
import axios from "axios"
import { onMounted } from "vue"
import { useAppStore } from '@/stores/app'

const app = useAppStore()
const cidrRegex = new RegExp('^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(\/(3[0-2]|[1-2][0-9]|[0-9]))$')
const items = ref([])
const dialog = ref(false)
const valid = ref(false)
const network = ref({ name: '', network: '' })
const headers = ref([{ title: 'Name', key: 'name', align: 'start' }, { title: 'Network (CIDR)', key: 'network', align: 'end' }, { title: 'Actions', key: 'actions', align: 'end' }])
const cidrRules = ref([
  value => {
    return cidrRegex.test(value) ? true : 'Network in CIDR format is required.'
  },
])

onMounted(async () => {
  refresh()
})

async function refresh() {
  const response = await axios.get("/api/data/networks")
  items.value = response.data
}

async function save() {
  if (!valid.value) return
  dialog.value = false
  const response = await axios.put("/api/data/networks", network.value)
  refresh()
}

async function editItem(item) {
  console.log('editing item: ' + JSON.stringify(item))
  network.value = item
  dialog.value = true
}

async function deleteItem(item) {
  var response = await axios.get("/api/data/devices/field/network_id/" + item.id)
  if (response.data && response.data.length > 0) {
    app.setalert('Network still has a number of devices. Cannot delete it!')
  } else {
    response = await axios.delete("/api/data/networks/" + item.id)
    if (response.status === 200) app.clearalert()
  }
  refresh()
}
</script>
