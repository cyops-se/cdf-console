<template>
    <v-card title="System" width="100%" height="100%">
        <template v-slot:append>
            <v-btn class="ml-2" icon="mdi-import" @click="importdialog = true"></v-btn>
            <v-btn class="ml-2" icon="mdi-content-save" @click="save"></v-btn>
        </template>

        <v-card>
            <v-card-text>
                <v-text-field v-model="app.system.name" label="Name" required></v-text-field>
                <HubList :items="app.system.hubs" />
            </v-card-text>
        </v-card>
    </v-card>

    <v-dialog max-width="500" v-model="importdialog">
        <template v-slot:default="{ isActive }">
            <v-card title="Import hub">
                <v-card-text>
                    <v-text-field v-model="importip" label="IP address to import" required></v-text-field>
                    <v-text-field v-model="importpwd" label="Password to use during import" required></v-text-field>
                </v-card-text>
                <v-card-actions>
                    <v-btn text="Cancel" @click="isActive.value = false"></v-btn>
                    <v-spacer></v-spacer>
                    <v-btn color="success" text="Import" @click="importhub"></v-btn>
                </v-card-actions>
            </v-card>
        </template>
    </v-dialog>
</template>

<script setup>
import { onMounted } from "vue"
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { parseCidr } from 'cidr-tools'
import { device } from '@/plugins/device'
import system from '@/plugins/system'

const router = useRouter()
const app = useAppStore()
const importdialog = ref(false)
const importip = ref('192.168.1.1')
const importpwd = ref('admin01')

onMounted(async () => {
    if (!app.system) app.system = {}
    refresh()
})

async function refresh() {
    console.log('edit refresh')
    await device.checklist(app.system.hubs)
}

async function save() {
    console.log('saving system: ' + JSON.stringify(app.system))
    const response = await system.save(app.system)
    app.showsnack('System saved!')
    refresh()
}

async function importhub() {
    const hub = await device.importhub(importip.value, importpwd.value)
    app.system.hubs.push(hub)
    save()

    importdialog.value = false
    app.showsnack('Import completed')
    refresh()
}

</script>