<template>
	<div class="body">
		<h1>My CFHN Presence Stats</h1>
		<em>(generiert um {{generated_at}})</em>
		<br style="clear: both;" />

		<h2 v-html="current_year"></h2>

		<div style="width: 33%; float: left;">
				<h3>Top 10 (Besuchszeit)</h3>

				<ol>
					<li v-for="v in year_user_total" v-bind="v">{{v.username}} ({{parseInt(v.visits * 5 / 60, 10)}}h {{v.visits * 5 % 60}}min)</li>
				</ol>
		</div>

		<div style="width: 33%; float: left;">
				<h3>Top 10 (Besuchszeit letzte Woche)</h3>

				<ol>
					<li v-for="v in year_user_lastweek" v-bind="v">{{v.username}} ({{parseInt(v.visits * 5 / 60, 10)}}h {{v.visits * 5 % 60}}min)</li>
				</ol>
		</div>

		<div style="width: 33%; float: left;">
				<h3>Top 10 (Besuche)</h3>

				<ol>
					<li v-for="v in year_user_visits" v-bind="v">{{v.username}} ({{v.visits}})</li>
				</ol>
		</div>

		<div style="width: 49%; float: left;">
				<h3>Top 10 (Forever alone)</h3>

				<ol>
					<li v-for="v in year_user_alone" v-bind="v">{{v.username}} ({{parseInt(v.visits * 5 / 60, 10)}}h {{v.visits * 5 % 60}}min)</li>
				</ol>
		</div>

		<div style="width: 49%; float: left;">
				<h3>Top 10 (Längster Streaks)</h3>

				<ol>
					<li v-for="v in year_user_streaks" v-bind="v">{{v.username}} ({{parseInt(v.visits * 5 / 60, 10)}}h {{v.visits * 5 % 60}}min)</li>
				</ol>
		</div>

		<div style="width: 49%; float: left;">
				<h3>Top 10 (Tage nach Besuchszeit)</h3>

				<ol>
					<li v-for="v in year_day_visits" v-bind="v">{{v.username}} ({{parseInt(v.visits * 5 / 60, 10)}}h {{v.visits * 5 % 60}}min)</li>
				</ol>
		</div>

		<div style="width: 49%; float: left;">
				<h3>Top 10 (Tage nach Besuche)</h3>

				<ol>
					<li v-for="v in year_day_users" v-bind="v">{{v.username}} ({{v.visits}} Besucher)</li>
				</ol>
		</div>

		<div style="width: 100%; float: left;">
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

		<h2>Allzeit</h2>

		<div style="width: 49%; float: left;">
				<h3>Top 10 (Besuchszeit)</h3>

				<ol>
					<li v-for="v in total_user_total" v-bind="v">{{v.username}} ({{parseInt(v.visits * 5 / 60, 10)}}h {{v.visits * 5 % 60}}min)</li>
				</ol>
		</div>

		<div style="width: 49%; float: left;">
				<h3>Top 10 (Besuche)</h3>

				<ol>
					<li v-for="v in total_user_visits" v-bind="v">{{v.username}} ({{v.visits}})</li>
				</ol>
		</div>

		<div style="width: 49%; float: left;">
				<h3>Top 10 (Forever alone)</h3>

				<ol>
					<li v-for="v in total_user_alone" v-bind="v">{{v.username}} ({{parseInt(v.visits * 5 / 60, 10)}}h {{v.visits * 5 % 60}}min)</li>
				</ol>
		</div>

		<div style="width: 49%; float: left;">
				<h3>Top 10 (Längster Streaks)</h3>

				<ol>
					<li v-for="v in total_user_streaks" v-bind="v">{{v.username}} ({{parseInt(v.visits * 5 / 60, 10)}}h {{v.visits * 5 % 60}}min)</li>
				</ol>
		</div>

		<div style="width: 49%; float: left;">
				<h3>Top 10 (Tage nach Besuchszeit)</h3>

				<ol>
					<li v-for="v in total_day_visits" v-bind="v">{{v.username}} ({{parseInt(v.visits * 5 / 60, 10)}}h {{v.visits * 5 % 60}}min)</li>
				</ol>
		</div>

		<div style="width: 49%; float: left;">
				<h3>Top 10 (Tage nach Besuche)</h3>

				<ol>
					<li v-for="v in total_day_users" v-bind="v">{{v.username}} ({{v.visits}} Besucher)</li>
				</ol>
		</div>
	</div>
</template>

<script>
	import axios from 'axios';

	export default {
		name: 'app',

		data: function () {
			return {
				generated_at: '',
				current_year: '',
				year_user_total: {},
				year_user_lastweek: {},
				year_user_alone: {},
				year_user_visits: {},
				year_user_streaks: {},
				year_day_users: {},
				year_day_visist: {},
				total_user_total: {},
				total_user_alone: {},
				total_user_visits: {},
				total_user_streaks: {},
				total_day_users: {},
				total_day_visist: {},
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
						this.generated_at = response.data.generated_at;
						this.current_year = response.data.current_year;
						this.year_user_total = response.data.year_user_total;
						this.year_user_lastweek = response.data.year_user_lastweek;
						this.year_user_alone = response.data.year_user_alone;
						this.year_user_visits = response.data.year_user_visits;
						this.year_user_streaks = response.data.year_user_streaks;
						this.year_day_users = response.data.year_day_users;
						this.year_day_visits = response.data.year_day_visits;
						this.total_user_total = response.data.total_user_total;
						this.total_user_lastweek = response.data.total_user_lastweek;
						this.total_user_alone = response.data.total_user_alone;
						this.total_user_visits = response.data.total_user_visits;
						this.total_user_streaks = response.data.total_user_streaks;
						this.total_day_users = response.data.total_day_users;
						this.total_day_visits = response.data.total_day_visits;
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
