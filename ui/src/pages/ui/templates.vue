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
          <v-text-field v-model="network.name" label="Name" required></v-text-field>
          <v-text-field v-model="network.network" label="Network" required></v-text-field>
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

const items = ref([])
const dialog = ref(false)
const network = ref({ name: '', network: '' })
// const headers = ref([{ title: 'Name', key: 'name' }, { title: 'CIDR', key: 'network' }])
const headers = ref([{ title: 'Name', key: 'name', align: 'start' }, { title: 'CIDR', key: 'network', align: 'end' }, { title: 'Actions', key: 'actions', align: 'end' }])

onMounted(async () => {
  refresh()
})

async function refresh() {
  const response = await axios.get("/api/data/networks")
  items.value = response.data
  console.log('items: ' + JSON.stringify(items.value))
}

async function save() {
  dialog.value = false
  console.log('saving item: ' + JSON.stringify(network.value))
  const response = await axios.put("/api/data/networks", network.value)
  console.log('response: ' + JSON.stringify(response.data))
  refresh()
}

async function editItem(item) {
  console.log('editing item: ' + JSON.stringify(item))
  network.value = item
  dialog.value = true
}

async function deleteItem(item) {
  console.log('deleting item: ' + JSON.stringify(item))
  const response = await axios.delete("/api/data/networks/" + item.id)
  console.log('response: ' + JSON.stringify(response.data))
  refresh()
}
</script>
