import axios from "axios"

export default {

    create: async (system) => {
        var response = await axios.post('/api/data/systems', system)
        return response.data
    },

    save: async (system) => {
        var response = await axios.put('/api/data/systems', system)
        return response.data
    },

    hubs: async (system) => {
        var response = await axios.get('/api/data/device/field/role/hub')
        const available = response.data.filter((e) => e.systemid === system.id)
        return available
    },

    freehubs: async () => {
        var response = await axios.get('/api/data/devices/field/role/hub')
        const available = response.data.filter((e) => e.systemid !== null)
        return available
    },

    allsystems: async () => {
        var response = await axios.get('/api/data/systems')
        return response.data
    },
}
