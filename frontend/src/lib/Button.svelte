<script lang="ts">
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();
</script>

<button class="rounded-xl bg-green px-4 py-2 text-black transition-all hover:bg-primary">
	{@render children()}
</button>
