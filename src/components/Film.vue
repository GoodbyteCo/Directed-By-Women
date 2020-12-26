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
			<div v-if="done">
				<p>Films Directed by women {{ women }}</p>
				<p>Total Films {{ total }}</p>
				<p>Percentage of films directed by women {{ percentage }}</p>
			</div>
		</div>
	</div>
</template>
<script>
export default {
	name: "Film",
	data() {
		return {
			users: "",
			info: "",
			done: false
		};
	},
	methods: {
		submit() {
			let inputted = this.users.split(/(?:,| )+/); //split input field on space or comma
			let userlist = inputted.filter(function(el) {
				return el;
			});
			let url =
				"http://127.0.0.1:8080/api?users=" + userlist.join("&users=");
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