<template>
    <v-card>
        <v-card-text>
            <v-list :items="app.device.allowedaddresses" density="compact">
                <v-list-subheader>Device IPs connected to this {{ type }}</v-list-subheader>
                <v-list-item v-for="(item, i) in app.device.allowedaddresses" :key="i" :value="item">
                    <template v-slot:prepend>
                        <v-icon :icon="item.icon"></v-icon>
                    </template>

                    <v-list-item-title v-text="item.address"></v-list-item-title>
                    <template v-slot:append v-slot:item.actions="{ item }">
                        <v-btn color="grey-lighten-1" icon="mdi-information" variant="text"></v-btn>
                        <v-icon size="small" @click.stop="removeDevice(item)" color="error">mdi-delete</v-icon>
                    </template>
                </v-list-item>
            </v-list>
            <v-row>
                <v-col cols="12" sm="8">
                <v-text-field :label="'Device connected to this ' + type" v-model="newdevice"></v-text-field>
                </v-col>
                <v-col cols="12" sm="4">
                    <v-btn class="ml-2" @click="addDevice" color="primary">Add new device</v-btn>
                </v-col>
            </v-row>
        </v-card-text>
    </v-card>
</template>

<script>
export default {
  props: {
    type: {
      type: String,
      required: true
    }
  }
}
</script>

<script setup>
import axios from "axios"
import { useAppStore } from '@/stores/app'
const app = useAppStore()
const newdevice = ref('')

async function addDevice() {
    if (!app.device.allowedaddresses) app.device.allowedaddresses = []

    for (var i = 0; i < app.device.allowedaddresses.length; i++) {
        if (app.device.allowedaddresses[i].address == newdevice.value) {
          return;
        }
    }

    app.device.allowedaddresses.push({ address: newdevice.value })
}

async function removeDevice(item) {
    // console.log('deleting allowed device: ' + JSON.stringify(item))

    for (var i = 0; i < app.device.allowedaddresses.length; i++) {
      if (app.device.allowedaddresses[i].address == item.address) {
          app.device.allowedaddresses.splice(i, 1);
          // console.log(JSON.stringify(item));
          await axios.delete(`/api/data/addresses/${item.id}`);
          break
      }
    }

    if (app.device.allowedaddresses == null) {
      app.device.allowedaddresses = [];
    }
}
</script>
