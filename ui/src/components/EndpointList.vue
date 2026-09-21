<template>
    <v-card width="100%" height="100%">
        <v-toolbar color="transparent">
            <v-toolbar-title class="text-h6" text="Devices"></v-toolbar-title>

            <v-text-field v-model="search" label="Search" prepend-inner-icon="mdi-magnify" variant="outlined"
                hide-details single-line></v-text-field>
            <template v-slot:append>
                <v-btn class="ml-4" color="primary" icon="mdi-refresh" @click="refresh"></v-btn>
                <v-btn class="ml-2" @click="scandialog = true">Scan</v-btn>
                <v-btn class="ml-2" color="success" icon="mdi-plus" @click="dialog = true"></v-btn>
            </template>
        </v-toolbar>
        <v-card-text>
            <v-data-table :headers="headers" :items="hub.endpoints">
                <template v-slot:item.actions="{ item }">
                    <!-- <v-icon class="me-2" size="small" @click="editItem(item)">mdi-pencil</v-icon> -->
                    <v-icon class="me-2" size="small" @click="configureItem(item)">mdi-pencil</v-icon>
                    <v-icon size="small" @click.stop="deleteEndpoint(item)">mdi-delete</v-icon>
                </template>
                <template v-slot:top>
                    <v-dialog v-model="dialog">
                        <template v-slot:activator="{ props }">
                            <v-toolbar>
                                <v-toolbar-title>Endpoints</v-toolbar-title>
                                <v-spacer />
                                <v-btn color="primary" @click="dialog = !dialog">Add endpoints</v-btn>
                            </v-toolbar>
                        </template>
                        <v-card title="Add enpoints">
                            <v-card-text>
                                <v-data-table v-model="selected" :headers="headers" :items="endpoints"
                                    show-select></v-data-table>
                            </v-card-text>
                            <v-card-actions>
                                <v-btn text="Cancel" @click="dialog = false"></v-btn>
                                <v-spacer></v-spacer>
                                <v-btn color="success" text="Add endpoint" @click="addEndpoint"></v-btn>
                            </v-card-actions>
                        </v-card>
                    </v-dialog>
                </template>
            </v-data-table>
        </v-card-text>
    </v-card>
    <v-dialog max-width="500" v-model="dialog">
        <template v-slot:default="{ isActive }">
            <v-card title="Device">
                <v-card-text>
                    <v-text-field v-model="app.device.name" label="Name" required></v-text-field>
                    <v-text-field v-model="app.device.hostname" label="Hostname" required></v-text-field>
                    <v-text-field v-model="app.device.serialno" label="Serial" required></v-text-field>
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
    import { useRouter } from 'vue-router'
    import { expandCidr, parseCidr } from "cidr-tools"
    import { useAppStore } from '@/stores/app'
    import device from '@/plugins/device'

    const props = defineProps(['hub'])

    const router = useRouter()
    const app = useAppStore()
    const items = ref([])
    const networks = ref([])
    const dialog = ref(false)
    const scandialog = ref(false)
    // const headers = ref([{ title: 'Hostname', key: 'hostname', align: 'start' }, { title: 'Role', key: 'role', align: 'end' }, { title: 'Hub', key: 'hubid', align: 'end' }, { title: 'Serial', key: 'serialno', align: 'end' }, { title: 'WAN', key: 'wan', align: 'end' }, { title: 'APN 4G', key: 'apn4g', align: 'end' }, { title: 'LAN', key: 'lan', align: 'end' }, { title: 'Public key', key: 'pubkey', align: 'end' }, { title: 'Status', key: 'status', align: 'end' }, { title: 'Actions', key: 'actions', align: 'end' }])
    // const headers = ref([{ title: 'Hostname', key: 'hostname', align: 'start' }, { title: 'Role', key: 'role', align: 'end' }, { title: 'Hub', key: 'hubid', align: 'end' }, { title: 'Serial', key: 'serialno', align: 'end' }, { title: 'WAN', key: 'wan', align: 'end' }, { title: 'LAN', key: 'lan', align: 'end' }, { title: 'Status', key: 'status', align: 'end' }, { title: 'Actions', key: 'actions', align: 'end' }])
    const headers = ref([{ title: 'Hostname', key: 'hostname', align: 'start' }, { title: 'Serial', key: 'serialno', align: 'end' }, { title: 'WAN', key: 'wan', align: 'end' }, { title: 'LAN', key: 'lan', align: 'end' }, { title: 'APN', key: 'apn4g', align: 'end' }, { title: 'Status', key: 'status', align: 'end' }, { title: 'Public key', key: 'pubkey', align: 'end' }, { title: 'Actions', key: 'actions', align: 'end' }])
    const search = ref('')
    const scancidr = ref('192.168.1.0/24')
    const scanning = ref('')
    const progress = ref(0)

    onMounted(async () => {
        refresh()
    })

    async function refresh() {
    }

    async function save() {
    }

    async function deleteItem(item) {
    }

    async function onClick(e, item) {
    }

</script>