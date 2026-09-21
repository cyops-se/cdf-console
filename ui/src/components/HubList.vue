<template>
    <v-card width="100%" height="100%">
        <v-toolbar color="transparent">
            <v-toolbar-title class="text-h6" text="Hubs"></v-toolbar-title>

            <v-text-field v-model="search" label="Search" prepend-inner-icon="mdi-magnify" variant="outlined"
                hide-details single-line></v-text-field>
            <template v-slot:append>
                <v-btn class="ml-4" color="primary" icon="mdi-refresh" @click="refresh"></v-btn>
                <v-btn color="primary" @click="dialog = !dialog">Add hub</v-btn>
            </template>
        </v-toolbar>
        <v-card-text>
            <v-data-table :headers="headers" :items="app.system.hubs" @click:row="onClick">
                <template v-slot:item.lan="{ item }">
                    <div class="text-end">
                        <v-chip v-if="item.lan" :color="item.lanpresent ? 'success' : 'error'" :text="item.lan"
                            class="text-uppercase" size="small" label></v-chip>
                        <span v-else class="text-grey">-</span>
                    </div>
                </template>
                <template v-slot:item.wan="{ item }">
                    <div class="text-end">
                        <v-chip v-if="item.wan" :color="item.wanpresent ? 'success' : 'error'" :text="item.wan"
                            class="text-uppercase" size="small" label></v-chip>
                        <span v-else class="text-grey">-</span>
                    </div>
                </template>
                <template v-slot:item.status="{ item }">
                    <div class="text-end">
                        <v-chip :color="item.status == 'ok' ? 'success' : 'error'" :text="item.status"
                            class="text-uppercase" size="small" label></v-chip>
                    </div>
                </template>
                <template v-slot:item.actions="{ item }">
                    <v-icon size="small" @click.stop="deleteItem(item)">mdi-delete</v-icon>
                </template>
                <template v-slot:top>
                    <v-dialog v-model="dialog">
                        <v-card title="Add hubs">
                            <v-card-text>
                                <v-data-table v-model="selected" :headers="headers" :items="newhubs"
                                    show-select></v-data-table>
                            </v-card-text>
                            <v-card-actions>
                                <v-btn text="Cancel" @click="dialog = false"></v-btn>
                                <v-spacer></v-spacer>
                                <v-btn color="success" text="Add hubs" @click="addHub"></v-btn>
                            </v-card-actions>
                        </v-card>
                    </v-dialog>
                </template>
            </v-data-table>
        </v-card-text>
    </v-card>
</template>

<script setup>
    import axios from "axios"
    import { onMounted } from "vue"
    import { useRouter } from 'vue-router'
    import { useAppStore } from '@/stores/app'
    import { device } from '@/plugins/device'
    import system from '@/plugins/system'

    const props = defineProps(['items'])

    const router = useRouter()
    const app = useAppStore()
    const items = ref([])
    const newhubs = ref([])
    const selected = ref([])
    const dialog = ref(false)
    const headers = ref([{ title: 'Hostname', key: 'hostname', align: 'start' }, { title: 'Serial', key: 'serialno', align: 'end' }, { title: 'WAN', key: 'wan', align: 'end' }, { title: 'LAN', key: 'lan', align: 'end' }, { title: 'Status', key: 'status', align: 'end' }, { title: 'Actions', key: 'actions', align: 'end' }])
    const search = ref('')

    onMounted(async () => {
        refresh()
    })

    async function refresh() {
        items.value = app.system.hubs
        // console.log('system hubs: ', JSON.stringify(app.system.hubs))

        var response = await system.freehubs()
        newhubs.value = response
    }

    async function addHub() {
        // move selected items from newhubs.value to app.device.hubs to avoid them being added more than once
        for (var i = 0; i < selected.value.length; i++) {
            for (var j = newhubs.value.length - 1; j >= 0; j--) {
                if (newhubs.value[j].id == selected.value[i]) {
                    var item = newhubs.value.splice(j, 1)
                    app.system.hubs.push(item[0])
                    break
                }
            }
        }

        selected.value = []
        dialog.value = false
        refresh()
    }

    async function save() {
        // console.log('saving system: ' + JSON.stringify(app.system))
        const response = await system.save(app.system)
        refresh()
    }

    async function deleteItem(item) {
        // Delete from system, not from database
        app.system.hubs = app.system.hubs.filter((hub) => { app.system.id !== hub.systemid })
        delete item.systemid
        await device.save(item)
        refresh()
    }

    async function onClick(e, item) {
        app.device = item.item
        if (app.device.endpoints == null) app.device.endpoints = []
        if (app.device.allowedaddresses == null) app.device.allowedaddresses = []
        // console.log('item clicked: ', JSON.stringify(app.device))
        router.push('/ui/device/hub')
    }
</script>