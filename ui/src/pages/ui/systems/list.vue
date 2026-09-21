<template>
  <v-card width="100%" height="100%">
    <v-toolbar color="transparent">
      <v-toolbar-title class="text-h6" text="Systems"></v-toolbar-title>

      <v-text-field v-model="search" label="Search" prepend-inner-icon="mdi-magnify" variant="outlined" hide-details
        single-line></v-text-field>
      <template v-slot:append>
        <v-btn class="ml-4" color="primary" icon="mdi-refresh" @click="refresh"></v-btn>
        <v-btn class="ml-2" color="success" icon="mdi-plus" @click="onNew"></v-btn>
      </template>
    </v-toolbar>
    <v-card-text>
      <v-data-table v-model:sort-by="sortBy" density="compact" :headers="headers" :items="items" :search="search"
        @click:row="onClick">
      </v-data-table>
    </v-card-text>
  </v-card>
  <v-dialog max-width="500" v-model="dialog">
    <template v-slot:default="{ isActive }">
      <v-card title="System">
        <v-card-text>
          <v-text-field v-model="app.system.name" label="Name" required></v-text-field>
        </v-card-text>
        <v-card-actions>
          <v-btn text="Cancel" @click="isActive.value = false"></v-btn>
          <v-spacer></v-spacer>
          <v-btn color="success" text="New" @click="save"></v-btn>
        </v-card-actions>
      </v-card>
    </template>
  </v-dialog>
</template>

<script setup>
import { onMounted } from "vue"
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import system from '@/plugins/system'

const router = useRouter()
const app = useAppStore()
const items = ref([])
const sortBy = ref([{ key: 'count', order: 'desc' }])
const search = ref('')
const dialog = ref(false)
const headers = ref([{ title: 'Name', key: 'name' }, { title: 'Actions', key: 'actions', align: 'end' }])

onMounted(async () => {
  if (!app.system) app.system = {}
  refresh()
})

async function refresh() {
  const response = await system.allsystems()
  console.log('systems: ' + JSON.stringify(response))
  items.value = response
}

async function onNew() {
  app.system = {}
  dialog.value = true
}

async function save() {
  dialog.value = false
  const response = await system.save(app.system)
  refresh()
}

async function onClick(e, item) {
  app.system = item.item
  if (app.system.networks == null) app.system.networks = []
  if (app.system.hubs == null) app.system.hubs = []
  router.push('/ui/systems/edit')
}

</script>