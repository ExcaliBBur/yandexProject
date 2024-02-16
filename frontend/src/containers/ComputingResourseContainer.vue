<template>
    <tr class="border">
        <td class="border">{{ id }}</td>

        <td class="border" v-if="is_alive==true" style="color: green">{{ is_alive }}</td>
        <td class="border" v-else style="color: red">{{ is_alive }} This agent will be deleted within two minutes</td>

        <td class="border">{{ last_heartbeat }}</td>
        <td class="border">{{ time_has_passed }}</td>
    </tr>
</template>

<script>

    export default  {
        props: {
            content: JSON
        },
        data() {
            return {
                id: this.content.id,
                is_alive: this.content.is_alive,
                last_heartbeat: new Date(this.content.last_heartbeat),
                time_has_passed: Math.abs(Date.now() - new Date(this.content.last_heartbeat).getTime()) / 1000
            }
        },
    }

</script>

<style scoped>
    td {
        text-align: center;
    }

</style>