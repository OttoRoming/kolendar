<script lang="ts">
	let { value }: { value: string } = $props();
</script>

<input
	type="submit"
	class="w-80 rounded-xl bg-green px-4 py-2 text-black transition-all hover:bg-primary"
	{value}
/>
