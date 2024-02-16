<template>
    <table class='table'>
        <tr>
            <div>
                <HeaderView />
            </div>
        </tr>
        <tr>
            <h1 style="text-align: center;">Calculator Settings</h1>
        </tr>
        <tr>
            <h2>Arithmetic operations bounds ∈ [0;10000] ms</h2>
            <h2>HeartBeat operation bounds ∈ [1;100] s</h2>
        </tr>
        <tr>
            <td>
                <p style="color:red">
                    {{ errorMsg }}
                </p>
                <p style="color:green">
                    {{ msg }}
                </p>

                <div class="form__group field">
                    <input type="input" v-model="plusDuration" class="form__field" placeholder="plus" name="plus" id='plus'
                        required />
                    <label for="plus" class="form__label">Plus duration in ms</label>
                </div>
                <div class="form__group field">

                    <input type="input" v-model="minusDuration" class="form__field" placeholder="minus" name="minus"
                        id='minus' required />
                    <label for="minus" class="form__label">Minus duration in ms</label>
                </div>
                <div class="form__group field">

                    <input type="input" v-model="multiplicationDuration" class="form__field" placeholder="mul" name="mul"
                        id='mul' required />
                    <label for="mul" class="form__label">Multiplication duration in ms</label>
                </div>
                <div class="form__group field">
                    <input type="input" v-model="divisionDuration" class="form__field" placeholder="div" name="div" id='div'
                        required />
                    <label for="div" class="form__label">Division duration in ms</label>
                </div>
                <div class="form__group field">
                    <input type="input" v-model="heartbeatDuration" class="form__field" placeholder="heart" name="heart"
                        id='heart' required />
                    <label for="heart" class="form__label">HeartBeat duration in s</label>
                </div>

                <button @click="sendDuration" style="--clr:green">
                    <span>
                        Send
                    </span>
                    <i></i>
                </button>
            </td>
        </tr>
    </table>
</template>

<script>
import HeaderView from "@/utility/HeaderView.vue"

import axios from 'axios'


export default {
    components: {
        HeaderView,
    },
    created() {
        this.getDurationOnCreation()
    },
    data() {
        return {
            errorMsg: null,
            msg: null,

            plusDuration: 200,
            minusDuration: 200,
            multiplicationDuration: 200,
            divisionDuration: 200,
            heartbeatDuration: 5,
        }
    },
    methods: {
        async sendDuration() {
            this.errorMsg = null
            this.msg = null

            console.log(this.plusDuration)

            let duration = {
                plus_duration: Number(this.plusDuration),
                minus_duration: Number(this.minusDuration),
                mul_duration: Number(this.multiplicationDuration),
                div_duration: Number(this.divisionDuration),
                heartbeat_duration: Number(this.heartbeatDuration)
            }
            await axios.put("http://localhost:8080/api/duration", duration).catch((error) => {
                this.errorMsg = "Got status " + error.response.status + " with message: " + error.response.data.message;
                return;
            }).then((token) => {
                this.msg = "Got status " + token.data.code + " with message: " + token.data.message
            })
        },
        getDurationOnCreation() {
            axios.get("http://localhost:8080/api/duration").catch((error) => {
                this.errorMsg = "Got status " + error.response.status + " with message: " + error.response.data.message;
                return;
            }).then((token) => {
                let obj = JSON.parse(token.data.message)
                this.plusDuration = obj.plus_duration
                this.minusDuration = obj.minus_duration
                this.multiplicationDuration = obj.mul_duration
                this.divisionDuration = obj.div_duration
                this.heartbeatDuration = obj.heartbeat_duration
            })
        }
    }
}

</script>

<style scoped lang="scss" src="@/assets/calculate.scss"></style>
<style scoped src="@/assets/button.css"></style>
<style scoped>
.form__field {
    width: 20%
}
</style>