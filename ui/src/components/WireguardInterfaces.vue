<template>
  <v-card>
    <v-card-title>Wireguard Interfaces</v-card-title>
    <v-card-text>
      <v-data-table
        :headers="headers"
        :items="interfaces"
        density="compact"
      >
        <template v-slot:item.status="{ item }">
          <v-chip
            :color="getStatusColor(item.status)"
            size="x-small"
            label
          >
            {{ item.status || 'unknown' }}
          </v-chip>
        </template>
        <template v-slot:item.publickey="{ item }">
          <code class="text-caption">{{ truncateKey(item.publickey) }}</code>
        </template>
        <template v-slot:item.actions="{ item }">
          <v-btn
            icon="mdi-pencil"
            size="x-small"
            variant="text"
            @click="editInterface(item)"
          ></v-btn>
          <v-btn
            icon="mdi-delete"
            size="x-small"
            variant="text"
            color="error"
            @click="deleteInterface(item)"
          ></v-btn>
        </template>
        <template v-slot:top>
          <v-toolbar density="compact">
            <v-spacer></v-spacer>
            <v-btn
              color="primary"
              size="small"
              @click="addInterface"
            >
              Add Wireguard Interface
            </v-btn>
          </v-toolbar>
        </template>
      </v-data-table>
    </v-card-text>

    <!-- Edit Dialog -->
    <v-dialog v-model="dialog" max-width="600px">
      <v-card>
        <v-card-title>{{ editedIndex === -1 ? 'New' : 'Edit' }} Wireguard Interface</v-card-title>
        <v-card-text>
          <v-row>
            <v-col cols="12" sm="6">
              <v-text-field
                v-model="editedItem.name"
                label="Interface Name (e.g., wg0)"
                :disabled="editedIndex !== -1"
              ></v-text-field>
            </v-col>
            <v-col cols="12" sm="6">
              <v-text-field
                v-model.number="editedItem.listenport"
                label="Listen Port"
                type="number"
              ></v-text-field>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" sm="6">
              <v-text-field
                v-model="editedItem.ipaddress"
                label="IP Address"
              ></v-text-field>
            </v-col>
            <v-col cols="12" sm="6">
              <v-text-field
                v-model="editedItem.netmask"
                label="Netmask"
              ></v-text-field>
            </v-col>
          </v-row>
          <v-text-field
            v-model="editedItem.publickey"
            label="Public Key"
            readonly
          ></v-text-field>
          <v-text-field
            v-model="editedItem.endpoint"
            label="Endpoint (host:port)"
          ></v-text-field>
          <v-text-field
            v-model="editedItem.allowedips"
            label="Allowed IPs (comma-separated)"
          ></v-text-field>
          <v-select
            v-model="editedItem.status"
            :items="['unknown', 'not configured', 'not active', 'active']"
            label="Status"
          ></v-select>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text="Cancel" @click="close"></v-btn>
          <v-btn color="primary" text="Save" @click="save"></v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-card>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue'])

const dialog = ref(false)
const editedIndex = ref(-1)
const editedItem = ref({
  name: '',
  ipaddress: '',
  netmask: '',
  listenport: 51820,
  publickey: '',
  privatekey: '',
  endpoint: '',
  allowedips: '',
  status: 'unknown'
})

const defaultItem = {
  name: '',
  ipaddress: '',
  netmask: '',
  listenport: 51820,
  publickey: '',
  privatekey: '',
  endpoint: '',
  allowedips: '',
  status: 'unknown'
}

const headers = ref([
  { title: 'Name', key: 'name', align: 'start' },
  { title: 'IP Address', key: 'ipaddress', align: 'start' },
  { title: 'Port', key: 'listenport', align: 'center' },
  { title: 'Public Key', key: 'publickey', align: 'start' },
  { title: 'Status', key: 'status', align: 'center' },
  { title: 'Actions', key: 'actions', align: 'end', sortable: false }
])

const interfaces = computed({
  get: () => props.modelValue || [],
  set: (value) => emit('update:modelValue', value)
})

function getStatusColor(status) {
  switch (status) {
    case 'active': return 'success'
    case 'not active': return 'warning'
    case 'not configured': return 'error'
    default: return 'grey'
  }
}

function truncateKey(key) {
  if (!key) return ''
  return key.length > 20 ? key.substring(0, 20) + '...' : key
}

function addInterface() {
  editedItem.value = Object.assign({}, defaultItem)
  editedIndex.value = -1
  dialog.value = true
}

function editInterface(item) {
  editedIndex.value = interfaces.value.indexOf(item)
  editedItem.value = Object.assign({}, item)
  dialog.value = true
}

function deleteInterface(item) {
  const index = interfaces.value.indexOf(item)
  if (confirm('Are you sure you want to delete this Wireguard interface?')) {
    const newInterfaces = [...interfaces.value]
    newInterfaces.splice(index, 1)
    interfaces.value = newInterfaces
  }
}

function close() {
  dialog.value = false
  setTimeout(() => {
    editedItem.value = Object.assign({}, defaultItem)
    editedIndex.value = -1
  }, 300)
}

function save() {
  const newInterfaces = [...interfaces.value]
  if (editedIndex.value > -1) {
    Object.assign(newInterfaces[editedIndex.value], editedItem.value)
  } else {
    newInterfaces.push(editedItem.value)
  }
  interfaces.value = newInterfaces
  close()
}
</script>
