import axios from 'axios'
import Vue from 'vue'

const apiHostname = 'localhost:8080'
const apiBasePath = '/api/v1/'
const apiScheme = 'http'

Vue.component("my-trip", {

  template: require('./my-trip.html'),

  data: function data() { return {
    riders: [],
    user: new this.Rider(),
    date: this.formatDate(new Date()),
    direction: "Inbound"
  }},

  watch: {
    date: {
      handler: function Handler() {this.refresh(false)},
    },
    direction: {
      handler: function Handler() {this.refresh(false)},
    },
  },

  mounted: function() {
    this.refresh(true)
  },

  methods: {

    Rider: function Rider(username, displayName, date, direction) {
      this.username = "" + username;
      this.displayName = "" + displayName;
      this.date = this.formatDate(date ? date : new Date());
      this.direction = direction;
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
      let baseUrl = `${apiScheme}://${apiHostname}${apiBasePath}riders`
      let queryParams = `?date=${this.formatDate(date)}&direction=${direction}`
      axios.get(`${baseUrl}${queryParams}`, {resposeType:'json'})
        .then(response => {
          console.log(`got riders: ${response.data}`)
          self.riders = response.data
        })
        .catch(err => {
          console.log(`failed to get riders: ${err}`)
          self.riders = []
        })
    },

    present: function present() {
      for (var rider in this.riders) {
        if (rider.username == this.user.username) {
          return true
        }
      }
      return false
    },

    refreshUser: function refreshUser(force) {
      var self = this
      if (!force && this.user && this.user.username) { return }
      let baseUrl = `${apiScheme}://${apiHostname}${apiBasePath}user`
      axios.get(baseUrl, {resposeType:'json'})
        .then(response => {
          console.log(`got current user: ${response.data}`)
          self.user = response.data
        })
        .catch(err => {
          console.log(`failed to get current user: ${err}`)
          self.user = new Rider()
        })
    },

    addRider: function addRider() {
      axios.put(
        `${apiScheme}://${apiHostname}${apiBasePath}riders`,
        JSON.stringify(this.getRider())
       )
      .catch(err => {
        console.log(`failed to add rider: ${err}`)
      })
    },

    removeRider: function removeRider() {
      axios.delete(
        `${apiScheme}://${apiHostname}${apiBasePath}riders`,
        JSON.stringify(this.getRider())
      )
      .catch(err => {
        console.log(`failed to remove rider: ${err}`)
      })
    },

    getRider: function getRider() {
      return new this.Rider(
        this.user.username,
        this.user.displayName,
        this.formatDate(this.date),
        this.parseDirection(this.direction),
      )
    },

    formatDate: function formatDate(_date) {
      console.log(`formatting date: ${_date}`)
      var date
      try { date = new Date(_date) } catch (e) {}
      if (date != null) {
        let year  = date.getUTCFullYear()
        let month = this.fixLength(date.getUTCMonth() + 1)
        let day   = this.fixLength(date.getUTCDate())
        let formattedDate = `${year}-${month}-${day}`
        console.log(`as: ${formattedDate}`)
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

    isInbound: function isInbound() {
      return (this.direction === "inbound")
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

