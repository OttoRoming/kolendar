<script lang="ts">
	let {
		label,
		type = 'text',
		value = $bindable()
	}: { label: string; value: string; type: 'text' | 'email' | 'password' } = $props();
</script>

<input
	{type}
	class="w-80 rounded-xl bg-black px-4 py-2"
	bind:value
	aria-label={label}
	placeholder={label}
/>
