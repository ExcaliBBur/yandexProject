<template>
    <table class='table'>
        <tr>
            <td colspan="4">
                <div>
                    <HeaderView />
                </div>
            </td>
        </tr>
        <tr>
            <td colspan="4">
                <h1 style="text-align: center;">Computing Resourses</h1>
            </td>
        </tr>
        <p style="color:red">
            {{ errorMsg }}
        </p>

        <tr class="border">
            <td class="centrize border">ID</td>
            <td class="centrize border">Is alive</td>
            <td class="centrize border">Last HeartBeat</td>
            <td class="centrize border">Seconds has passed from last heartbeat</td>
        </tr>

        <ComputingResourseContainer v-for="resourse in computingResourses" :key="resourse" :content="resourse" />

    </table>
</template>

<script>
import HeaderView from "@/utility/HeaderView.vue"
import ComputingResourseContainer from "@/containers/ComputingResourseContainer.vue"

import axios from 'axios'

export default {
    components: {
        HeaderView,
        ComputingResourseContainer,
    },
    created() {
        this.getComputingResoursesOnCreate()
    },
    data() {
        return {
            errorMsg: null,

            computingResourses: [],
        }
    },
    methods: {
        getComputingResoursesOnCreate() {
            axios.get("http://localhost:8080/api/workers").catch((error) => {
                this.errorMsg = "Got status " + error.response.status + " with message: " + error.response.data.message;
                return;
            }).then((token) => {
                this.computingResourses = JSON.parse(token.data.message)
            })
        }
    }
}

</script>

<style scoped>
.centrize {
    text-align: center;
}

</style>