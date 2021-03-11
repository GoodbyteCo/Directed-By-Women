<template>
	<div id="app">
		<Film />
	</div>
</template>

<script>
import Film from "./components/Film";

export default {
	name: "App",
	components: {
		Film
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
</style>
