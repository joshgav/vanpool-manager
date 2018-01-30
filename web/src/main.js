import axios from 'axios'
import Vue from 'vue'

const apiBasePath = '/api/v1'

Vue.component("my-trip", {

  template: require('./my-trip.html'),

  data: function data() { return {
    riders: [],
    user: {},
    date: this.formatDate(new Date()),
    direction: "Inbound"
  }},

  watch: {
    date: {
      handler: function() {this.refresh(false)},
    },
    direction: {
      handler: function() {this.refresh(false)},
    },
  },

  mounted: function mounted() {
    this.refresh(true)
  },

  computed: {
    present: function present() {
      for (let rider of this.riders) {
        if (rider.username == this.user.username) {
          return true
        }
      }
      return false
    },

    isInbound: function isInbound() {
      return (this.direction === "Inbound" || this.direction === "I")
    },
  },

  methods: {
    makeRider: function makeRider(username, displayName, date, direction) {
      return {
        username: username,
        displayName: displayName,
        date: date,
        direction: direction
      }
    },

    refresh: function refresh(force) {
      if ((this.date != null) && (this.direction != null)) {
        this.refreshUser(force)
        this.refreshRiders(this.formatDate(this.date), this.parseDirection(this.direction))
      }
    },

    // @param date date in format 'yyyy-mm-dd'
    // @param direction 'O' or 'I'
    refreshRiders: function refreshRiders(date, direction) {
      var self = this
      let baseUrl = `${apiBasePath}/riders`
      let queryParams = `?date=${this.formatDate(date)}&direction=${direction}`
      axios.get(`${baseUrl}${queryParams}`, {resposeType:'json'})
        .then(response => {
          // console.log(`got riders: ${response.data}`)
          self.riders = response.data
        })
        .catch(err => {
          console.log(`failed to get riders: ${err}`)
          self.riders = []
        })
    },

    refreshUser: function refreshUser(force) {
      var self = this
      if (!force && this.user && this.user.username) {
        // console.log(`user already set: ${this.user.username}`)
        return
      }
      let baseUrl = `${apiBasePath}/user`
      axios.get(baseUrl, {resposeType:'json'})
        .then(response => {
         // console.log(`got current user: ${response.data.username}`)
          self.user = response.data
        })
        .catch(err => {
          console.log(`failed to get current user: ${err}`)
          self.user = new Rider()
        })
    },

    addRider: function addRider() {
      axios.put(
        `${apiBasePath}/riders`, this.getRider()
      )
      .then(response => {
        this.refresh(false)
      })
      .catch(err => {
        console.log(`failed to add rider: ${err}`)
      })
    },

    removeRider: function removeRider() {
      axios.post(
        `${apiBasePath}/riders/delete`, this.getRider()
      )
      .then(response => {
        this.refresh(false)
      })
      .catch(err => {
        console.log(`failed to remove rider: ${err}`)
      })
    },

    getRider: function getRider() {
      return this.makeRider(
        this.user.username,
        this.user.displayName,
        this.date,
        this.parseDirection(this.direction)
      )
    },

    formatDate: function formatDate(_date) {
      var date
      try { date = new Date(_date) } catch (e) {}
      if (date != null) {
        let year  = date.getUTCFullYear()
        let month = this.fixLength(date.getUTCMonth() + 1)
        let day   = this.fixLength(date.getUTCDate())
        let formattedDate = `${year}-${month}-${day}`
        return formattedDate
      }
      else {
        return _date
      }
    },

    fixLength: function fixLength(_in) {
      _in = _in.toString();
      if (_in.length == 1) {
        _in = `0${_in}`
      }
      return _in
    },

    parseDirection: function parseDirection(dir) {
      if (dir == "Outbound" || dir == "O") {
        return "O"
      } else {
        return "I"
      }
    },

  } // methods
}) // Vue.component

var vm = new Vue({
  el: '#vue-app',
});

