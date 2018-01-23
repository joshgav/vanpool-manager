<template>
  <div class="trip">
    <div class="trip-selector">
      <label for="date">Date: </label>
      <input
        type="date"
        id="date"
        v-model="date"
      />

      <div>
        <input
          type="radio"
          name="direction"
          value="Inbound"
          v-bind:checked="isInbound()"
          v-model="direction"
        />
        <label for="directionIn">Inbound</label>

        <input type="radio"
          id="directionOut"
          name="direction"
          value="Outbound"
          v-bind:checked="!isInbound()"
          v-model="direction"
        />
        <label for="directionOut">Outbound</label>
      </div>
    </div>

    <div class="trip-details">
      <div class="trip-title">
        {{ this.date }} - {{ this.direction }}
      </div>
      <div class="trip-button">
        <button type="button" v-if="!present()"
          v-on:click="addRider()">Add self</button>
        <button type="button" v-if="present()"
          v-on:click="removeRider()">Remove self</button>
      </div>
      <div class="trip-riders">
        <ul>
          <li v-for="rider in riders">{{ rider.displayName }}</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script>
const axios = require('axios');

exports = module.exports = {
  data: data,
  props: props,
  methods: {
    present,
    isInbound,
    Rider,
    addRider,
    removeRider,
    refresh,
    refreshCurrentUser,
    refreshRiders,
  },
  mounted: mounted, // without this `this` is undefined in hooks
  updated: updated
}

// note: parent is not updated when props change
let props = ['date', 'direction'];

function data () { return {
  riders: [],
  currentUser: new Rider(),
  date: this.date || formatDate(new Date()),
  direction: this.direction || 'Inbound',
}}

const apiHostname = 'localhost'
const apiBasePath = '/api/v1/'
const apiScheme = 'http'

function Rider(username, displayName, date, direction) {
  this.username = "" + username;
  this.displayName = "" + displayName;
  this.date = formatDate(date ? date : new Date());
  this.direction = direction ? direction : "Inbound";
}

function mounted() {
  this.refresh(true)
}

function updated() {
  this.refresh(false)
}

function refresh(force) {
  if ((this.date != null) && (this.direction != null)) {
    this.refreshCurrentUser(force)
    this.refreshRiders(this.date, this.direction)
  }
}

// @param date date in format 'yyyy-mm-dd'
// @param direction 'out' or 'in'
function refreshRiders (date, direction) {
  var self = this
  let baseUrl = `${apiScheme}://${apiHostname}${apiBasePath}riders`
  let queryParams = `?date=${formatDate(date)}&direction=${direction}`
  axios.get(`${baseUrl}${queryParams}`)
    .then(response => {
      console.log(`got riders`)
      self.riders = response.riders
    })
    .catch(err => {
      console.log(`failed to get riders: ${err}`)
      self.riders = []
    })
}

function present() {
  for (rider in this.riders) {
    if (rider.username == this.currentUser.username) {
      return true
    }
  }
  return false
}

function refreshCurrentUser (force) {
  var self = this
  if (!force && this.currentUser.username) { return }
  let baseUrl = `${apiScheme}://${apiHostname}${apiBasePath}self`
  axios.get(baseUrl)
    .then(response => {
      console.log(`got current user`)
      self.currentUser = response.self
    })
    .catch(err => {
      console.log(`failed to get current user: ${err}`)
      self.currentUser = new Rider()
    })
}

function addRider () {
  axios.put(
    `${apiScheme}://${apiHostname}${apiBasePath}riders`,
    JSON.stringify(getRider())
   )
  .catch(err => {
    console.log(`failed to add rider: ${err}`)
  })
}

function removeRider () {
  axios.delete(
    `${apiScheme}://${apiHostname}${apiBasePath}riders`,
    JSON.stringify(getRider())
  )
  .catch(err => {
    console.log(`failed to remove rider: ${err}`)
  })
}

function getRider () {
  return new Rider(
    this.currentUser.username,
    this.currentUser.displayName,
    this.date,
    this.direction
  )
}

function formatDate (date) {
  var typedDate // : Date
  try { typedDate = new Date(date) } catch (e) {}
  if (typedDate != null) {
    let year  = typedDate.getFullYear()
    let month = fixLength(typedDate.getMonth() + 1)
    let day   = fixLength(typedDate.getDate())
    return `${year}-${month}-${day}`
  }
  else {
    return date
  }
}

function fixLength (_in) {
  _in = _in.toString();
  if (_in.length == 1) {
    _in = `0${_in}`
  }
  return _in
}

function isInbound () {
  return (this.direction === "inbound")
}

</script>

<style>
</style>
