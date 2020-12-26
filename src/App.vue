<template>
	<div id="app">
		<div id="main">
			<Film />
		</div>
		<Footer />
	</div>
</template>

<script>
import Film from "./components/Film";
import Footer from "./components/Footer";

export default {
	name: "App",
	components: {
		Film,
		Footer
	}
};
</script>

<style>
#app {
	font-family: "IBM Plex Mono", monospace;
	-webkit-font-smoothing: antialiased;
	-moz-osx-font-smoothing: grayscale;
	text-align: left;
	color: #2c3e50;
	margin-top: 60px;
}

#main {
	min-height: calc(75vh - 15px);
	min-height: -webkit-calc(75vh - 15px);
	min-height: -moz-calc(75vh - 15px);
}
</style>
