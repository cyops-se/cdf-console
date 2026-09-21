<template>
  <v-card>
    <v-card-title>Network Interfaces</v-card-title>
    <v-card-text>
      <v-data-table
        :headers="headers"
        :items="interfaces"
        density="compact"
      >
        <template v-slot:item.present="{ item }">
          <v-chip
            :color="item.present ? 'success' : 'error'"
            size="x-small"
            label
          >
            {{ item.present ? 'Present' : 'Not Present' }}
          </v-chip>
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
              Add Interface
            </v-btn>
          </v-toolbar>
        </template>
      </v-data-table>
    </v-card-text>

    <!-- Edit Dialog -->
    <v-dialog v-model="dialog" max-width="500px">
      <v-card>
        <v-card-title>{{ editedIndex === -1 ? 'New' : 'Edit' }} Interface</v-card-title>
        <v-card-text>
          <v-select
            v-model="editedItem.name"
            :items="['wan', 'lan', 'apn4g']"
            label="Interface Type"
            :disabled="editedIndex !== -1"
          ></v-select>
          <v-text-field
            v-model="editedItem.ipaddress"
            label="IP Address"
          ></v-text-field>
          <v-text-field
            v-model="editedItem.netmask"
            label="Netmask"
          ></v-text-field>
          <v-checkbox
            v-model="editedItem.present"
            label="Present"
          ></v-checkbox>
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
  present: false
})

const defaultItem = {
  name: '',
  ipaddress: '',
  netmask: '',
  present: false
}

const headers = ref([
  { title: 'Type', key: 'name', align: 'start' },
  { title: 'IP Address', key: 'ipaddress', align: 'start' },
  { title: 'Netmask', key: 'netmask', align: 'start' },
  { title: 'Status', key: 'present', align: 'center' },
  { title: 'Actions', key: 'actions', align: 'end', sortable: false }
])

const interfaces = computed({
  get: () => props.modelValue || [],
  set: (value) => emit('update:modelValue', value)
})

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
  if (confirm('Are you sure you want to delete this interface?')) {
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
