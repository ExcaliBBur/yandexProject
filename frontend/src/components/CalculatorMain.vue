<template>
  <head>
    <link href="https://use.fontawesome.com/releases/v5.6.1/css/all.css" rel="stylesheet">
    <link href="https://fonts.googleapis.com/css2?family=Lato&display=swap" rel="stylesheet">
  </head>

  <table class="table">
    <tr>
      <td colspan="5">
        <div>
          <HeaderView />
        </div>
      </td>
    </tr>
    <tr>
      <td colspan="5">
        <h1 style="text-align: center;">Calculator</h1>
      </td>
    </tr>
    <tr>
      <td colspan="2">
        <p style="color:red">
          {{ errorMsg }}
        </p>
        <p style="color:green">
          {{ msg }}
        </p>

        <div class="form__group field">
          <input type="input" v-model="expression" class="form__field" placeholder="Name" name="name" id='name'
            required />
          <label for="name" class="form__label">Type your expression here</label>
        </div>
          <button @click="sendExpression" style="--clr:green">
            <span>
              Send
            </span>
            <i></i>
          </button>

          <button @click="resetExpression" style="--clr:red">
            <span>
              Reset
            </span>
            <i></i>
          </button>
      </td>

      <td colspan="1"></td>

      <td colspan="2" align="right">
        <div class="search-container">
          <input type="text" v-model="expressionId" name="search" placeholder="Search id" class="search-input">
          <a @click="handleSearch" class="search-btn">
            <i class="fas fa-search"></i>
          </a>
        </div>
      </td>

    </tr>
    <tr class="border" style="line-height: 50px;">
      <td width="10%" class="centrize border">ID</td>
      <td width="60%" class="centrize border">Expression</td>
      <td width="10%" class="centrize border">Status</td>
      <td width="10%" class="centrize border">Computation time</td>
      <td width="10%" class="centrize border">Error</td>
    </tr>

    <ExpressionContainer v-for="expression in expressions" :key="expression" :content="expression" />

    <tr>
      <td colspan="5">
        <button @click="handlePreviousPage" style="--clr:green">
          <span>
            Previous page
          </span>
          <i></i>
        </button>

        <button @click="handleNextPage" style="--clr:red">
          <span>
            Next page
          </span>
          <i></i>
        </button>

        <p>
          Current page number: {{ currentPage + 1 }};
          Page size: {{ pageSize }}
        </p>
      </td>
    </tr>
  </table>
</template>

<script>
import HeaderView from "@/utility/HeaderView.vue"
import ExpressionContainer from "@/containers/ExpressionContainer.vue"

import axios from 'axios'


export default {
  name: "CalculatorMain",
  components: {
    HeaderView,
    ExpressionContainer,
  },
  data() {
    return {
      errorMsg: null,
      msg: null,

      expression: null,
      expressions: [],
      expressionId: null,

      currentPage: 0,
      pageSize: 10,
    }
  },
  created() {
    this.getExpressionsOnCreation()
  },
  methods: {
    async sendExpression() {
      console.log(this.expressions)
      this.errorMsg = null
      this.msg = null

      let expression = {
        expression: this.expression
      }
      if (expression.expression == null) {
        this.errorMsg = "Expression can not be empty"
        return
      }
      expression.expression = expression.expression.replaceAll(" ", "")
      if (expression.expression == "") {
        this.errorMsg = "Expression can not be empty"
        return
      }

      expression.expression = expression.expression.replaceAll(",", ".")
      const headers = {
        'X-Request-ID': expression.expression,
      }
      await axios.post("http://localhost:8080/api/calculate", expression, {
        headers: headers,
      }).catch((error) => {
        this.errorMsg = "Got status " + error.response.status + " with message: " + error.response.data.message;
        return;
      }).then((token) => {
        this.msg = "Got status " + token.data.code + " with message: " + token.data.message
      })
    },
    resetExpression() {
      this.msg = null
      this.errorMsg = null
      this.expression = null
    },
    getExpressionsOnCreation() {
      let params = new URLSearchParams()
      params.append('page_number', this.currentPage)
      params.append('page_size', this.pageSize)

      axios.get("http://localhost:8080/api/expressions", { params }).catch((error) => {
        this.errorMsg = "Got status " + error.response.status + " with message: " + error.response.data.message;
        return;
      }).then((token) => {
        this.expressions = JSON.parse(token.data.message)

        for (let i = 0; i < this.expressions.length; i++) {
          if (this.expressions[i].is_error) {
            this.expressions[i].error = "divide by zero"
            this.expressions[i].result = "???"
          }

          if (this.expressions[i].is_finished) {
            this.expressions[i].status = "finished"
            this.expressions[i].expression = this.expressions[i].expression + " = "
              + this.expressions[i].result
          } else {
            this.expressions[i].status = "computationing"
            this.expressions[i].expression = this.expressions[i].expression + " = ???"
          }

          let date_start = new Date(Date.parse(this.expressions[i].date_start));
          let date_finish = new Date(Date.parse(this.expressions[i].date_finish));

          this.expressions[i].computation_time = Math.abs(date_finish - date_start)
        }
      })
    },
    handleNextPage() {
      this.errorMsg = null
      this.msg = null
      this.currentPage++

      this.getExpressionsOnCreation()
    },
    handlePreviousPage() {
      this.errorMsg = null
      this.msg = null
      if (this.currentPage != 0) {
        this.currentPage--
      }
      this.getExpressionsOnCreation()
    },
    async handleSearch() {
      this.errorMsg = null
      this.msg = null

      await axios.get("http://localhost:8080/api/expression/"+this.expressionId).catch((error) => {
        this.errorMsg = "Got status " + error.response.status + " with message: " + error.response.data.message;
        return;
      }).then((token) => {
        this.msg = "Got status " + token.data.code + " with message: " + token.data.message
      })
      
      this.expressionId = null
    }
  }
}

</script>

<style scoped lang="scss" src="@/assets/calculate.scss"></style>
<style scoped src="@/assets/button.css"></style>
<style scoped src="@/assets/searchForm.css"></style>


<style scoped>
.centrize {
  text-align: center;
}
</style>