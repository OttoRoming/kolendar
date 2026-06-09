<script lang="ts">
	import '@fontsource-variable/eb-garamond/wght.css';
	import './layout.css';

	import favicon from '$lib/assets/favicon.svg';
	import { resolve } from '$app/paths';

	import { Toaster } from 'svelte-sonner';

	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<div class="grid min-h-screen grid-rows-[auto_1fr]">
	<nav class="mb-2 flex justify-between bg-black px-5 py-2 text-foreground shadow-xl">
		<a class="text-2xl font-bold" href={resolve('/')}>Kolendar</a>
	</nav>
	<main class="flex flex-col items-center text-foreground">
		{@render children()}
	</main>
	<Toaster position="top-right" />
</div>

<style>
	:global(body) {
		background-color: var(--color-background);
	}
</style>
