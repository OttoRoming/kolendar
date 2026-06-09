<script lang="ts">
	let { label, value = $bindable() }: { label: string; value: string } = $props();
</script>

<input
	type="text"
	class="w-80 rounded-xl bg-black px-4 py-2"
	bind:value
	aria-label={label}
	placeholder={label}
/>
