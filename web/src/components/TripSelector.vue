// bind input elements to data elements
// use as input to Trip component

<template>
  <div>
    <input type="date" id="date" />
    <input type="radio" id="direction" />
  </div>
</template>

<script>
module.exports = {
  data
}

let date = Date.now()
let direction = null

function data () { return {
  date: date,
  direction: direction
}}

</script>

<style>
</style>
