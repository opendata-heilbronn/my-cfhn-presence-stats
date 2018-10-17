<template>
	<div class="body">
		<h1>My CFHN Presence Stats</h1>

		<h3>Top 10 (gesamt)</h3>

		<ul>
			<li v-for="v in user_total" v-bind="v">{{v.username}} ({{parseInt(v.visits * 5 / 60, 10)}}h {{v.visits * 5 % 60}}min)</li>
		</ul>

		<h3>Top 10 (letzte Woche)</h3>

		<ul>
			<li v-for="v in user_lastweek" v-bind="v">{{v.username}} ({{parseInt(v.visits * 5 / 60, 10)}}h {{v.visits * 5 % 60}}min)</li>
		</ul>
		</ul>

		<h3>Top 10 (Forever alone)</h3>

		<ul>
			<li v-for="v in user_alone" v-bind="v">{{v.username}} ({{parseInt(v.visits * 5 / 60, 10)}}h {{v.visits * 5 % 60}}min)</li>
		</ul>

		<h3>Die letzte Woche im Überblick</h3>

		<ul>
			<li style="list-style: none;">
				<span style="display: inline-block; width: 200px; text-align: right;">Uhrzeit</span>
				<template v-for="hour in hours">
					<span style="display: inline-block; width: 20px; height: 20px; margin-right: 4px; text-align: center;">{{hour}}</span>
				</template>
			</li>
			<li v-for="d in overview_lastweek" v-bind="d" style="list-style: none;">
				<span style="display: inline-block; width: 200px; text-align: right;">{{d.day}}</span>
				<template v-for="h in d.visits">
					<span style="display: inline-block; width: 20px; height: 20px; margin-right: 4px;" :style="{ 'background-color': h.color }" :title="h.visits"></span>
				</template>
			</li>
		</ul>
	</div>
</template>

<script>
	import axios from 'axios';

	export default {
		name: 'app',

		data: function () {
			return {
				user_total: {},
				user_lastweek: {},
				overview_lastweek: {},
				hours: new Array(24).fill(0).map(function(x, i) { return i; })
			};
		},

		computed: {
		},

		methods: {
			getColor: function(visits) {
				switch (visits) {
					case 0:
						return "#ebedf0";
					case 1:
						return "#c6e48b";
					case 2:
					case 3:
						return "#7bc96f";
					case 4:
					case 5:
					case 6:
						return "#239a3b";
					default:
						return "#196127";
				}
			},
			_fetch: function() {
				axios
					.get('/api')
					.then(response => {
						this.user_total = response.data.user_total;
						this.user_lastweek = response.data.user_lastweek;
						this.user_alone = response.data.user_alone;
						this.overview_lastweek = {};
						response.data.overview_lastweek.forEach(function(timeVisit) {
							if (!this.overview_lastweek.hasOwnProperty(timeVisit.day)) {
								this.overview_lastweek[timeVisit.day] = {
									day: timeVisit.day,
									visits: new Array(24).fill({visits: 0, color: this.getColor(0)})
								};
							}
							this.overview_lastweek[timeVisit.day].visits[timeVisit.hour] = {
								visits: timeVisit.visits,
								color: this.getColor(timeVisit.visits)
							};
						}.bind(this));
						console.log(this.overview_lastweek);
					})
					.catch(err => {
						console.error(err);
							console.log(arguments);
					});
			}
		},

		mounted: function () {
			this._fetch();
		}
	}
</script>
