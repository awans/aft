<script>
	import { navStore } from './stores.js';
	import { canRoute } from '../app/router.js';


	let selected;
	navStore.subscribe(value => {
		selected = value;
	});
	let items = [
		{name:"Schema", path:'/schema', id:"schema"}, 
		{name:"Datatypes", path: '/datatypes', id:"datatype"},
		{name:"RPCs", path: '/rpcs', id:"rpcs"},
		{name:"Access", path:"/roles", id:"access"},
		{name:"Modules", path:'/modules', id:"modules"}, 
		{name:"Terminal", path: '/terminal',id:"terminal"},
		{name:"Records", path:"/records",id:"records"},
		{name:"Log", path:"/log/request",id:"log"},
	];

	let hide = false;
</script>

<style>
	.nav {
		height: 100%;
		flex-grow: 1;
		width: 9em;
	}
	ul {
		margin: 0;
		padding-top: 1em;
		padding-bottom: 1em;
		padding-left: 0;
		padding-right: 0;
		list-style-type: none;
	}
	li {
		padding-left: 1.5em;
		padding-bottom: .5em;
	}
	.nav-item {
		color: var(--text-color-darker);
		font-weight: 400;
		display:flex;
		align-items:center;
		flex-direction:row;
		transition: color .1s;
	}
	.nav-item:hover {
		color: var(--text-color);
		transition: color .1s;
	}

	.active {
		color: var(--text-color);
	}
	.noselect {
		cursor: pointer;
		-webkit-touch-callout: none; /* iOS Safari */
		-webkit-user-select: none; /* Safari */
		-khtml-user-select: none; /* Konqueror HTML */
		-moz-user-select: none; /* Old versions of Firefox */
		-ms-user-select: none; /* Internet Explorer/Edge */
		user-select: none; /* Non-prefixed version, currently supported by Chrome, Edge, Opera and Firefox */
	}
	.wrapper {
		height: 100%;
		display: flex;
		flex-direction: row;
	}
	.collapse-bar {
		height: 100%;
		width: 1em;
		background: var(--background);
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
	}
	.collapse-bar:hover {
		background: var(--background-highlight);
	}
	.affordance {
		color: var(--text-color-darker);
	}
	.collapse-bar:hover > .affordance {
		color: var(--text-color);
	}
</style>

<div class="wrapper">
	{#if !hide}
	<div class="nav">
		<ul>
			{#each items as item}
			<li>
				<div class="nav-item {selected ===item.id? 'active' : ''}">
					<a href="{item.path}" class="noselect" on:click={canRoute}>{item.name}</a>		
				</div>
			</li>
			{/each}
		</ul>
	</div>
	{/if}
	<div on:click={() => hide = !hide} class="collapse-bar">
		<div class="affordance">
			{#if hide}
			&rsaquo;
			{:else}
			&lsaquo;
			{/if}
		</div>
	</div>
</div>
