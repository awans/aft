<style>
	button {
		-webkit-appearance: none;
		-moz-appearance: none;
		appearance: none;
		background-color: var(--background);
		border-radius: 0;
		margin: 0;
		color: var(--text-color-darker);
		border-color: var(--border-color);
		transition: background-color .1s, color .1s;
		cursor: pointer;
	}
	button:active {
		background-color: var(--background-highlight);
	}
	button:hover {
		background-color: var(--background-highlight);
		transition: background-color .1s, color .1s;
		color: var(--text-color);
	}

</style>
<button on:click>
	<slot />
</button>