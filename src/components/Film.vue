<template>
	<div>
		<div class="input">
			<label for="userbox">Username(s):</label>
			<input
				class="userfield"
				placeholder="ex: holopollock, qjack"
				id="userbox"
				type="text"
				v-on:keyup.enter="submit()"
				v-model="users"
			/>
			<button v-on:click="submit()">Submit</button>
			<div v-if="started" class="output">
				<p>Loading...</p>
			</div>
			<div v-if="done" class="output">
				<p>Films Directed by women {{ women }}</p>
				<p>Total Films {{ total }}</p>
				<p>Percentage of films directed by women {{ percentage }}</p>
			</div>
		</div>
	</div>
</template>

<style scoped>
label
{
	display: block;
	margin: 10px 0;
}

input
{
	font-family: "IBM Plex Mono", monospace;
	font-size: 1rem;
	line-height: 2;
	color: inherit;
	border: none;
	border-bottom: 1px currentColor solid;
	outline: none;
}

button
{
	cursor: pointer;
	padding: 0 1ch;
	margin-left: 20px;
	
	font-family: "IBM Plex Mono", monospace;
	font-size: 1rem;
	line-height: 2;
	color: inherit;
	background: transparent;
	border: 1px solid currentColor;
	outline: none;
}

button:hover, button:focus-visible
{
	background: #c4c4c4;
}

.output
{
	margin-top: 60px;
}
</style>

<script>
export default {
	name: "Film",
	data() {
		return {
			users: "",
			info: "",
			started: false,
			done: false
		};
	},
	methods: {
		submit() {
			this.started = true;
			let inputted = this.users.split(/(?:,| )+/); //split input field on space or comma
			let userlist = inputted.filter(function(el) {
				return el;
			});
			let url = "/api?users=" + userlist.join("&users=");
			try {
				let vue = this;
				fetch(url)
					.then(function(res) {
						if (res.status != 200) {
							return "";
						}
						return res.json();
					})
					.then(function(json) {
						vue.info = json;
						console.log(vue.info);
						vue.done = true;
					});
			} catch (e) {
				this.$alert(
					"Something went wrong. Please try again in a moment. Error:" +
						e,
					"An error occured"
				);
			}
		}
	},
	computed: {
		women: function() {
			return this.info.women;
		},
		total: function() {
			return this.info.total;
		},
		percentage: function() {
			return this.info.percentage;
		}
	}
};
</script>
