<template>
  <v-card width="100%" height="100%">
    <v-toolbar color="transparent">
      <v-toolbar-title class="text-h6" text="Flows"></v-toolbar-title>

      <v-text-field v-model="search" label="Search" prepend-inner-icon="mdi-magnify" variant="outlined" hide-details
        single-line></v-text-field>
      <template v-slot:append>
        <v-btn class="ml-4" color="primary" icon="mdi-refresh" @click="refresh"></v-btn>
      </template>
    </v-toolbar>
    <!-- <template v-slot:append>
    </template> -->
    <v-card-text>
      <v-data-table v-model:sort-by="sortBy" density="compact" :items="items" :search="search">
      </v-data-table>
    </v-card-text>
  </v-card>
</template>

<script setup>
import axios from "axios"
import { onMounted } from "vue"

const items = ref([])
const sortBy = ref([{ key: 'count', order: 'desc' }])
const search = ref('')

onMounted(async () => {
  refresh()
})

async function refresh() {
  const response = await axios.get("/api/flows/grey")
  items.value = response.data.entries
}

</script>