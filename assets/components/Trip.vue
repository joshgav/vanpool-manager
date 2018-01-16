<template>
  <div class="trip trip-{{ date }} trip-{{ direction }}">
    <div class="trip-title">{{ date }} - {{ direction }}</div>
    <div class="trip-button">
      <button v-if="!(present())" action="addRider()">add self</button>
      <button v-if="present()" action="removeRider()">remove self</button>
    </div>
    <div>
      <ul v-for="rider in riders">
        <li>{{ rider.name }}</li>
      </ul>
    </div>
  </div>
</template>

<script>
module.exports = {
  data,
  props
}

let axios = require('axios')

const apiHostname = 'localhost'
const apiBasePath = '/api/v1/
const apiScheme = 'http'

const directions = {
  in: "Inbound",
  out: "Outbound"
}

function Rider(name, date, direction) {
  this.name = name
  this.date = this.formatDate(date)
  this.direction = direction
}

let riders = []

function data () { return {
  riders: riders
  currentUser: new Rider()
}}

let props = ['date', 'direction']

function mounted() {
  // abort if props aren't set, but refresh when they are
  this.refreshCurrentUser()
  this.refreshRiders()
}

// @param date date in format 'yyyy-mm-dd'
// @param direction 'out' or 'in'
function refreshRiders (date, direction) {
  let baseUrl = `${apiScheme}://${apiHostname}${apiBasePath}riders`
  let queryParams = `?date=${formatDate(date)}&direction=${direction}`
  axios.get(`${baseUrl}${queryParams}`)
    .then(response => {
      console.log(`got riders`)
      this.riders = response.riders
    })
    .catch(err => {
      console.log(`failed to get riders: ${err}`)
      this.riders = []
    })
}

function present() {
  for (rider in riders) {
    if (rider.name == this.currentUser.name) {
      return true
    }
  }
  return false
}

function refreshCurrentUser () {
  let baseUrl = `${apiScheme}://${apiHostname}${apiBasePath}currentUser`
  axios.get(baseUrl)
    .then(response => {
      console.log(`got current user`)
      this.currentUser = response.currentUser
    })
    .catch(err => {
      console.log(`failed to get current user: ${err}`)
      this.currentUser = new Rider()
    })
}

function addRider (rider) {
  let rider = new Rider(currentUser.name, this.date, this.direction)
  axios.put(`${apiScheme}://${apiHostname}${apiBasePath}riders`, JSON.stringify(rider))
    .catch(err => { console.log(`failed to add rider: ${err}`) })
}

function removeRider (rider) {
  let rider = new Rider(currentUser.name, this.date, this.direction)
  // how to do parameterized delete?
  axios.delete(`${apiScheme}://${apiHostname}${apiBasePath}riders`, JSON.stringify(rider))
    .catch(err => { console.log(`failed to remove rider: ${err}`) })
}

// API expects "yyyy-mm-dd"
function formatDate (date) {
  return `${date.getFullYear()}-${fixLength(date.getMonth())}-${fixLength(date.getDate())}`
}

function fixLength (inString) {
  if inString.length == 1 {
    inString = `0${inString}`
  }
  return inString
}

</script>

<style>
</style>
