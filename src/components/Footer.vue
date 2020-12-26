<template>
	<footer>
		<a href="https://ko-fi.com/goodbyte" class="support"
			>Support this site</a
		>
		<a href="http://goodbyte.ca/">Goodbyte</a> &copy; 2020
	</footer>
</template>

<script>
export default {
	name: "Footer"
};
</script>

<style scoped>
footer {
	text-align: center;
	margin: 0 auto 2rem;
	font-size: 12px;
	letter-spacing: 0.03em;
}

a {
	color: inherit;
	text-decoration: none;
}

a:hover {
	text-decoration: underline;
}

a.support {
	display: block;
	font-family: "IBM Plex Mono", monospace;
	font-size: 1rem;
	padding: 0 1ch;
	text-align: center;
	margin: 0.5rem auto;
	width: 180px;
	text-decoration: none;
	background: transparent;
	border: 1px solid currentColor;
	line-height: 2;
	outline: none;
	font-weight: 400;
}

a.support:hover,
a.support:focus-visible {
	background: #eaeaea;
}
</style>
