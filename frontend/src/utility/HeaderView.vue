<template>
    <table>
      <tr>
        <td>
          <RouterLink class='routerLink' to="calculator">Calculator</RouterLink>
        </td>
        <td>
          <RouterLink class='routerLink' to="calculator-settings">Calculator Settings</RouterLink>
        </td>
        <td>
          <RouterLink class='routerLink' to="computing-resourses">Computing Resourses</RouterLink>
        </td>
      </tr>
    </table>
</template>

<style scoped>

  .routerLink {
    font-size: 30px;
    margin-left: 6px;
    margin-right: 6px;
  }

</style>