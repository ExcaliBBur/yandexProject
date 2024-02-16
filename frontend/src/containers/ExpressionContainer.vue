<template>
    <tr class="border">
        <td class="border">{{ id }}</td>
        <td class="border">{{ expression }}</td>
        <td class="border" v-if="status==`finished`" style="color: green">{{ status }}</td>
        <td class="border" v-else style="color: yellow">{{ status }}</td>

        <td class="border">{{ computation_time }} ms</td>

        <td class="border" v-if="error!=null" style="color: red;">{{ error }}</td>
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
                expression: this.content.expression,
                result: this.content.result,
                status: this.content.status,
                error: this.content.error,
                date_start: this.content.date_start,
                date_finish: this.content.date_finish,
                computation_time: this.content.computation_time,
            }
        },
    }

</script>

<style scoped>
    td {
        text-align: center;
    }
</style>