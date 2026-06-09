<script lang="ts">
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();
</script>

<header class="pb-4 text-4xl text-green">
	{@render children()}
</header>
