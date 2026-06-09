<script lang="ts">
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();
</script>

<header class="pb-8 text-6xl text-primary">
	{@render children()}
</header>
